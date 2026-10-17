package provider_test

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/libdns/libdns"
)

var errCNAMEConflict = errors.New("CNAME and other data")

// memoryProvider is an in-memory libdns provider implementing the libdns v1 semantics.
// Like a real DNS server it rejects CNAME records next to other records of the same name.
type memoryProvider struct {
	mu    sync.Mutex
	zones map[string][]libdns.RR
	// minTTL raises lower TTLs of appended records
	minTTL time.Duration
	// ops records the called methods, e.g. "delete example.com."
	ops []string
}

func newMemoryProvider(zones ...string) *memoryProvider {
	m := &memoryProvider{zones: map[string][]libdns.RR{}}
	for _, zone := range zones {
		m.zones[zone] = nil
	}

	return m
}

func (m *memoryProvider) seed(t *testing.T, zone string, rrs ...libdns.RR) {
	t.Helper()

	m.mu.Lock()
	defer m.mu.Unlock()

	m.zones[zone] = append(m.zones[zone], rrs...)
}

func (m *memoryProvider) rrs(zone string) []libdns.RR {
	m.mu.Lock()
	defer m.mu.Unlock()

	return slices.Clone(m.zones[zone])
}

func (m *memoryProvider) GetRecords(_ context.Context, zone string) ([]libdns.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	rrs, ok := m.zones[zone]
	if !ok {
		return nil, errors.New("zone not found: " + zone)
	}

	records := make([]libdns.Record, 0, len(rrs))

	for _, rr := range rrs {
		record, err := rr.Parse()
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", rr.Name, err)
		}

		records = append(records, record)
	}

	return records, nil
}

func (m *memoryProvider) AppendRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ops = append(m.ops, "append "+zone)

	for _, record := range records {
		rr := record.RR()
		for _, existing := range m.zones[zone] {
			if existing.Name == rr.Name && existing.Type != rr.Type && (existing.Type == "CNAME" || rr.Type == "CNAME") {
				return nil, errCNAMEConflict
			}
		}

		rr.TTL = max(rr.TTL, m.minTTL)
		m.zones[zone] = append(m.zones[zone], rr)
	}

	return records, nil
}

func (m *memoryProvider) SetRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ops = append(m.ops, "set "+zone)

	for _, record := range records {
		rr := record.RR()
		m.zones[zone] = slices.DeleteFunc(m.zones[zone], func(existing libdns.RR) bool {
			return existing.Name == rr.Name && existing.Type == rr.Type
		})
	}

	for _, record := range records {
		m.zones[zone] = append(m.zones[zone], record.RR())
	}

	return records, nil
}

func (m *memoryProvider) DeleteRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ops = append(m.ops, "delete "+zone)

	var deleted []libdns.Record

	for _, record := range records {
		rr := record.RR()
		m.zones[zone] = slices.DeleteFunc(m.zones[zone], func(existing libdns.RR) bool {
			match := existing.Name == rr.Name &&
				(rr.Type == "" || rr.Type == existing.Type) &&
				(rr.TTL == 0 || rr.TTL == existing.TTL) &&
				(rr.Data == "" || rr.Data == existing.Data)
			if match {
				deleted = append(deleted, record)
			}

			return match
		})
	}

	return deleted, nil
}

func (m *memoryProvider) ListZones(_ context.Context) ([]libdns.Zone, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	zones := make([]libdns.Zone, 0, len(m.zones))
	for zone := range m.zones {
		zones = append(zones, libdns.Zone{Name: zone})
	}

	slices.SortFunc(zones, func(a, b libdns.Zone) int { return strings.Compare(a.Name, b.Name) })

	return zones, nil
}
