package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/rs/zerolog/log"
)

type WebhookProvider struct {
	// zones managed by the provider, normalized and sorted by most specific first
	zones          []string
	libdnsProvider libdnsregistry.Provider
}

func New(zones []string, libdnsProvider libdnsregistry.Provider) *WebhookProvider {
	return &WebhookProvider{
		zones:          normalizeZones(zones),
		libdnsProvider: libdnsProvider,
	}
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

func (p *WebhookProvider) ApplyChanges(ctx context.Context, changes *externaldns.Changes) error {
	creates, err := p.zoneRecords(changes.Create)
	if err != nil {
		return err
	}

	deletes, err := p.zoneRecords(changes.Delete)
	if err != nil {
		return err
	}

	updates, err := p.zoneRecords(changes.UpdateNew)
	if err != nil {
		return err
	}

	var errs []error

	for zone, records := range creates {
		log.Ctx(ctx).Info().Str("zone", zone).Any("records", records).Msg("Creating records")

		if _, err := p.libdnsProvider.AppendRecords(ctx, fqdn(zone), records); err != nil {
			errs = append(errs, fmt.Errorf("failed to create records in zone %s: %w", zone, err))
		}
	}

	for zone, records := range deletes {
		log.Ctx(ctx).Info().Str("zone", zone).Any("records", records).Msg("Deleting records")

		if _, err := p.libdnsProvider.DeleteRecords(ctx, fqdn(zone), records); err != nil {
			errs = append(errs, fmt.Errorf("failed to delete records in zone %s: %w", zone, err))
		}
	}

	for zone, records := range updates {
		log.Ctx(ctx).Info().Str("zone", zone).Any("records", records).Msg("Updating records")

		if _, err := p.libdnsProvider.SetRecords(ctx, fqdn(zone), records); err != nil {
			errs = append(errs, fmt.Errorf("failed to update records in zone %s: %w", zone, err))
		}
	}

	return errors.Join(errs...)
}

// zoneRecords converts endpoints to libdns records grouped by zone.
func (p *WebhookProvider) zoneRecords(endpoints externaldns.Endpoints) (map[string][]libdns.Record, error) {
	zoneRecords := map[string][]libdns.Record{}

	for _, endpoint := range endpoints {
		zone, name, ok := splitDNSName(endpoint.DNSName, p.zones)
		if !ok {
			return nil, fmt.Errorf("no matching zone found for %s", endpoint.DNSName)
		}

		records, err := endpointToRecords(endpoint, name)
		if err != nil {
			return nil, err
		}

		zoneRecords[zone] = append(zoneRecords[zone], records...)
	}

	return zoneRecords, nil
}

func (p *WebhookProvider) AdjustEndpoints(_ context.Context, endpoints externaldns.Endpoints) (externaldns.Endpoints, error) {
	return endpoints, nil
}

func (p *WebhookProvider) GetDomainFilter(_ context.Context) externaldns.DomainFilter {
	return externaldns.DomainFilter{
		Include: p.zones,
	}
}
