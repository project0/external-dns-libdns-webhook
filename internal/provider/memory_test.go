package provider_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/fakememory"
)

// memoryProvider wraps the fakememory provider and records the called methods.
type memoryProvider struct {
	*fakememory.Provider

	mu sync.Mutex
	// ops records the called methods, e.g. "delete example.com."
	ops []string
}

func newMemoryProvider(zones ...string) *memoryProvider {
	return &memoryProvider{Provider: &fakememory.Provider{Zones: zones}}
}

func (m *memoryProvider) withMinTTL(ttl time.Duration) *memoryProvider {
	m.MinTTL = int(ttl.Seconds())

	return m
}

func (m *memoryProvider) record(op, zone string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ops = append(m.ops, op+" "+zone)
}

func (m *memoryProvider) seed(t *testing.T, zone string, rrs ...libdns.RR) {
	t.Helper()

	records := make([]libdns.Record, 0, len(rrs))
	for _, rr := range rrs {
		records = append(records, rr)
	}

	if _, err := m.Provider.AppendRecords(t.Context(), zone, records); err != nil {
		t.Fatal(err)
	}
}

func (m *memoryProvider) rrs(t *testing.T, zone string) []libdns.RR {
	t.Helper()

	records, err := m.GetRecords(t.Context(), zone)
	if err != nil {
		t.Fatal(err)
	}

	rrs := make([]libdns.RR, 0, len(records))
	for _, record := range records {
		rrs = append(rrs, record.RR())
	}

	return rrs
}

func (m *memoryProvider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.record("append", zone)

	return wrap(m.Provider.AppendRecords(ctx, zone, records))
}

func (m *memoryProvider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.record("set", zone)

	return wrap(m.Provider.SetRecords(ctx, zone, records))
}

func (m *memoryProvider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.record("delete", zone)

	return wrap(m.Provider.DeleteRecords(ctx, zone, records))
}

func wrap(records []libdns.Record, err error) ([]libdns.Record, error) {
	if err != nil {
		return records, fmt.Errorf("fakememory: %w", err)
	}

	return records, nil
}
