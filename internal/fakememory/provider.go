// Package fakememory implements a libdns provider keeping records in memory. It follows the libdns v1
// semantics and is meant for testing, records are lost on restart.
package fakememory

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/libdns/libdns"
)

var (
	ErrZoneNotFound  = errors.New("zone not found")
	ErrCNAMEConflict = errors.New("CNAME records cannot coexist with other records of the same name")
)

// Provider keeps records of the configured zones in memory.
// Like a DNS server it rejects CNAME records next to other records of the same name.
type Provider struct {
	// Zones served by the provider, e.g. ["example.com"]. Records of other zones are rejected.
	Zones []string `json:"zones,omitempty"`

	// MinTTL in seconds, lower TTLs are raised on append, emulating providers with a minimum TTL.
	MinTTL int `json:"min_ttl,omitempty"`

	mu      sync.Mutex
	records map[string][]libdns.RR
}

func zoneKey(zone string) string {
	return strings.ToLower(strings.TrimSuffix(zone, ".")) + "."
}

// zone returns the records of a zone, the caller must hold the lock.
func (p *Provider) zone(zone string) (string, error) {
	if p.records == nil {
		p.records = map[string][]libdns.RR{}
		for _, configured := range p.Zones {
			p.records[zoneKey(configured)] = nil
		}
	}

	key := zoneKey(zone)
	if _, ok := p.records[key]; !ok {
		return "", fmt.Errorf("%w: %s", ErrZoneNotFound, zone)
	}

	return key, nil
}

// GetRecords returns all records of the zone.
func (p *Provider) GetRecords(_ context.Context, zone string) ([]libdns.Record, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, err := p.zone(zone)
	if err != nil {
		return nil, err
	}

	records := make([]libdns.Record, 0, len(p.records[key]))

	for _, rr := range p.records[key] {
		record, err := rr.Parse()
		if err != nil {
			return nil, fmt.Errorf("parse %s %s: %w", rr.Name, rr.Type, err)
		}

		records = append(records, record)
	}

	return records, nil
}

// AppendRecords adds the records to the zone.
func (p *Provider) AppendRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, err := p.zone(zone)
	if err != nil {
		return nil, err
	}

	appended := make([]libdns.Record, 0, len(records))

	for _, record := range records {
		rr := record.RR()
		for _, existing := range p.records[key] {
			if existing.Name == rr.Name && existing.Type != rr.Type && (existing.Type == "CNAME" || rr.Type == "CNAME") {
				return appended, fmt.Errorf("%w: %s", ErrCNAMEConflict, rr.Name)
			}
		}

		rr.TTL = max(rr.TTL, time.Duration(p.MinTTL)*time.Second)
		p.records[key] = append(p.records[key], rr)
		appended = append(appended, record)
	}

	return appended, nil
}

// SetRecords replaces all records of the given name and type pairs.
func (p *Provider) SetRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, err := p.zone(zone)
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		rr := record.RR()
		p.records[key] = slices.DeleteFunc(p.records[key], func(existing libdns.RR) bool {
			return existing.Name == rr.Name && existing.Type == rr.Type
		})
	}

	for _, record := range records {
		rr := record.RR()
		rr.TTL = max(rr.TTL, time.Duration(p.MinTTL)*time.Second)
		p.records[key] = append(p.records[key], rr)
	}

	return records, nil
}

// DeleteRecords deletes exactly matching records, empty type, TTL and data match any value.
func (p *Provider) DeleteRecords(_ context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, err := p.zone(zone)
	if err != nil {
		return nil, err
	}

	var deleted []libdns.Record

	for _, record := range records {
		rr := record.RR()
		p.records[key] = slices.DeleteFunc(p.records[key], func(existing libdns.RR) bool {
			match := existing.Name == rr.Name &&
				(rr.Type == "" || rr.Type == existing.Type) &&
				(rr.TTL == 0 || rr.TTL == existing.TTL) &&
				(rr.Data == "" || rr.Data == existing.Data)
			if match {
				deleted = append(deleted, existing)
			}

			return match
		})
	}

	return deleted, nil
}

// ListZones returns the configured zones.
func (p *Provider) ListZones(_ context.Context) ([]libdns.Zone, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	_, _ = p.zone(".") // initialize

	zones := make([]libdns.Zone, 0, len(p.records))
	for zone := range p.records {
		zones = append(zones, libdns.Zone{Name: zone})
	}

	slices.SortFunc(zones, func(a, b libdns.Zone) int { return strings.Compare(a.Name, b.Name) })

	return zones, nil
}

// Interface guards.
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
	_ libdns.ZoneLister     = (*Provider)(nil)
)
