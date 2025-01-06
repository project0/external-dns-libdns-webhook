package webhook

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
	log "github.com/sirupsen/logrus"
	"sigs.k8s.io/external-dns/endpoint"
)

const (
	identifierLabelWeight   = "weight"
	identifierLabelPriority = "priority"
)

// small wrapper to keep track of the zone of a record
type libdnsZoneRecord struct {
	record libdns.Record
	zone   string
}

// relativeName returns the name part of a DNS name in a zone.
func relativeName(name, zone string) string {
	return libdns.RelativeName(name, zone)
}

// absoluteName returns the FQDN of a DNS name in a zone.
func absoluteName(name, zone string) string {
	return libdns.AbsoluteName(name, zone)
}

// splitDNSName splits a DNS name into a name and a zone.
func splitDNSName(dnsName string, zones []string) (rr string, domain string) {
	name := strings.TrimSuffix(dnsName, ".")

	// sort zones by dot count; make sure subdomains sort earlier
	sort.Slice(zones, func(i, j int) bool {
		return strings.Count(zones[i], ".") > strings.Count(zones[j], ".")
	})
	for _, filter := range zones {
		if strings.HasSuffix(name, "."+filter) {
			rr = name[0 : len(name)-len(filter)-1]
			domain = filter
			break
		} else if name == filter {
			domain = filter
			rr = ""
		}
	}
	return rr, domain
}

func toExternalDNSEndpoint(record libdns.Record, zone string) (*endpoint.Endpoint, error) {
	endpoint := &endpoint.Endpoint{
		DNSName:          absoluteName(record.Name, zone),
		Targets:          []string{record.Value},
		RecordType:       record.Type,
		RecordTTL:        endpoint.TTL(record.TTL.Seconds()),
		SetIdentifier:    record.ID,
		ProviderSpecific: []endpoint.ProviderSpecificProperty{},
	}
	endpoint.WithProviderSpecific(identifierLabelWeight, string(record.Weight))
	endpoint.WithProviderSpecific(identifierLabelPriority, string(record.Priority))
	return endpoint, nil
}

func toLibdnsRecord(endpoint *endpoint.Endpoint, zone string) (libdns.Record, error) {
	record := libdns.Record{
		ID:    endpoint.SetIdentifier,
		Type:  endpoint.RecordType,
		Name:  relativeName(endpoint.DNSName, zone),
		Value: endpoint.Targets[0],
		TTL:   time.Duration(int64(endpoint.RecordTTL)) * time.Second,
	}

	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelWeight); ok {
		weight, err := strconv.ParseUint(prop, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse weight: %v", err)
		} else {
			record.Weight = uint(weight)
		}
	}

	if prop, ok := endpoint.GetProviderSpecificProperty(identifierLabelPriority); ok {
		prio, err := strconv.ParseUint(prop, 10, 64)
		if err != nil {
			log.Errorf("Failed to parse priority: %v", err)
		} else {
			record.Priority = uint(prio)
		}
	}

	return record, nil
}

func endpointsToLibdnsZoneRecords(endpoints []*endpoint.Endpoint, zones []string) (map[string][]libdns.Record, error) {
	zoneRecords := map[string][]libdns.Record{}

	for _, endpoint := range endpoints {
		_, zone := splitDNSName(endpoint.DNSName, zones)
		if zone == "" {
			return nil, fmt.Errorf("no matching zone found for %s", endpoint.DNSName)
		}

		record, err := toLibdnsRecord(endpoint, relativeName(endpoint.DNSName, zone))
		if err != nil {
			return nil, fmt.Errorf("failed to convert endpoint to libdns record: %w", err)
		}
		if _, ok := zoneRecords[zone]; !ok {
			zoneRecords[zone] = []libdns.Record{}
		}
		zoneRecords[zone] = append(zoneRecords[zone], record)

	}
	return zoneRecords, nil
}
