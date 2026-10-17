package fakememory_test

import (
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/fakememory"
)

const zone = "example.com."

func rrs(t *testing.T, provider *fakememory.Provider) []libdns.RR {
	t.Helper()

	records, err := provider.GetRecords(t.Context(), zone)
	if err != nil {
		t.Fatal(err)
	}

	result := make([]libdns.RR, 0, len(records))
	for _, record := range records {
		result = append(result, record.RR())
	}

	return result
}

func TestZones(t *testing.T) {
	t.Parallel()

	provider := &fakememory.Provider{Zones: []string{"Example.com", "example.org."}}

	zones, err := provider.ListZones(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if want := []libdns.Zone{{Name: "example.com."}, {Name: "example.org."}}; !slices.Equal(zones, want) {
		t.Errorf("zones %v, want %v", zones, want)
	}

	if _, err := provider.GetRecords(t.Context(), "example.net."); !errors.Is(err, fakememory.ErrZoneNotFound) {
		t.Errorf("expected zone not found, got %v", err)
	}
}

func TestAppendAndParse(t *testing.T) {
	t.Parallel()

	provider := &fakememory.Provider{Zones: []string{zone}, MinTTL: 300}

	_, err := provider.AppendRecords(t.Context(), zone, []libdns.Record{
		libdns.MX{Name: "@", TTL: time.Minute, Preference: 10, Target: "mail.example.com."},
		libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	records, err := provider.GetRecords(t.Context(), zone)
	if err != nil {
		t.Fatal(err)
	}

	mx, ok := records[0].(libdns.MX)
	if !ok || mx.Preference != 10 || mx.TTL != 5*time.Minute {
		t.Errorf("expected typed MX with TTL raised to 5m, got %#v", records[0])
	}

	if _, ok := records[1].(libdns.Address); !ok {
		t.Errorf("expected typed address, got %#v", records[1])
	}
}

func TestCNAMEConflict(t *testing.T) {
	t.Parallel()

	provider := &fakememory.Provider{Zones: []string{zone}}

	_, err := provider.AppendRecords(t.Context(), zone, []libdns.Record{
		libdns.RR{Name: "www", Type: "A", Data: "192.0.2.1"},
		libdns.RR{Name: "www", Type: "CNAME", Data: "example.com."},
	})
	if !errors.Is(err, fakememory.ErrCNAMEConflict) {
		t.Errorf("expected CNAME conflict, got %v", err)
	}
}

func TestSetRecords(t *testing.T) {
	t.Parallel()

	provider := &fakememory.Provider{Zones: []string{zone}}

	_, err := provider.AppendRecords(t.Context(), zone, []libdns.Record{
		libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
		libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.2"},
		libdns.RR{Name: "www", Type: "AAAA", TTL: time.Hour, Data: "2001:db8::1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = provider.SetRecords(t.Context(), zone, []libdns.Record{libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.3"}})
	if err != nil {
		t.Fatal(err)
	}

	want := []libdns.RR{
		{Name: "www", Type: "AAAA", TTL: time.Hour, Data: "2001:db8::1"},
		{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.3"},
	}
	if got := rrs(t, provider); !slices.Equal(got, want) {
		t.Errorf("records %v, want %v", got, want)
	}
}

func TestDeleteRecordsWildcards(t *testing.T) {
	t.Parallel()

	provider := &fakememory.Provider{Zones: []string{zone}}

	_, err := provider.AppendRecords(t.Context(), zone, []libdns.Record{
		libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
		libdns.RR{Name: "www", Type: "A", TTL: time.Hour, Data: "192.0.2.2"},
		libdns.RR{Name: "www", Type: "TXT", TTL: time.Hour, Data: "keep"},
		libdns.RR{Name: "other", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
	})
	if err != nil {
		t.Fatal(err)
	}

	// wrong TTL does not match
	deleted, err := provider.DeleteRecords(t.Context(), zone, []libdns.Record{libdns.RR{Name: "www", Type: "A", TTL: time.Minute, Data: "192.0.2.1"}})
	if err != nil || len(deleted) != 0 {
		t.Fatalf("deleted %v, err %v", deleted, err)
	}

	// empty TTL and data match all A records of the name
	deleted, err = provider.DeleteRecords(t.Context(), zone, []libdns.Record{libdns.RR{Name: "www", Type: "A"}})
	if err != nil || len(deleted) != 2 {
		t.Fatalf("deleted %v, err %v", deleted, err)
	}

	want := []libdns.RR{
		{Name: "www", Type: "TXT", TTL: time.Hour, Data: "keep"},
		{Name: "other", Type: "A", TTL: time.Hour, Data: "192.0.2.1"},
	}
	if got := rrs(t, provider); !slices.Equal(got, want) {
		t.Errorf("records %v, want %v", got, want)
	}
}
