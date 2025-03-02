package externaldns

type Initialization struct {
	Filters DomainFilters `json:"filters,omitempty"`
}

// external-dns will only create DNS records for host names
// (specified in ingress objects and services with the external-dns annotation)
// related to zones that match filters.
// They can be set in external-dns deployment manifest.
type DomainFilters []string
