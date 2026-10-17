package provider_test

import (
	"slices"
	"testing"

	"github.com/project0/external-dns-libdns-webhook/internal/provider"
)

func TestCanonicalTarget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		recordType string
		target     string
		want       string
		wantErr    bool
	}{
		{recordType: "A", target: "192.0.2.1", want: "192.0.2.1"},
		{recordType: "A", target: "not-an-ip", wantErr: true},
		{recordType: "AAAA", target: "2001:DB8:0:0::1", want: "2001:db8::1"},
		{recordType: "CNAME", target: "Target.Example.com.", want: "target.example.com"},
		{recordType: "NS", target: "ns1.example.net", want: "ns1.example.net"},
		{recordType: "NS", target: "", wantErr: true},
		{recordType: "MX", target: "10 Mail.example.com.", want: "10 mail.example.com"},
		{recordType: "MX", target: "10   mail.example.com", want: "10 mail.example.com"},
		{recordType: "MX", target: "0 .", want: "0 ."},
		{recordType: "MX", target: "mail.example.com", wantErr: true},
		{recordType: "MX", target: "70000 mail.example.com", wantErr: true},
		{recordType: "SRV", target: "10 5 5060 sip.example.com", want: "10 5 5060 sip.example.com."},
		{recordType: "SRV", target: "10 5 5060 sip.example.com.", want: "10 5 5060 sip.example.com."},
		{recordType: "SRV", target: "0 0 0 .", want: "0 0 0 ."},
		{recordType: "SRV", target: "10 5 sip.example.com", wantErr: true},
		{recordType: "TXT", target: `"heritage=external-dns,external-dns/owner=default"`, want: "heritage=external-dns,external-dns/owner=default"},
		{recordType: "TXT", target: "v=spf1 -all", want: "v=spf1 -all"},
		{recordType: "PTR", target: "host.example.com", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.recordType+" "+tc.target, func(t *testing.T) {
			t.Parallel()

			got, err := provider.CanonicalTarget(tc.recordType, tc.target)
			if tc.wantErr {
				if err == nil {
					t.Errorf("expected error, got %q", got)
				}

				return
			}

			if err != nil {
				t.Fatal(err)
			}

			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestSplitDNSName(t *testing.T) {
	t.Parallel()

	zones := provider.NormalizeZones([]string{"example.com.", "sub.Example.com", "example.org", "example.com"})
	if want := []string{"sub.example.com", "example.com", "example.org"}; !slices.Equal(zones, want) {
		t.Fatalf("zones %v, want %v", zones, want)
	}

	tests := []struct {
		dnsName  string
		zone     string
		name     string
		notFound bool
	}{
		{dnsName: "example.com", zone: "example.com", name: "@"},
		{dnsName: "WWW.example.com.", zone: "example.com", name: "www"},
		{dnsName: "a.b.sub.example.com", zone: "sub.example.com", name: "a.b"},
		{dnsName: "_sip._tcp.example.org", zone: "example.org", name: "_sip._tcp"},
		{dnsName: "badexample.com", notFound: true},
		{dnsName: "example.net", notFound: true},
	}

	for _, tc := range tests {
		zone, name, ok := provider.SplitDNSName(tc.dnsName, zones)
		if ok == tc.notFound || zone != tc.zone || name != tc.name {
			t.Errorf("%s: got (%q, %q, %v)", tc.dnsName, zone, name, ok)
		}
	}
}
