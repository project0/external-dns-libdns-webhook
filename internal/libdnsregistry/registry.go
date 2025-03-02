package libdnsregistry

// This file is auto generated, do not modify directly.
import (
	libdnsacmeproxy "github.com/libdns/acmeproxy"
	libdnsalidns "github.com/libdns/alidns"
	libdnsautodns "github.com/libdns/autodns"
	libdnsbunny "github.com/libdns/bunny"
	libdnsddnss "github.com/libdns/ddnss"
	libdnsdesec "github.com/libdns/desec"
	libdnsdinahosting "github.com/libdns/dinahosting"
	libdnsdirectadmin "github.com/libdns/directadmin"
	libdnsdnsexit "github.com/libdns/dnsexit"
	libdnsdnsimple "github.com/libdns/dnsimple"
	libdnsdnsmadeeasy "github.com/libdns/dnsmadeeasy"
	libdnsdnspod "github.com/libdns/dnspod"
	libdnsdnsupdate "github.com/libdns/dnsupdate"
	libdnsdreamhost "github.com/libdns/dreamhost"
	libdnsduckdns "github.com/libdns/duckdns"
	libdnsdynu "github.com/libdns/dynu"
	libdnsdynv6 "github.com/libdns/dynv6"
	libdnseasydns "github.com/libdns/easydns"
	libdnshe "github.com/libdns/he"
	libdnshexonet "github.com/libdns/hexonet"
	libdnshosttech "github.com/libdns/hosttech"
	libdnsinfomaniak "github.com/libdns/infomaniak"
	libdnsinwx "github.com/libdns/inwx"
	libdnskatapult "github.com/libdns/katapult"
	libdnsleaseweb "github.com/libdns/leaseweb"
	libdnsloopia "github.com/libdns/loopia"
	libdnsluadns "github.com/libdns/luadns"
	libdnsmailinabox "github.com/libdns/mailinabox"
	libdnsmetaname "github.com/libdns/metaname"
	libdnsmythicbeasts "github.com/libdns/mythicbeasts"
	libdnsnamecheap "github.com/libdns/namecheap"
	libdnsnamedotcom "github.com/libdns/namedotcom"
	libdnsnamesilo "github.com/libdns/namesilo"
	libdnsnetlify "github.com/libdns/netlify"
	libdnsnfsn "github.com/libdns/nfsn"
	libdnsnjalla "github.com/libdns/njalla"
	libdnsporkbun "github.com/libdns/porkbun"
	libdnsrfc2136 "github.com/libdns/rfc2136"
	libdnsselectel "github.com/libdns/selectel"
	libdnstimeweb "github.com/libdns/timeweb"
	libdnstotaluptime "github.com/libdns/totaluptime"
	libdnstransip "github.com/libdns/transip"
	libdnsvercel "github.com/libdns/vercel"
)

var registry = RegistryStore{
	"dnsexit": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsexit.Provider](conf)
		},
		URL:         "https://www.dnsexit.com/",
		Description: "DNSExit provides DNS hosting services.",
	},
	"rfc2136": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsrfc2136.Provider](conf)
		},
		URL:         "https://tools.ietf.org/html/rfc2136",
		Description: "RFC2136 is a standard for dynamic DNS updates.",
	},
	"inwx": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinwx.Provider](conf)
		},
		URL:         "https://www.inwx.com/",
		Description: "INWX is a domain registrar and DNS hosting provider.",
	},
	"katapult": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnskatapult.Provider](conf)
		},
		URL:         "https://katapult.io/",
		Description: "Katapult provides cloud infrastructure services.",
	},
	"autodns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsautodns.Provider](conf)
		},
		URL:         "https://www.internetx.com/",
		Description: "AutoDNS is a DNS management service by InterNetX.",
	},
	"hosttech": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshosttech.Provider](conf)
		},
		URL:         "https://www.hosttech.ch/",
		Description: "Hosttech provides web hosting and domain services.",
	},
	"he": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshe.Provider](conf)
		},
		URL:         "https://dns.he.net/",
		Description: "Hurricane Electric provides DNS hosting services.",
	},
	"dnsmadeeasy": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsmadeeasy.Provider](conf)
		},
		URL:         "https://dnsmadeeasy.com/",
		Description: "DNS Made Easy provides enterprise DNS services.",
	},
	"namedotcom": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamedotcom.Provider](conf)
		},
		URL:         "https://www.name.com/",
		Description: "Name.com is a domain registrar and DNS hosting provider.",
	},
	"leaseweb": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsleaseweb.Provider](conf)
		},
		URL:         "https://www.leaseweb.com/",
		Description: "Leaseweb provides cloud hosting and domain services.",
	},
	"porkbun": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsporkbun.Provider](conf)
		},
		URL:         "https://porkbun.com/",
		Description: "Porkbun is a domain registrar and DNS hosting provider.",
	},
	"dreamhost": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdreamhost.Provider](conf)
		},
		URL:         "https://www.dreamhost.com/",
		Description: "DreamHost provides web hosting and domain services.",
	},
	"timeweb": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstimeweb.Provider](conf)
		},
		URL:         "https://timeweb.com/",
		Description: "Timeweb provides web hosting and domain services.",
	},
	"selectel": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsselectel.Provider](conf)
		},
		URL:         "https://selectel.ru/en/",
		Description: "Selectel provides cloud infrastructure services.",
	},
	"luadns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsluadns.Provider](conf)
		},
		URL:         "https://www.luadns.com/",
		Description: "LuaDNS provides DNS hosting services.",
	},
	"dynu": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynu.Provider](conf)
		},
		URL:         "https://www.dynu.com/",
		Description: "Dynu provides dynamic DNS services.",
	},
	"easydns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnseasydns.Provider](conf)
		},
		URL:         "https://www.easydns.com/",
		Description: "EasyDNS provides DNS hosting services.",
	},
	"acmeproxy": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsacmeproxy.Provider](conf)
		},
		URL:         "https://github.com/libdns/acmeproxy",
		Description: "ACME Proxy is a DNS-01 challenge proxy for ACME clients.",
	},
	"transip": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstransip.Provider](conf)
		},
		URL:         "https://www.transip.eu/",
		Description: "TransIP provides web hosting and domain services.",
	},
	"directadmin": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdirectadmin.Provider](conf)
		},
		URL:         "https://www.directadmin.com/",
		Description: "DirectAdmin is a web hosting control panel.",
	},
	"infomaniak": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinfomaniak.Provider](conf)
		},
		URL:         "https://www.infomaniak.com/",
		Description: "Infomaniak provides web hosting and domain services.",
	},
	"desec": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdesec.Provider](conf)
		},
		URL:         "https://desec.io/",
		Description: "deSEC is a free DNS hosting service.",
	},
	"dnsimple": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsimple.Provider](conf)
		},
		URL:         "https://dnsimple.com/",
		Description: "DNSimple provides DNS hosting services.",
	},
	"netlify": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetlify.Provider](conf)
		},
		URL:         "https://www.netlify.com/",
		Description: "Netlify provides web hosting and domain services.",
	},
	"nfsn": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnfsn.Provider](conf)
		},
		URL:         "https://www.nearlyfreespeech.net/",
		Description: "NearlyFreeSpeech.NET provides web hosting and domain services.",
	},
	"namesilo": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamesilo.Provider](conf)
		},
		URL:         "https://www.namesilo.com/",
		Description: "NameSilo is a domain registrar and DNS hosting provider.",
	},
	"bunny": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsbunny.Provider](conf)
		},
		URL:         "https://bunny.net/",
		Description: "BunnyCDN provides content delivery network services.",
	},
	"dnsupdate": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsupdate.Provider](conf)
		},
		URL:         "https://dnsupdate.org/",
		Description: "DNSUpdate provides dynamic DNS services.",
	},
	"loopia": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsloopia.Provider](conf)
		},
		URL:         "https://www.loopia.se/",
		Description: "Loopia provides web hosting and domain services.",
	},
	"mailinabox": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmailinabox.Provider](conf)
		},
		URL:         "https://mailinabox.email/",
		Description: "Mail-in-a-Box is an easy-to-deploy mail server.",
	},
	"hexonet": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshexonet.Provider](conf)
		},
		URL:         "https://www.hexonet.net/",
		Description: "HEXONET provides domain registration and DNS services.",
	},
	"alidns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsalidns.Provider](conf)
		},
		URL:         "https://www.alibabacloud.com/product/dns",
		Description: "AliDNS is a DNS hosting service by Alibaba Cloud.",
	},
	"duckdns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsduckdns.Provider](conf)
		},
		URL:         "https://www.duckdns.org/",
		Description: "DuckDNS provides free dynamic DNS services.",
	},
	"njalla": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnjalla.Provider](conf)
		},
		URL:         "https://njal.la/",
		Description: "Njalla provides privacy-focused domain registration and DNS services.",
	},
	"totaluptime": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstotaluptime.Provider](conf)
		},
		URL:         "https://totaluptime.com/",
		Description: "Total Uptime provides cloud-based DNS services.",
	},
	"ddnss": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsddnss.Provider](conf)
		},
		URL:         "https://ddnss.de/",
		Description: "DDNSS provides free dynamic DNS services.",
	},
	"vercel": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsvercel.Provider](conf)
		},
		URL:         "https://vercel.com/",
		Description: "Vercel provides web hosting and domain services.",
	},
	"dinahosting": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdinahosting.Provider](conf)
		},
		URL:         "https://dinahosting.com/",
		Description: "Dinahosting provides web hosting and domain services.",
	},
	"mythicbeasts": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmythicbeasts.Provider](conf)
		},
		URL:         "https://www.mythic-beasts.com/",
		Description: "Mythic Beasts provides web hosting and domain services.",
	},
	"dynv6": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynv6.Provider](conf)
		},
		URL:         "https://dynv6.com/",
		Description: "Dynv6 provides free dynamic DNS services.",
	},
	"namecheap": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamecheap.Provider](conf)
		},
		URL:         "https://www.namecheap.com/",
		Description: "Namecheap is a domain registrar and DNS hosting provider.",
	},
	"dnspod": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnspod.Provider](conf)
		},
		URL:         "https://www.dnspod.com/",
		Description: "DNSPod provides DNS hosting services.",
	},
	"metaname": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmetaname.Provider](conf)
		},
		URL:         "https://metaname.net/",
		Description: "Metaname provides domain registration and DNS services.",
	},
}
