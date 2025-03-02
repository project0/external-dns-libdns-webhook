package provider

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
	"github.com/rs/zerolog/log"
)

const (
	identifierLabelWeight   = "weight"
	identifierLabelPriority = "priority"
	identifierLabelID       = "record-id"
)

// relativeName returns the name part of a DNS name in a zone.
func relativeName(name, zone string) string {
	return libdns.RelativeName(name, zone)
}

// absoluteName returns the FQDN of a DNS name in a zone.
func absoluteName(name, zone string) string {
	return libdns.AbsoluteName(name, zone)
}

// splitDNSName splits a DNS name into a name and a zone.
func splitDNSName(dnsName string, zones []string) (string, string) {
	name := strings.TrimSuffix(dnsName, ".")

	domain := ""
	resourceRecord := ""
	// sort zones by dot count; make sure subdomains sort earlier
	sort.Slice(zones, func(i, j int) bool {
		return strings.Count(zones[i], ".") > strings.Count(zones[j], ".")
	})

	for _, filter := range zones {
		if strings.HasSuffix(name, "."+filter) {
			domain = filter
			resourceRecord = name[0 : len(name)-len(filter)-1]

			break
		} else if name == filter {
			domain = filter
			resourceRecord = ""
		}
	}

	return resourceRecord, domain
}

func toExternalDNSEndpoint(record libdns.Record, zone string) *externaldns.Endpoint {
	endpoint := externaldns.NewEndpointWithTTL(absoluteName(record.Name, zone), record.Type, int64(record.TTL.Seconds()), record.Value)
	endpoint.WithProviderSpecific(identifierLabelWeight, strconv.FormatUint(uint64(record.Weight), 10))
	endpoint.WithProviderSpecific(identifierLabelPriority, strconv.FormatUint(uint64(record.Priority), 10))
	endpoint.WithProviderSpecific(identifierLabelID, record.ID)

	return endpoint
}

func toLibdnsRecord(endpoint *externaldns.Endpoint, zone string) libdns.Record {
	record := libdns.Record{
		Type:  endpoint.RecordType,
		Name:  relativeName(endpoint.DNSName, zone),
		Value: endpoint.Targets[0],
		TTL:   time.Duration(endpoint.RecordTTL) * time.Second,
	}

	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelID); ok {
		record.ID = prop
	}

	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelWeight); ok {
		weight, err := strconv.ParseUint(prop, 10, 64)
		if err != nil {
			log.Err(err).Str("weigth", prop).Msg("Failed to parse weight")
		} else {
			record.Weight = uint(weight)
		}
	}

	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelPriority); ok {
		prio, err := strconv.ParseUint(prop, 10, 64)
		if err != nil {
			log.Err(err).Str("priority", prop).Msg("Failed to parse priority")
		} else {
			record.Priority = uint(prio)
		}
	}

	return record
}

func endpointsToLibdnsZoneRecords(endpoints []*externaldns.Endpoint, zones []string) (map[string][]libdns.Record, error) {
	zoneRecords := map[string][]libdns.Record{}

	for _, endpoint := range endpoints {
		_, zone := splitDNSName(endpoint.DNSName, zones)
		if zone == "" {
			return nil, fmt.Errorf("no matching zone found for %s", endpoint.DNSName)
		}

		record := toLibdnsRecord(endpoint, zone)

		if _, ok := zoneRecords[zone]; !ok {
			zoneRecords[zone] = []libdns.Record{}
		}

		zoneRecords[zone] = append(zoneRecords[zone], record)
	}

	return zoneRecords, nil
}
