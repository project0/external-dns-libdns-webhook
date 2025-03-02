package externaldns

import (
	"slices"
	"strings"
)

type (
	Endpoints []*Endpoint
	Endpoint  struct {
		// The hostname of the DNS record
		DNSName string `json:"dnsName,omitempty"`
		// This is the list of targets that this DNS record points to. So for an A record it will be a list of IP addresses.
		Targets []string `json:"targets,omitempty"`
		// RecordType type of record, e.g. CNAME, A, AAAA, SRV, TXT etc
		RecordType string `json:"recordType,omitempty"`
		// Identifier to distinguish multiple records with the same name and type (e.g. Route53 records with routing policies other than 'simple')
		SetIdentifier string `json:"setIdentifier,omitempty"`
		// TTL for the record
		RecordTTL int64 `json:"recordTTL,omitempty"`
		// Labels stores labels defined for the Endpoint
		// +optional
		Labels map[string]string `json:"labels,omitempty"`
		// ProviderSpecific stores provider specific configThis is the list of changes send by external-dns that need to be applied. There are four lists of endpoints. The create and delete lists are lists of records to create and delete respectively. The updateOld and updateNew lists are paired. For each entry there's the old version of the record and a new version of the record.
		// +optional
		ProviderSpecific ProviderSpecific `json:"providerSpecific,omitempty"`
	}
)

// ProviderSpecificProperty holds the name and value of a configuration which is specific to individual DNS providers.
type ProviderSpecificProperty struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// ProviderSpecific holds configuration which is specific to individual DNS providers.
type ProviderSpecific []ProviderSpecificProperty

// NewEndpoint initialization method to be used to create an endpoint.
func NewEndpoint(dnsName, recordType string, targets ...string) *Endpoint {
	return NewEndpointWithTTL(dnsName, recordType, 0, targets...)
}

// NewEndpointWithTTL initialization method to be used to create an endpoint with a TTL struct.
func NewEndpointWithTTL(dnsName, recordType string, ttl int64, targets ...string) *Endpoint {
	cleanTargets := make([]string, len(targets))
	for idx, target := range targets {
		cleanTargets[idx] = strings.TrimSuffix(target, ".")
	}

	return &Endpoint{
		DNSName:    strings.TrimSuffix(dnsName, "."),
		Targets:    cleanTargets,
		RecordType: recordType,
		Labels:     map[string]string{},
		RecordTTL:  ttl,
	}
}

// WithSetIdentifier applies the given set identifier to the endpoint.
func (e *Endpoint) WithSetIdentifier(setIdentifier string) *Endpoint {
	e.SetIdentifier = setIdentifier

	return e
}

// WithProviderSpecific attaches a key/value pair to the Endpoint and returns the Endpoint.
// This can be used to pass additional data through the stages of ExternalDNS's Endpoint processing.
// The assumption is that most of the time this will be provider specific metadata that doesn't
// warrant its own field on the Endpoint object itself. It differs from Labels in the fact that it's
// not persisted in the Registry but only kept in memory during a single record synchronization.
func (e *Endpoint) WithProviderSpecific(key, value string) *Endpoint {
	e.SetProviderSpecificProperty(key, value)

	return e
}

// GetProviderSpecificProperty returns the value of a ProviderSpecificProperty if the property exists.
func (e *Endpoint) GetProviderSpecificProperty(key string) (string, bool) {
	for _, providerSpecific := range e.ProviderSpecific {
		if providerSpecific.Name == key {
			return providerSpecific.Value, true
		}
	}

	return "", false
}

// SetProviderSpecificProperty sets the value of a ProviderSpecificProperty.
func (e *Endpoint) SetProviderSpecificProperty(key string, value string) {
	for i, providerSpecific := range e.ProviderSpecific {
		if providerSpecific.Name == key {
			e.ProviderSpecific[i] = ProviderSpecificProperty{
				Name:  key,
				Value: value,
			}

			return
		}
	}

	e.ProviderSpecific = append(e.ProviderSpecific, ProviderSpecificProperty{Name: key, Value: value})
}

// DeleteProviderSpecificProperty deletes any ProviderSpecificProperty of the specified name.
func (e *Endpoint) DeleteProviderSpecificProperty(key string) {
	for i, providerSpecific := range e.ProviderSpecific {
		if providerSpecific.Name == key {
			e.ProviderSpecific = slices.Delete(e.ProviderSpecific, i, i+1)

			return
		}
	}
}
