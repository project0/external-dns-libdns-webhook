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
}

func NewWebhookProvider(zones []string, libdnsProvider ilibdns.LibdnsProvider) *WebhookProvider {
	return &WebhookProvider{
		DomainFilter:   endpoint.NewDomainFilter(zones),
		libdnsProvider: libdnsProvider,
	}
}

func (p *WebhookProvider) Records(ctx context.Context) ([]*endpoint.Endpoint, error) {
	endpoints := []*endpoint.Endpoint{}

	// return all records for configured zones
	for _, zone := range p.DomainFilter.Filters {
		log.WithContext(ctx).WithField("zone", zone).Debug("Getting records for zone")
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
			log.Info(ep.Targets[0])
			log.WithContext(ctx).WithField("endpoint", ep).WithField("record", record).Trace("Converted record to endpoint")
			endpoints = append(endpoints, ep)
		}
	}
	return endpoints, nil
}

func (p *WebhookProvider) ApplyChanges(ctx context.Context, changes *plan.Changes) error {
	log.
		WithContext(ctx).
		WithField("changes_create", changes.Create).
		Debug("Endpoint creations")
	creates, err := endpointsToLibdnsZoneRecords(changes.Create, p.DomainFilter.Filters)
	if err != nil {
		return err
	}

	log.
		WithContext(ctx).
		WithField("changes_delete", changes.Delete).
		Debug("Endpoint deletions")
	deletes, err := endpointsToLibdnsZoneRecords(changes.Delete, p.DomainFilter.Filters)
	if err != nil {
		return err
	}

	log.
		WithContext(ctx).
		WithField("changes_update_new", changes.UpdateNew).
		WithField("changes_update_new", changes.UpdateOld).
		Debug("Endpoint updates")
	updates, err := endpointsToLibdnsZoneRecords(changes.UpdateNew, p.DomainFilter.Filters)
	if err != nil {
		return err
	}

	if len(creates) > 0 {
		for zone, records := range creates {
			log.
				WithContext(ctx).
				WithField("records", records).
				WithField("zone", zone).
				Debug("Creating records")
			_, err := p.libdnsProvider.AppendRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to create records: %w", err)
			}
		}
	}

	if len(deletes) > 0 {
		for zone, records := range deletes {
			log.WithContext(ctx).
				WithField("records", records).
				WithField("zone", zone).
				Debug("Deleting records")
			_, err := p.libdnsProvider.DeleteRecords(ctx, zone, records)
			if err != nil {
				return fmt.Errorf("failed to delete records: %w", err)
			}
		}
	}

	if len(updates) > 0 {
		for zone, records := range updates {
			log.
				WithContext(ctx).
				WithField("records", records).
				WithField("zone", zone).
				Debug("Updating records")
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
