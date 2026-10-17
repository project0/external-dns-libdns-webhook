package externaldns

// DomainFilter is returned on webhook negotiation. external-dns will only send endpoints
// for domains matching this filter. The JSON encoding matches the one of upstream
// external-dns (endpoint.DomainFilter): either domain lists or regular expressions.
type DomainFilter struct {
	Include      []string `json:"include,omitempty"`
	Exclude      []string `json:"exclude,omitempty"`
	RegexInclude string   `json:"regexInclude,omitempty"`
	RegexExclude string   `json:"regexExclude,omitempty"`
}
