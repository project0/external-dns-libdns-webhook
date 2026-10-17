package provider_test

import (
	"slices"
	"testing"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/project0/external-dns-libdns-webhook/internal/provider"
)

const zone = "example.com."

func newWebhook(t *testing.T, zones []string, libdnsProvider libdnsregistry.Provider, opts provider.Options) *provider.WebhookProvider {
	t.Helper()

	webhook, err := provider.New(t.Context(), zones, libdnsProvider, opts)
	if err != nil {
		t.Fatal(err)
	}

	return webhook
}

func ep(name, recordType string, ttl int64, targets ...string) *externaldns.Endpoint {
	return &externaldns.Endpoint{DNSName: name, RecordType: recordType, RecordTTL: ttl, Targets: targets}
}

// endpointKey identifies an endpoint the way external-dns compares records.
func endpointKey(e *externaldns.Endpoint) string {
	return e.DNSName + "/" + e.RecordType
}

// assertEndpoints compares name, type, TTL and the set of targets.
func assertEndpoints(t *testing.T, got, want externaldns.Endpoints) {
	t.Helper()

	index := map[string]*externaldns.Endpoint{}
	for _, e := range got {
		index[endpointKey(e)] = e
	}

	for _, wantEp := range want {
		gotEp, ok := index[endpointKey(wantEp)]
		if !ok {
			t.Errorf("missing endpoint %s", endpointKey(wantEp))

			continue
		}

		delete(index, endpointKey(wantEp))

		gotTargets, wantTargets := slices.Sorted(slices.Values(gotEp.Targets)), slices.Sorted(slices.Values(wantEp.Targets))
		if !slices.Equal(gotTargets, wantTargets) || gotEp.RecordTTL != wantEp.RecordTTL || len(gotEp.ProviderSpecific) != 0 {
			t.Errorf("endpoint %s: got ttl=%d targets=%q providerSpecific=%v, want ttl=%d targets=%q",
				endpointKey(wantEp), gotEp.RecordTTL, gotTargets, gotEp.ProviderSpecific, wantEp.RecordTTL, wantTargets)
		}
	}

	for key := range index {
		t.Errorf("unexpected endpoint %s", key)
	}
}

func TestRecordsGrouping(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	mem.seed(t, zone,
		libdns.RR{Name: "@", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
		libdns.RR{Name: "@", Type: "A", TTL: time.Hour, Data: "192.0.2.2"},
		libdns.RR{Name: "@", Type: "MX", TTL: time.Hour, Data: "20 backup.example.com."},
		libdns.RR{Name: "@", Type: "MX", TTL: time.Hour, Data: "10 mail.example.com."},
		libdns.RR{Name: "@", Type: "SOA", TTL: time.Hour, Data: "ns1.example.net. hostmaster.example.com. 1 7200 3600 1209600 3600"},
		libdns.RR{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns1.example.net."},
		libdns.RR{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns2.example.net"},
		libdns.RR{Name: "_sip._tcp", Type: "SRV", TTL: 5 * time.Minute, Data: "10 5 5060 sip.example.com."},
		libdns.RR{Name: "www", Type: "CNAME", TTL: 5 * time.Minute, Data: "example.com."},
		libdns.RR{Name: "WWW6", Type: "AAAA", TTL: 5 * time.Minute, Data: "2001:db8::1"},
		libdns.RR{Name: "txt", Type: "TXT", TTL: 5 * time.Minute, Data: "heritage=external-dns,external-dns/owner=default"},
	)

	endpoints, err := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour}).Records(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	assertEndpoints(t, endpoints, externaldns.Endpoints{
		ep("example.com", "A", 3600, "192.0.2.1", "192.0.2.2"),
		ep("example.com", "MX", 3600, "10 mail.example.com", "20 backup.example.com"),
		ep("sub.example.com", "NS", 3600, "ns1.example.net", "ns2.example.net"),
		ep("_sip._tcp.example.com", "SRV", 300, "10 5 5060 sip.example.com."),
		ep("www.example.com", "CNAME", 300, "example.com"),
		ep("www6.example.com", "AAAA", 300, "2001:db8::1"),
		ep("txt.example.com", "TXT", 300, "heritage=external-dns,external-dns/owner=default"),
	})
}

func TestRecordsZoneError(t *testing.T) {
	t.Parallel()

	_, err := newWebhook(t, []string{"example.com", "missing.org"}, newMemoryProvider(zone), provider.Options{}).Records(t.Context())
	if err == nil {
		t.Fatal("expected error for unknown zone")
	}
}

func TestCreateRoundTrip(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider(zone)
	webhook := newWebhook(t, []string{"example.com"}, mem, provider.Options{DefaultTTL: time.Hour})

	desired := externaldns.Endpoints{
		ep("example.com", "MX", 3600, "10 mail.example.com", "20 backup.example.com."),
		ep("_sip._tcp.example.com", "SRV", 300, "10 5 5060 sip.example.com."),
		ep("_xmpp._tcp.chat.example.com", "SRV", 300, "5 0 5222 xmpp.example.com."),
		ep("sub.example.com", "NS", 3600, "ns1.example.net", "ns2.example.net"),
		ep("example.com", "TXT", 300, `"v=spf1 -all"`),
		ep("www.example.com", "A", 300, "192.0.2.1", "192.0.2.2"),
	}

	if err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{Create: desired}); err != nil {
		t.Fatal(err)
	}

	stored := mem.rrs(t, zone)
	for _, want := range []libdns.RR{
		{Name: "@", Type: "MX", TTL: time.Hour, Data: "10 mail.example.com."},
		{Name: "_sip._tcp", Type: "SRV", TTL: 5 * time.Minute, Data: "10 5 5060 sip.example.com."},
		{Name: "_xmpp._tcp.chat", Type: "SRV", TTL: 5 * time.Minute, Data: "5 0 5222 xmpp.example.com."},
		{Name: "sub", Type: "NS", TTL: time.Hour, Data: "ns1.example.net."},
		{Name: "@", Type: "TXT", TTL: 5 * time.Minute, Data: "v=spf1 -all"},
	} {
		if !slices.Contains(stored, want) {
			t.Errorf("record %+v not stored, got %+v", want, stored)
		}
	}

	endpoints, err := webhook.Records(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	assertEndpoints(t, endpoints, externaldns.Endpoints{
		ep("example.com", "MX", 3600, "10 mail.example.com", "20 backup.example.com"),
		ep("_sip._tcp.example.com", "SRV", 300, "10 5 5060 sip.example.com."),
		ep("_xmpp._tcp.chat.example.com", "SRV", 300, "5 0 5222 xmpp.example.com."),
		ep("sub.example.com", "NS", 3600, "ns1.example.net", "ns2.example.net"),
		ep("example.com", "TXT", 300, "v=spf1 -all"),
		ep("www.example.com", "A", 300, "192.0.2.1", "192.0.2.2"),
	})
}

// withoutZoneLister hides the ListZones method of a provider.
type withoutZoneLister struct {
	libdnsregistry.Provider
}

func TestZoneDiscovery(t *testing.T) {
	t.Parallel()

	mem := newMemoryProvider("example.com.", "sub.example.com.", "example.org.")

	webhook := newWebhook(t, nil, mem, provider.Options{})
	if got, want := webhook.GetDomainFilter(t.Context()).Include, []string{"sub.example.com", "example.com", "example.org"}; !slices.Equal(got, want) {
		t.Errorf("discovered zones %v, want %v", got, want)
	}

	webhook = newWebhook(t, []string{"example.org"}, mem, provider.Options{})
	if got := webhook.GetDomainFilter(t.Context()).Include; !slices.Equal(got, []string{"example.org"}) {
		t.Errorf("configured zones must take precedence, got %v", got)
	}

	if _, err := provider.New(t.Context(), nil, withoutZoneLister{mem}, provider.Options{}); err == nil {
		t.Error("expected error without zones and zone lister")
	}
}
