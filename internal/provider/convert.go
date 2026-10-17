package provider

import (
	"cmp"
	"errors"
	"fmt"
	"net/netip"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/libdns/libdns"
	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
)

const (
	mxFields  = 2
	srvFields = 4
	nullHost  = "."
)

var errUnsupportedRecordType = errors.New("unsupported record type")

// supportedRecordTypes are the record types external-dns and libdns have in common.
var supportedRecordTypes = []string{
	externaldns.RecordTypeA,
	externaldns.RecordTypeAAAA,
	externaldns.RecordTypeCNAME,
	externaldns.RecordTypeTXT,
	externaldns.RecordTypeMX,
	externaldns.RecordTypeNS,
	externaldns.RecordTypeSRV,
}

func isSupportedRecordType(recordType string) bool {
	return slices.Contains(supportedRecordTypes, recordType)
}

// normalizeName returns a lower case DNS name without trailing dot.
func normalizeName(name string) string {
	return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), "."))
}

// fqdn returns the zone name in the fully qualified form libdns expects.
func fqdn(zone string) string {
	return zone + "."
}

// normalizeZones normalizes and deduplicates zones, sorted by most specific (longest) first.
func normalizeZones(zones []string) []string {
	normalized := make([]string, 0, len(zones))

	for _, zone := range zones {
		if zone = normalizeName(zone); zone != "" && !slices.Contains(normalized, zone) {
			normalized = append(normalized, zone)
		}
	}

	slices.SortFunc(normalized, func(a, b string) int {
		return cmp.Or(
			cmp.Compare(strings.Count(b, "."), strings.Count(a, ".")),
			cmp.Compare(len(b), len(a)),
			strings.Compare(a, b),
		)
	})

	return normalized
}

// splitDNSName returns the most specific zone and the libdns relative record name of a DNS name.
func splitDNSName(dnsName string, zones []string) (string, string, bool) {
	name := normalizeName(dnsName)

	for _, zone := range zones {
		if name == zone {
			return zone, "@", true
		}

		if rel, ok := strings.CutSuffix(name, "."+zone); ok {
			return zone, rel, true
		}
	}

	return "", "", false
}

// absoluteRecordName returns the normalized FQDN of a record name returned by a libdns provider.
// Names should be relative, but some providers return them fully qualified.
func absoluteRecordName(name, zone string) string {
	name = normalizeName(name)

	switch {
	case name == "" || name == "@":
		return zone
	case name == zone || strings.HasSuffix(name, "."+zone):
		return name
	default:
		return name + "." + zone
	}
}

func canonicalHost(host string) string {
	if host == nullHost {
		return host
	}

	return normalizeName(host)
}

func parseUint16(fields []string) ([]uint16, error) {
	values := make([]uint16, 0, len(fields))

	for _, field := range fields {
		value, err := strconv.ParseUint(field, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %w", field, err)
		}

		values = append(values, uint16(value))
	}

	return values, nil
}

// canonicalTarget returns the target in the form returned to external-dns, regardless of the input spelling.
// It must produce the same output for desired endpoints (AdjustEndpoints) and records read from the provider,
// otherwise external-dns plans updates on every sync.
//
//   - A/AAAA: normalized IP address
//   - CNAME/NS: lower case host without trailing dot
//   - MX: "<preference> <host>", host without trailing dot ("0 ." for null MX)
//   - SRV: "<priority> <weight> <port> <host>.", external-dns requires the trailing dot
//   - TXT: unquoted text
func canonicalTarget(recordType, target string) (string, error) {
	target = strings.TrimSpace(target)

	switch recordType {
	case externaldns.RecordTypeA, externaldns.RecordTypeAAAA:
		addr, err := netip.ParseAddr(target)
		if err != nil {
			return "", fmt.Errorf("invalid %s target %q: %w", recordType, target, err)
		}

		return addr.String(), nil

	case externaldns.RecordTypeCNAME, externaldns.RecordTypeNS:
		host := normalizeName(target)
		if host == "" {
			return "", fmt.Errorf("empty %s target", recordType)
		}

		return host, nil

	case externaldns.RecordTypeMX:
		fields := strings.Fields(target)
		if len(fields) != mxFields {
			return "", fmt.Errorf("invalid MX target %q, expected '<preference> <host>'", target)
		}

		values, err := parseUint16(fields[:1])
		if err != nil {
			return "", fmt.Errorf("invalid MX target %q: %w", target, err)
		}

		return fmt.Sprintf("%d %s", values[0], canonicalHost(fields[1])), nil

	case externaldns.RecordTypeSRV:
		fields := strings.Fields(target)
		if len(fields) != srvFields {
			return "", fmt.Errorf("invalid SRV target %q, expected '<priority> <weight> <port> <host>'", target)
		}

		values, err := parseUint16(fields[:3])
		if err != nil {
			return "", fmt.Errorf("invalid SRV target %q: %w", target, err)
		}

		host := canonicalHost(fields[3])
		if host != nullHost {
			host += "."
		}

		return fmt.Sprintf("%d %d %d %s", values[0], values[1], values[2], host), nil

	case externaldns.RecordTypeTXT:
		if len(target) >= 2 && strings.HasPrefix(target, `"`) && strings.HasSuffix(target, `"`) {
			return target[1 : len(target)-1], nil
		}

		return target, nil

	default:
		return "", fmt.Errorf("%w: %s", errUnsupportedRecordType, recordType)
	}
}

// recordData converts a canonical target to libdns RR data. Hostnames are fully qualified with trailing dot.
func recordData(recordType, target string) string {
	withDot := func(host string) string {
		if host == nullHost {
			return host
		}

		return host + "."
	}

	switch recordType {
	case externaldns.RecordTypeCNAME, externaldns.RecordTypeNS:
		return withDot(target)
	case externaldns.RecordTypeMX:
		preference, host, _ := strings.Cut(target, " ")

		return preference + " " + withDot(host)
	default:
		return target
	}
}

// recordTarget returns the canonical target of a libdns record.
func recordTarget(record libdns.Record) (string, error) {
	rr := record.RR()
	if rr.Type == externaldns.RecordTypeTXT {
		// libdns TXT data is always unquoted, quotes are part of the text
		return rr.Data, nil
	}

	return canonicalTarget(rr.Type, rr.Data)
}

// recordsToEndpoints groups libdns records of a zone by name and type into endpoints with multiple targets.
func recordsToEndpoints(records []libdns.Record, zone string) (externaldns.Endpoints, error) {
	type key struct{ name, recordType string }

	grouped := map[key]*externaldns.Endpoint{}

	var errs []error

	for _, record := range records {
		rr := record.RR()
		if !isSupportedRecordType(rr.Type) {
			continue
		}

		target, err := recordTarget(record)
		if err != nil {
			errs = append(errs, fmt.Errorf("record %s %s in zone %s: %w", rr.Name, rr.Type, zone, err))

			continue
		}

		recordKey := key{absoluteRecordName(rr.Name, zone), rr.Type}

		endpoint, exists := grouped[recordKey]
		if !exists {
			endpoint = &externaldns.Endpoint{
				DNSName:    recordKey.name,
				RecordType: recordKey.recordType,
				RecordTTL:  int64(rr.TTL.Seconds()),
				Labels:     map[string]string{},
			}
			grouped[recordKey] = endpoint
		}

		if !slices.Contains(endpoint.Targets, target) {
			endpoint.Targets = append(endpoint.Targets, target)
		}
	}

	endpoints := make(externaldns.Endpoints, 0, len(grouped))
	for _, endpoint := range grouped {
		slices.Sort(endpoint.Targets)
		endpoints = append(endpoints, endpoint)
	}

	slices.SortFunc(endpoints, func(a, b *externaldns.Endpoint) int {
		return cmp.Or(strings.Compare(a.DNSName, b.DNSName), strings.Compare(a.RecordType, b.RecordType))
	})

	return endpoints, errors.Join(errs...)
}

// endpointToRecords converts an endpoint to typed libdns records, one per target.
// Typed records are required as providers type switch on them (e.g. libdns.MX).
func endpointToRecords(endpoint *externaldns.Endpoint, name string) ([]libdns.Record, error) {
	records := make([]libdns.Record, 0, len(endpoint.Targets))

	for _, target := range endpoint.Targets {
		canonical, err := canonicalTarget(endpoint.RecordType, target)
		if err != nil {
			return nil, fmt.Errorf("endpoint %s %s: %w", endpoint.DNSName, endpoint.RecordType, err)
		}

		record, err := libdns.RR{
			Name: name,
			Type: endpoint.RecordType,
			TTL:  time.Duration(endpoint.RecordTTL) * time.Second,
			Data: recordData(endpoint.RecordType, canonical),
		}.Parse()
		if err != nil {
			return nil, fmt.Errorf("endpoint %s %s: %w", endpoint.DNSName, endpoint.RecordType, err)
		}

		records = append(records, record)
	}

	return records, nil
}
