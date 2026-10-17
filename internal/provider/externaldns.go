package provider

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Options of the webhook provider.
type Options struct {
	// DefaultTTL is used for endpoints without TTL.
	DefaultTTL time.Duration
	// MinTTL raises lower TTLs, set it to the minimum the DNS provider accepts to avoid endless updates.
	MinTTL time.Duration
}

type WebhookProvider struct {
	// zones managed by the provider, normalized and sorted by most specific first
	zones          []string
	libdnsProvider libdnsregistry.Provider
	defaultTTL     int64
	minTTL         int64
}

// New creates the webhook provider. Without configured zones, the zones are discovered once
// if the libdns provider implements libdns.ZoneLister.
func New(ctx context.Context, zones []string, libdnsProvider libdnsregistry.Provider, opts Options) (*WebhookProvider, error) {
	if len(zones) == 0 {
		lister, ok := libdnsProvider.(libdns.ZoneLister)
		if !ok {
			return nil, errors.New("no zones configured and the provider does not support listing zones")
		}

		listed, err := lister.ListZones(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to list zones: %w", err)
		}

		for _, zone := range listed {
			zones = append(zones, zone.Name)
		}

		log.Ctx(ctx).Info().Strs("zones", zones).Msg("Discovered zones")
	}

	normalized := normalizeZones(zones)
	if len(normalized) == 0 {
		return nil, errors.New("no zones to manage")
	}

	return &WebhookProvider{
		zones:          normalized,
		libdnsProvider: libdnsProvider,
		defaultTTL:     int64(opts.DefaultTTL.Seconds()),
		minTTL:         int64(opts.MinTTL.Seconds()),
	}, nil
}

// Records returns all supported records of the managed zones, grouped by name and type.
func (p *WebhookProvider) Records(ctx context.Context) (externaldns.Endpoints, error) {
	endpoints := externaldns.Endpoints{}

	for _, zone := range p.zones {
		logger := log.Ctx(ctx).With().Str("zone", zone).Logger()
		logger.Debug().Msg("Retrieving records for zone")

		// fail entirely, otherwise external-dns would plan to recreate records of the missing zone
		records, err := p.libdnsProvider.GetRecords(ctx, fqdn(zone))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve records for zone %s: %w", zone, err)
		}

		zoneEndpoints, err := recordsToEndpoints(records, zone)
		if err != nil {
			// a single malformed record must not break the whole sync
			logger.Warn().Err(err).Msg("Skipped records that cannot be converted")
		}

		logger.Trace().Any("records", records).Any("endpoints", zoneEndpoints).Msg("Records converted to endpoints")

		endpoints = append(endpoints, zoneEndpoints...)
	}

	return endpoints, nil
}

// AdjustEndpoints canonicalizes the desired endpoints to exactly match what Records returns once applied,
// any difference makes external-dns plan updates on every sync.
func (p *WebhookProvider) AdjustEndpoints(ctx context.Context, endpoints externaldns.Endpoints) (externaldns.Endpoints, error) {
	adjusted := make(externaldns.Endpoints, 0, len(endpoints))

	for _, endpoint := range endpoints {
		logger := log.Ctx(ctx).With().Str("dns_name", endpoint.DNSName).Str("record_type", endpoint.RecordType).Logger()

		if reason := p.unsupported(endpoint); reason != "" {
			logger.Warn().Msgf("Ignoring endpoint: %s", reason)

			continue
		}

		targets := make([]string, 0, len(endpoint.Targets))

		for _, target := range endpoint.Targets {
			canonical, err := canonicalTarget(endpoint.RecordType, target)
			if err != nil {
				logger.Warn().Err(err).Msg("Ignoring invalid target")

				continue
			}

			if !slices.Contains(targets, canonical) {
				targets = append(targets, canonical)
			}
		}

		if len(targets) == 0 {
			logger.Warn().Msg("Ignoring endpoint without valid targets")

			continue
		}

		slices.Sort(targets)

		endpoint.DNSName = normalizeName(endpoint.DNSName)
		endpoint.Targets = targets
		endpoint.RecordTTL = p.ttl(endpoint.RecordTTL)
		// provider specific properties are not supported and never returned by Records
		endpoint.ProviderSpecific = nil

		adjusted = append(adjusted, endpoint)
	}

	return adjusted, nil
}

// unsupported returns the reason why an endpoint cannot be managed, or an empty string.
func (p *WebhookProvider) unsupported(endpoint *externaldns.Endpoint) string {
	switch {
	case endpoint.SetIdentifier != "":
		return "set identifiers (routing policies) are not supported"
	case !isSupportedRecordType(endpoint.RecordType):
		return "record type is not supported"
	default:
		if _, _, ok := splitDNSName(endpoint.DNSName, p.zones); !ok {
			return "no matching zone"
		}

		return ""
	}
}

func (p *WebhookProvider) ttl(ttl int64) int64 {
	if ttl <= 0 {
		ttl = p.defaultTTL
	}

	return max(ttl, p.minTTL)
}

// recordKey identifies a record set within a zone.
type recordKey struct {
	name       string
	recordType string
}

// zoneChanges are the records to delete and append within a zone.
type zoneChanges struct {
	// desired target sets (and their TTL) by record set, replace=true removes targets not desired anymore
	desired map[recordKey]desiredSet
	// deletes contains targets to delete by record set
	deletes map[recordKey][]string
}

type desiredSet struct {
	name    string
	ttl     int64
	targets []string
	replace bool
}

// ApplyChanges reconciles the changes against the current records of each zone. Instead of SetRecords,
// whose implementations differ a lot between providers, only DeleteRecords and AppendRecords are used:
// Deletes run first (required e.g. when replacing an A by a CNAME record), then the missing records are appended.
// Deleted records are the ones returned by the provider, so exact matching works regardless of its data format.
func (p *WebhookProvider) ApplyChanges(ctx context.Context, changes *externaldns.Changes) error {
	zones := map[string]*zoneChanges{}

	var errs []error

	forZone := func(endpoint *externaldns.Endpoint) (*zoneChanges, recordKey, string, bool) {
		if _, _, ok := splitDNSName(endpoint.DNSName, p.zones); !ok {
			// like other providers, ignore records outside of the managed zones instead of failing all changes
			event := log.Ctx(ctx).Warn().Str("dns_name", endpoint.DNSName).Str("record_type", endpoint.RecordType)
			if endpoint.RecordType == externaldns.RecordTypeTXT {
				event = event.Str("hint", "registry TXT records of zone apex records need --txt-prefix containing %{record_type}, e.g. edns-%{record_type}.")
			}

			event.Msg("Ignoring endpoint outside of the managed zones")

			return nil, recordKey{}, "", false
		}

		if reason := p.unsupported(endpoint); reason != "" {
			errs = append(errs, fmt.Errorf("endpoint %s %s: %s", endpoint.DNSName, endpoint.RecordType, reason))

			return nil, recordKey{}, "", false
		}

		zone, name, _ := splitDNSName(endpoint.DNSName, p.zones)
		if zones[zone] == nil {
			zones[zone] = &zoneChanges{desired: map[recordKey]desiredSet{}, deletes: map[recordKey][]string{}}
		}

		return zones[zone], recordKey{absoluteRecordName(name, zone), endpoint.RecordType}, name, true
	}

	canonicalTargets := func(endpoint *externaldns.Endpoint) ([]string, bool) {
		targets := make([]string, 0, len(endpoint.Targets))

		for _, target := range endpoint.Targets {
			canonical, err := canonicalTarget(endpoint.RecordType, target)
			if err != nil {
				errs = append(errs, fmt.Errorf("endpoint %s %s: %w", endpoint.DNSName, endpoint.RecordType, err))

				return nil, false
			}

			targets = append(targets, canonical)
		}

		return targets, true
	}

	for _, endpoint := range changes.Delete {
		if changes, key, _, ok := forZone(endpoint); ok {
			if targets, ok := canonicalTargets(endpoint); ok {
				changes.deletes[key] = append(changes.deletes[key], targets...)
			}
		}
	}

	desire := func(endpoints externaldns.Endpoints, replace bool) {
		for _, endpoint := range endpoints {
			if changes, key, name, ok := forZone(endpoint); ok {
				if targets, ok := canonicalTargets(endpoint); ok {
					changes.desired[key] = desiredSet{name: name, ttl: p.ttl(endpoint.RecordTTL), targets: targets, replace: replace}
				}
			}
		}
	}
	// updates replace the whole record set, the old state is taken from the provider instead of UpdateOld
	desire(changes.UpdateNew, true)
	desire(changes.Create, false)

	for zone, changes := range zones {
		if err := p.applyZone(ctx, zone, changes); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (p *WebhookProvider) applyZone(ctx context.Context, zone string, changes *zoneChanges) error {
	logger := log.Ctx(ctx).With().Str("zone", zone).Logger()

	records, err := p.libdnsProvider.GetRecords(ctx, fqdn(zone))
	if err != nil {
		return fmt.Errorf("failed to retrieve records for zone %s: %w", zone, err)
	}

	type currentRecord struct {
		record libdns.Record
		target string
		ttl    int64
	}

	current := map[recordKey][]currentRecord{}

	for _, record := range records {
		rr := record.RR()
		if !isSupportedRecordType(rr.Type) {
			continue
		}

		target, err := recordTarget(record)
		if err != nil {
			continue
		}

		key := recordKey{absoluteRecordName(rr.Name, zone), rr.Type}
		current[key] = append(current[key], currentRecord{record: record, target: target, ttl: int64(rr.TTL.Seconds())})
	}

	var deletes, appends []libdns.Record

	for key, targets := range changes.deletes {
		for _, cur := range current[key] {
			if slices.Contains(targets, cur.target) {
				deletes = append(deletes, cur.record)
			}
		}
	}

	var errs []error

	for key, set := range changes.desired {
		existing := map[string]bool{}

		for _, cur := range current[key] {
			wanted := slices.Contains(set.targets, cur.target)

			switch {
			case wanted && cur.ttl == set.ttl && !existing[cur.target]:
				existing[cur.target] = true
			case wanted || set.replace:
				// TTL changed, duplicate or not desired anymore
				deletes = append(deletes, cur.record)
			}
		}

		missing := &externaldns.Endpoint{DNSName: key.name, RecordType: key.recordType, RecordTTL: set.ttl}

		for _, target := range set.targets {
			if !existing[target] && !slices.Contains(missing.Targets, target) {
				missing.Targets = append(missing.Targets, target)
			}
		}

		toAppend, err := endpointToRecords(missing, set.name)
		if err != nil {
			errs = append(errs, err)

			continue
		}

		appends = append(appends, toAppend...)
	}

	if len(deletes) > 0 {
		logChanges(logger.Info(), deletes).Msg("Deleting records")

		if _, err := p.libdnsProvider.DeleteRecords(ctx, fqdn(zone), deletes); err != nil {
			// appending could conflict with records not deleted
			return errors.Join(append(errs, fmt.Errorf("failed to delete records in zone %s: %w", zone, err))...)
		}
	}

	if len(appends) > 0 {
		logChanges(logger.Info(), appends).Msg("Creating records")

		if _, err := p.libdnsProvider.AppendRecords(ctx, fqdn(zone), appends); err != nil {
			errs = append(errs, fmt.Errorf("failed to create records in zone %s: %w", zone, err))
		}
	}

	return errors.Join(errs...)
}

func logChanges(event *zerolog.Event, records []libdns.Record) *zerolog.Event {
	rrs := make([]libdns.RR, 0, len(records))
	for _, record := range records {
		rrs = append(rrs, record.RR())
	}

	return event.Any("records", rrs)
}

func (p *WebhookProvider) GetDomainFilter(_ context.Context) externaldns.DomainFilter {
	return externaldns.DomainFilter{
		Include: p.zones,
	}
}
