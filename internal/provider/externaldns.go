package provider

import (
	"context"
	"fmt"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/project0/external-dns-libdns-webhook/internal/libdnsregistry"
	"github.com/rs/zerolog/log"
)

type WebhookProvider struct {
	// DomainFilter, as zoneList is not implemented in the majority of libdns providers it is a fixed list of zones
	zones              []string
	libdnsProvider     libdnsregistry.Provider
	cachedZonesRecords map[string][]libdns.Record
}

func New(zones []string, libdnsProvider libdnsregistry.Provider) *WebhookProvider {
	return &WebhookProvider{
		zones:              zones,
		libdnsProvider:     libdnsProvider,
		cachedZonesRecords: map[string][]libdns.Record{},
	}
}

func (p *WebhookProvider) Records(ctx context.Context) (externaldns.Endpoints, error) {
	endpoints := []*externaldns.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.zones {
		logger := log.Ctx(ctx).With().Str("zone", zone).Logger()
		logger.Debug().Msg("Retrieving records for zone")

		records, err := p.libdnsProvider.GetRecords(ctx, zone)
		if err != nil {
			logger.Err(err).Msg("Failed to retrieve records for zone")

			return nil, fmt.Errorf("failed to retrieve records for zone: %w", err)
		}

		// as there is no real concurrent sync in progress we can cache between the calls to avoid calling api too many times
		p.cachedZonesRecords[zone] = records

		for _, record := range records {
			endpoint := toExternalDNSEndpoint(record, zone)
			logger.Trace().
				Any("record", record).
				Any("endpoint", endpoint).
				Msg("Record converted to endpoint")

			endpoints = append(endpoints, endpoint)
		}
	}

	return endpoints, nil
}

func (p *WebhookProvider) ApplyChanges(ctx context.Context, changes *externaldns.Changes) error {
	log.Ctx(ctx).Debug().
		Any("changes_create", changes.Create).
		Msg("Convet creation change endpoints to records")

	creates, err := endpointsToLibdnsZoneRecords(changes.Create, p.zones)
	if err != nil {
		return err
	}

	log.Ctx(ctx).Debug().
		Any("changes_delete", changes.Delete).
		Msg("Convert deletion change endpoints to records")

	deletes, err := endpointsToLibdnsZoneRecords(changes.Delete, p.zones)
	if err != nil {
		return err
	}

	log.Ctx(ctx).Debug().
		Any("changes_update_new", changes.UpdateNew).
		Any("changes_update_old", changes.UpdateOld).
		Msg("Convert updates change endpoints to records")

	updates, err := endpointsToLibdnsZoneRecords(changes.UpdateNew, p.zones)
	if err != nil {
		return err
	}

	if len(creates) > 0 {
		for zone, records := range creates {
			log.Ctx(ctx).Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Creating records")

			_, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to create records: %w", err)
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			log.Ctx(ctx).Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Deleting records")

			_, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to delete records: %w", err)
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			log.Ctx(ctx).Info().
				Any("records", records).
				Any("zone", zone).
				Msg("Updating records")

			_, err := p.libdnsProvider.SetRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to update records: %w", err)
			}
		}
	}

	return nil
}

func (p *WebhookProvider) AdjustEndpoints(ctx context.Context, endpoints externaldns.Endpoints) (externaldns.Endpoints, error) {
	for _, endpoint := range endpoints {
		// defaults
		if _, exist := endpoint.GetProviderSpecificProperty(identifierLabelWeight); !exist {
			endpoint.WithProviderSpecific(identifierLabelWeight, "0")
		}

		if _, exist := endpoint.GetProviderSpecificProperty(identifierLabelPriority); !exist {
			endpoint.WithProviderSpecific(identifierLabelPriority, "0")
		}

		// set missing record IDs
		if _, exist := endpoint.GetProviderSpecificProperty(identifierLabelID); !exist {
			endpoint.WithProviderSpecific(identifierLabelID, "")

			recordZones, err := endpointsToLibdnsZoneRecords([]*externaldns.Endpoint{endpoint}, p.zones)
			if err != nil {
				return nil, err
			}
			// there should be only one item in the map
			for zone, records := range recordZones {
				// there will be only one records
				for _, record := range records {
					if cachedRecords, exist := p.cachedZonesRecords[zone]; exist {
						for _, zoneRecord := range cachedRecords {
							if record.Name == zoneRecord.Name && record.Type == zoneRecord.Type {
								log.Ctx(ctx).Debug().
									Any("record", record).
									Any("record_zone", zoneRecord).
									Msg("Set record ID for endpoint")
								endpoint.SetProviderSpecificProperty(identifierLabelID, zoneRecord.ID)
							}
						}
					}
				}
			}
		}
	}

	return endpoints, nil
}

func (p *WebhookProvider) GetInitialization(_ context.Context) externaldns.Initialization {
	return externaldns.Initialization{
		Filters: p.zones,
	}
}
