//go:build integration

package provider_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/project0/external-dns-libdns-webhook/internal/provider"
)

// TestINWX runs against a real INWX zone, records are created below a random label and removed afterwards.
//
//	INWX_USERNAME=... INWX_PASSWORD=... INWX_ZONE=example.com go test -tags integration -run INWX -v ./internal/provider
//
// Optional: INWX_SHARED_SECRET (2FA), INWX_ENDPOINT_URL (e.g. https://api.ote.domrobot.com/jsonrpc/).
// A unique TXT record is also created on the zone apex.
//
//nolint:paralleltest // steps build on each other and share the live zone
func TestINWX(t *testing.T) {
	zoneName := strings.TrimSuffix(os.Getenv("INWX_ZONE"), ".")
	if os.Getenv("INWX_USERNAME") == "" || zoneName == "" {
		t.Skip("INWX_USERNAME, INWX_PASSWORD and INWX_ZONE are required")
	}

	config, err := json.Marshal(map[string]string{
		"username":      os.Getenv("INWX_USERNAME"),
		"password":      os.Getenv("INWX_PASSWORD"),
		"shared_secret": os.Getenv("INWX_SHARED_SECRET"),
		"endpoint_url":  os.Getenv("INWX_ENDPOINT_URL"),
	})
	if err != nil {
		t.Fatal(err)
	}

	libdnsProvider, err := libdnsregistry.New("inwx", [][]byte{config})
	if err != nil {
		t.Fatal(err)
	}

	webhook := newWebhook(t, []string{zoneName}, libdnsProvider, provider.Options{DefaultTTL: time.Hour, MinTTL: 5 * time.Minute})

	random := make([]byte, 4)
	_, _ = rand.Read(random)
	label := "edns-" + hex.EncodeToString(random) + "." + zoneName
	apexText := "external-dns-libdns-webhook test " + label

	owned := func(e *externaldns.Endpoint) bool {
		return strings.HasSuffix(e.DNSName, label) ||
			(e.DNSName == zoneName && e.RecordType == "TXT" && strings.Contains(strings.Join(e.Targets, " "), apexText))
	}

	t.Cleanup(func() { cleanupINWX(t, libdnsProvider, zoneName, label, apexText) })

	// records returns the test endpoints only, apex TXT reduced to the test target
	records := func(t *testing.T) externaldns.Endpoints {
		t.Helper()

		all, err := webhook.Records(t.Context())
		if err != nil {
			t.Fatal(err)
		}

		var result externaldns.Endpoints

		for _, endpoint := range all {
			if !owned(endpoint) {
				continue
			}

			if endpoint.DNSName == zoneName {
				endpoint.Targets = []string{apexText}
			}

			result = append(result, endpoint)
		}

		return result
	}

	// apply adjusts the endpoints, applies the changes and asserts that a following sync has nothing to do
	apply := func(t *testing.T, changes externaldns.Changes, desired externaldns.Endpoints) {
		t.Helper()

		var err error

		for _, list := range []*externaldns.Endpoints{&changes.Create, &changes.UpdateNew, &changes.UpdateOld, &changes.Delete} {
			if *list, err = webhook.AdjustEndpoints(t.Context(), *list); err != nil {
				t.Fatal(err)
			}
		}

		if err := webhook.ApplyChanges(t.Context(), &changes); err != nil {
			t.Fatal(err)
		}

		desired, err = webhook.AdjustEndpoints(t.Context(), desired)
		if err != nil {
			t.Fatal(err)
		}

		current := records(t)
		assertEndpoints(t, current, desired)

		if planned := planChanges(current, desired); len(planned) > 0 {
			t.Errorf("sync would plan changes: %v", planned)
		}
	}

	name := func(prefix string) string { return prefix + "." + label }

	initial := func() externaldns.Endpoints {
		return externaldns.Endpoints{
			ep(zoneName, "TXT", 300, `"`+apexText+`"`),
			ep(name("a"), "A", 300, "192.0.2.1", "192.0.2.2"),
			ep(name("aaaa"), "AAAA", 300, "2001:db8::1"),
			ep(name("cname"), "CNAME", 300, "www.example.com."),
			ep(name("mx"), "MX", 300, "10 mail.example.com.", "20 backup.example.com"),
			ep(name("_sip._tcp"), "SRV", 300, "10 5 5060 sip.example.com.", "20 0 5061 sip2.example.com."),
			ep(name("delegated"), "NS", 3600, "ns1.example.net", "ns2.example.net"),
			ep(name("txt"), "TXT", 300, `"v=spf1 -all"`, "hello world"),
		}
	}

	t.Run("create", func(t *testing.T) {
		apply(t, externaldns.Changes{Create: initial()}, initial())
	})

	t.Run("create again is idempotent", func(t *testing.T) {
		apply(t, externaldns.Changes{Create: initial()}, initial())
	})

	updated := func() externaldns.Endpoints {
		endpoints := initial()
		endpoints[1] = ep(name("a"), "A", 600, "192.0.2.2", "192.0.2.3")
		endpoints[4] = ep(name("mx"), "MX", 300, "5 mail.example.com", "20 backup.example.com")
		endpoints[5] = ep(name("_sip._tcp"), "SRV", 300, "10 5 5060 sip.example.com.")
		endpoints[6] = ep(name("delegated"), "NS", 3600, "ns1.example.net", "ns3.example.net")
		endpoints[3] = ep(name("cname"), "CNAME", 300, "www2.example.com")

		return endpoints
	}

	t.Run("update targets, priorities, TTL and CNAME", func(t *testing.T) {
		old, updates := initial(), updated()

		apply(t, externaldns.Changes{
			UpdateOld: externaldns.Endpoints{old[1], old[3], old[4], old[5], old[6]},
			UpdateNew: externaldns.Endpoints{updates[1], updates[3], updates[4], updates[5], updates[6]},
		}, updated())
	})

	t.Run("replace A by CNAME", func(t *testing.T) {
		desired := updated()
		desired[1] = ep(name("a"), "CNAME", 300, "www.example.com")

		apply(t, externaldns.Changes{
			Delete: externaldns.Endpoints{updated()[1]},
			Create: externaldns.Endpoints{ep(name("a"), "CNAME", 300, "www.example.com")},
		}, desired)
	})

	t.Run("delete", func(t *testing.T) {
		desired := updated()
		desired[1] = ep(name("a"), "CNAME", 300, "www.example.com")

		if err := webhook.ApplyChanges(t.Context(), &externaldns.Changes{Delete: desired}); err != nil {
			t.Fatal(err)
		}

		if left := records(t); len(left) > 0 {
			t.Errorf("records left after delete: %v", left)
		}
	})
}

func cleanupINWX(t *testing.T, libdnsProvider libdnsregistry.Provider, zoneName, label, apexText string) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	records, err := libdnsProvider.GetRecords(ctx, zoneName+".")
	if err != nil {
		t.Errorf("cleanup: %v", err)

		return
	}

	var leftovers []libdns.Record

	for _, record := range records {
		rr := record.RR()
		if strings.HasSuffix(libdns.AbsoluteName(rr.Name, zoneName), label) || (rr.Type == "TXT" && rr.Data == apexText) {
			leftovers = append(leftovers, record)
		}
	}

	if len(leftovers) == 0 {
		return
	}

	t.Logf("cleanup: deleting %d leftover records", len(leftovers))

	if _, err := libdnsProvider.DeleteRecords(ctx, zoneName+".", leftovers); err != nil {
		t.Errorf("cleanup: %v", err)
	}
}
