package webhook

import (
	"context"
	"fmt"

	ilibdns "github.com/project0/external-dns-webhook-libdns/internal/libdns"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/plan"
	"sigs.k8s.io/external-dns/provider"
)

type WebhookProvider struct {
	provider.BaseProvider
	// DomainFilter, as zoneList is not implemented in the majority of libdns providers it is a fixed list of zones
	DomainFilter   endpoint.DomainFilter
	libdnsProvider ilibdns.LibdnsProvider
	log            *log.Logger
}

func NewWebhookProvider(zones []string, libdnsProvider ilibdns.LibdnsProvider, logger *log.Logger) *WebhookProvider {
	return &WebhookProvider{
		DomainFilter:   endpoint.NewDomainFilter(zones),
		libdnsProvider: libdnsProvider,
		log:            logger,
	}
}

func (p *WebhookProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := []*endpoint.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.DomainFilter.Filters {
		records, err := p.libdnsProvider.GetRecords(ctx, zone)
		if err != nil {
			log.WithContext(ctx).WithError(err).WithField("zone", zone).Errorf("failed to get records for zone %s", zone)
			return nil, err
		}
		for _, record := range records {
			ep, err := toExternalDNSEndpoint(record, zone)
			if err != nil {
				log.WithContext(ctx).WithError(err).WithField("zone", zone).WithField("record", record).Errorf("failed to convert record to endpoints")
				return nil, err
			}
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints, nil
}

func (p *WebhookProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	creates, err := endpointsToLibdnsZoneRecords(changes.Create, p.DomainFilter.Filters)
	if err != nil {
		return err
	}
	deletes, err := endpointsToLibdnsZoneRecords(changes.Delete, p.DomainFilter.Filters)
	if err != nil {
		return err
	}
	updates, err := endpointsToLibdnsZoneRecords(changes.UpdateNew, p.DomainFilter.Filters)
	if err != nil {
		return err
	}

	if len(creates) > 0 {
		for zone, records := range creates {
			_, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to create records: %w", err)
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			_, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to delete records: %w", err)
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			_, err := p.libdnsProvider.SetRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to update records: %w", err)
			}
		}
	}

	return nil
}

func (p *WebhookProvider) AdjustEndpoints(endpoints []*endpoint.Endpoint) ([]*endpoint.Endpoint, error) {
	// not implemented
	return endpoints, nil
}

func (p *WebhookProvider) GetDomainFilter() endpoint.DomainFilterInterface {
	return p.DomainFilter
}
