package provider_test

import (
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/provider"
)

// planChanges mimics the external-dns plan: it returns the endpoints that would be created or updated.
// Targets are compared case insensitive regardless of order, TTL only if configured, provider specific exactly.
func planChanges(current, desired externaldns.Endpoints) []string {
	index := map[string]*externaldns.Endpoint{}
	for _, e := range current {
		index[strings.ToLower(e.DNSName)+"/"+e.RecordType] = e
	}

	var changes []string

	for _, want := range desired {
		cur, ok := index[strings.ToLower(want.DNSName)+"/"+want.RecordType]
		if !ok {
			changes = append(changes, "create "+endpointKey(want))

			continue
		}

		curTargets := slices.Sorted(slices.Values(cur.Targets))
		wantTargets := slices.Sorted(slices.Values(want.Targets))
		sameTargets := len(curTargets) == len(wantTargets)

		for i := range curTargets {
			if sameTargets && !strings.EqualFold(curTargets[i], wantTargets[i]) {
				sameTargets = false
			}
		}

		if !sameTargets || (want.RecordTTL > 0 && want.RecordTTL != cur.RecordTTL) ||
			!slices.Equal(cur.ProviderSpecific, want.ProviderSpecific) {
			changes = append(changes, "update "+endpointKey(want))
		}
	}

	return changes
}

func TestAdjustEndpoints(t *testing.T) {
	t.Parallel()

	webhook := newWebhook(t, []string{"example.com"}, newMemoryProvider(zone),
		provider.Options{DefaultTTL: time.Hour, MinTTL: 5 * time.Minute})

	withSetID := ep("weighted.example.com", "A", 300, "192.0.2.1")
	withSetID.SetIdentifier = "eu"

	withProviderSpecific := ep("Mail.Example.com.", "MX", 0, "20 backup.example.com.", "10 mail.example.com", "invalid", "10 Mail.example.com.")
	withProviderSpecific.ProviderSpecific = externaldns.ProviderSpecific{{Name: "webhook/priority", Value: "10"}}

	adjusted, err := webhook.AdjustEndpoints(t.Context(), externaldns.Endpoints{
		withSetID,
		withProviderSpecific,
		ep("ptr.example.com", "PTR", 300, "host.example.com"),
		ep("www.example.org", "A", 300, "192.0.2.1"),
		ep("bad.example.com", "A", 300, "not-an-ip"),
		ep("_sip._tcp.example.com", "SRV", 60, "10 5 5060 sip.example.com"),
		ep("txt.example.com", "TXT", 7200, `"heritage=external-dns"`),
	})
	if err != nil {
		t.Fatal(err)
	}

	assertEndpoints(t, adjusted, externaldns.Endpoints{
		ep("mail.example.com", "MX", 3600, "10 mail.example.com", "20 backup.example.com"),
		ep("_sip._tcp.example.com", "SRV", 300, "10 5 5060 sip.example.com."),
		ep("txt.example.com", "TXT", 7200, "heritage=external-dns"),
	})
}

func TestNoOpPlanAfterApply(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone, "sub.example.org.")
	// unsupported record types are neither returned nor touched
	soa := libdns.RR{Name: "@", Type: "SOA", TTL: time.Hour, Data: "ns1.example.net. hostmaster.example.org. 1 7200 3600 1209600 3600"}
	mem.seed(t, "sub.example.org.", soa)
	webhook := newWebhook(t, []string{"example.com", "sub.example.org"}, mem, provider.Options{DefaultTTL: time.Hour})

	desiredSource := func() externaldns.Endpoints {
		withProviderSpecific := ep("www.example.com", "A", 0, "192.0.2.2", "192.0.2.1")
		withProviderSpecific.ProviderSpecific = externaldns.ProviderSpecific{{Name: "webhook/foo", Value: "bar"}}

		return externaldns.Endpoints{
			withProviderSpecific,
			ep("v6.example.com", "AAAA", 300, "2001:DB8::1"),
			ep("alias.example.com", "CNAME", 300, "WWW.example.com."),
			ep("example.com", "MX", 3600, "10 mail.example.com.", "20 backup.example.com"),
			ep("_sip._tcp.example.com", "SRV", 300, "10 5 5060 sip.example.com.", "20 0 5060 sip2.example.com"),
			ep("delegated.example.com", "NS", 86400, "ns1.example.net.", "ns2.example.net"),
			ep("example.com", "TXT", 300, `"v=spf1 include:example.net -all"`),
			ep("DEEP.sub.example.org", "TXT", 0, "plain text"),
		}
	}

	adjusted, err := webhook.AdjustEndpoints(t.Context(), desiredSource())
	if err != nil {
		t.Fatal(err)
	}

	if err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{Create: adjusted}); err != nil {
		t.Fatal(err)
	}

	wantOrg := []libdns.RR{soa, {Name: "deep", Type: "TXT", TTL: time.Hour, Data: "plain text"}}
	if got := mem.rrs("sub.example.org."); !slices.Equal(got, wantOrg) {
		t.Errorf("sub.example.org records %+v, want %+v", got, wantOrg)
	}

	current, err := webhook.Records(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	// every sync adjusts a fresh copy of the desired endpoints
	adjusted, err = webhook.AdjustEndpoints(t.Context(), desiredSource())
	if err != nil {
		t.Fatal(err)
	}

	if changes := planChanges(current, adjusted); len(changes) > 0 {
		t.Errorf("expected no changes after apply, got %v\ncurrent: %v", changes, current)
	}
}

func TestNoOpPlanWithProviderMinTTL(t *testing.T) {
	t.Parallel()

	// emulates a provider raising TTLs to its minimum, like INWX does with 300 seconds
	mem := newMemoryProvider(zone)
	mem.minTTL = 5 * time.Minute
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour, MinTTL: 5 * time.Minute})

	adjusted, err := webhook.AdjustEndpoints(t.Context(), externaldns.Endpoints{ep("www.example.com", "A", 60, "192.0.2.1")})
	if err != nil {
		t.Fatal(err)
	}

	if err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{Create: adjusted}); err != nil {
		t.Fatal(err)
	}

	current, err := webhook.Records(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if changes := planChanges(current, adjusted); len(changes) > 0 {
		t.Errorf("expected no changes, got %v", changes)
	}
}

func TestApplyChangesReplaceAWithCNAME(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	mem.seed(t, zone, libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.1"})
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

	err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{
		Create: externaldns.Endpoints{ep("www.example.com", "CNAME", 3600, "example.com")},
		Delete: externaldns.Endpoints{ep("www.example.com", "A", 3600, "192.0.2.1")},
	})
	if err != nil {
		t.Fatal(err)
	}

	if want := []string{"delete " + zone, "append " + zone}; !slices.Equal(mem.ops, want) {
		t.Errorf("ops %v, want %v", mem.ops, want)
	}

	want := []libdns.RR{{Name: "www", Type: "CNAME", TTL: time.Hour, Data: "example.com."}}
	if got := mem.rrs(zone); !slices.Equal(got, want) {
		t.Errorf("records %+v, want %+v", got, want)
	}
}

func TestApplyChangesUpdate(t *testing.T) {
	t.Parallel()

	seed := []libdns.RR{
		{Name: "www", Type: "A", TTL: 5 * time.Minute, Data: "192.0.2.1"},
		{Name: "www", Type: "A", TTL: 5 * time.Minute, Data: "192.0.2.2"},
		{Name: "@", Type: "MX", TTL: time.Hour, Data: "10 mail.example.com."},
		{Name: "other", Type: "A", TTL: 5 * time.Minute, Data: "192.0.2.9"},
	}

	tests := []struct {
		name    string
		old     *externaldns.Endpoint
		new     *externaldns.Endpoint
		want    []libdns.RR
		deleted int
	}{
		{
			name: "replace one target",
			old:  ep("www.example.com", "A", 300, "192.0.2.1", "192.0.2.2"),
			new:  ep("www.example.com", "A", 300, "192.0.2.2", "192.0.2.3"),
			want: []libdns.RR{
				{Name: "www", Type: "A", TTL: 5 * time.Minute, Data: "192.0.2.2"},
				{Name: "www", Type: "A", TTL: 5 * time.Minute, Data: "192.0.2.3"},
			},
		},
		{
			name: "change TTL",
			old:  ep("www.example.com", "A", 300, "192.0.2.1", "192.0.2.2"),
			new:  ep("www.example.com", "A", 600, "192.0.2.1", "192.0.2.2"),
			want: []libdns.RR{
				{Name: "www", Type: "A", TTL: 10 * time.Minute, Data: "192.0.2.1"},
				{Name: "www", Type: "A", TTL: 10 * time.Minute, Data: "192.0.2.2"},
			},
		},
		{
			name: "change MX priority and add target",
			old:  ep("example.com", "MX", 3600, "10 mail.example.com"),
			new:  ep("example.com", "MX", 3600, "5 mail.example.com", "20 backup.example.com"),
			want: []libdns.RR{
				{Name: "@", Type: "MX", TTL: time.Hour, Data: "5 mail.example.com."},
				{Name: "@", Type: "MX", TTL: time.Hour, Data: "20 backup.example.com."},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mem := newMemoryProvider(zone)
			mem.seed(t, zone, seed...)
			webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

			err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{
				UpdateOld: externaldns.Endpoints{tc.old},
				UpdateNew: externaldns.Endpoints{tc.new},
			})
			if err != nil {
				t.Fatal(err)
			}

			var got []libdns.RR

			for _, rr := range mem.rrs(zone) {
				if rr.Name == "other" {
					continue
				}

				if rr.Type == tc.new.RecordType {
					got = append(got, rr)
				}
			}

			for _, want := range tc.want {
				if !slices.Contains(got, want) {
					t.Errorf("missing %+v in %+v", want, got)
				}
			}

			if len(got) != len(tc.want) {
				t.Errorf("got %d records, want %d: %+v", len(got), len(tc.want), got)
			}

			if !slices.Contains(mem.rrs(zone), seed[3]) {
				t.Error("unrelated record was modified")
			}
		})
	}
}

func TestApplyChangesDeleteListedTargetsOnly(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	mem.seed(t, zone,
		// providers may return host names without trailing dot
		libdns.RR{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns1.example.net"},
		libdns.RR{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns2.example.net."},
	)
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

	err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{
		Delete: externaldns.Endpoints{ep("sub.example.com", "NS", 0, "ns1.example.net")},
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []libdns.RR{{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns2.example.net."}}
	if got := mem.rrs(zone); !slices.Equal(got, want) {
		t.Errorf("records %+v, want %+v", got, want)
	}
}

func TestApplyChangesCreateIsIdempotent(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	mem.seed(t, zone, libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.1"})
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

	err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{
		Create: externaldns.Endpoints{ep("www.example.com", "A", 0, "192.0.2.1", "192.0.2.2")},
	})
	if err != nil {
		t.Fatal(err)
	}

	if got := mem.rrs(zone); len(got) != 2 {
		t.Errorf("expected 2 records, got %+v", got)
	}
}

func TestApplyChangesPartialFailure(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

	err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{
		Create: externaldns.Endpoints{
			ep("www.example.org", "A", 0, "192.0.2.1"),
			ep("ptr.example.com", "PTR", 0, "host.example.com"),
			ep("www.example.com", "A", 0, "192.0.2.1"),
		},
	})
	if err == nil {
		t.Fatal("expected error for unsupported endpoints")
	}

	if got := mem.rrs(zone); len(got) != 1 {
		t.Errorf("valid endpoint must be applied, got %+v", got)
	}
}
