package externaldns_test

import (
	"encoding/json"
	"testing"

	"github.com/project0/external-dns-libdns-webhook/internal/externaldns"
)

func TestDomainFilterJSON(t *testing.T) {
	t.Parallel()

	out, err := json.Marshal(externaldns.DomainFilter{Include: []string{"example.com", "example.org"}})
	if err != nil {
		t.Fatal(err)
	}

	if want := `{"include":["example.com","example.org"]}`; string(out) != want {
		t.Errorf("got %s, want %s", out, want)
	}
}

func TestNewEndpointWithTTLTrailingDots(t *testing.T) {
	t.Parallel()

	tests := []struct {
		recordType string
		target     string
		want       string
	}{
		{externaldns.RecordTypeCNAME, "target.example.com.", "target.example.com"},
		{externaldns.RecordTypeMX, "10 mail.example.com.", "10 mail.example.com"},
		{externaldns.RecordTypeNS, "ns1.example.com.", "ns1.example.com"},
		{externaldns.RecordTypeSRV, "10 5 5060 sip.example.com.", "10 5 5060 sip.example.com."},
		{externaldns.RecordTypeTXT, "ends with dot.", "ends with dot."},
	}

	for _, tc := range tests {
		t.Run(tc.recordType, func(t *testing.T) {
			t.Parallel()

			ep := externaldns.NewEndpointWithTTL("www.example.com.", tc.recordType, 300, tc.target)
			if ep.DNSName != "www.example.com" {
				t.Errorf("dnsName: got %q", ep.DNSName)
			}

			if ep.Targets[0] != tc.want {
				t.Errorf("target: got %q, want %q", ep.Targets[0], tc.want)
			}
		})
	}
}
