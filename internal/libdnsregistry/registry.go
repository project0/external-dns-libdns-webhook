// This file is auto generated, DO NOT EDIT.
package libdnsregistry

import (
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
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.dnsexit.com/",
				Description:   "DNSExit provides DNS hosting services.",
				Configuration: configurationDetails[libdnsdnsexit.Provider](),
			}
		},
	},
	"rfc2136": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsrfc2136.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://tools.ietf.org/html/rfc2136",
				Description:   "RFC2136 is a standard for dynamic DNS updates.",
				Configuration: configurationDetails[libdnsrfc2136.Provider](),
			}
		},
	},
	"inwx": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinwx.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.inwx.com/",
				Description:   "INWX is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsinwx.Provider](),
			}
		},
	},
	"katapult": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnskatapult.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://katapult.io/",
				Description:   "Katapult provides cloud infrastructure services.",
				Configuration: configurationDetails[libdnskatapult.Provider](),
			}
		},
	},
	"autodns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsautodns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.internetx.com/",
				Description:   "AutoDNS is a DNS management service by InterNetX.",
				Configuration: configurationDetails[libdnsautodns.Provider](),
			}
		},
	},
	"hosttech": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshosttech.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.hosttech.ch/",
				Description:   "Hosttech provides web hosting and domain services.",
				Configuration: configurationDetails[libdnshosttech.Provider](),
			}
		},
	},
	"he": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshe.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dns.he.net/",
				Description:   "Hurricane Electric provides DNS hosting services.",
				Configuration: configurationDetails[libdnshe.Provider](),
			}
		},
	},
	"dnsmadeeasy": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsmadeeasy.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dnsmadeeasy.com/",
				Description:   "DNS Made Easy provides enterprise DNS services.",
				Configuration: configurationDetails[libdnsdnsmadeeasy.Provider](),
			}
		},
	},
	"namedotcom": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamedotcom.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.name.com/",
				Description:   "Name.com is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsnamedotcom.Provider](),
			}
		},
	},
	"leaseweb": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsleaseweb.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.leaseweb.com/",
				Description:   "Leaseweb provides cloud hosting and domain services.",
				Configuration: configurationDetails[libdnsleaseweb.Provider](),
			}
		},
	},
	"porkbun": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsporkbun.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://porkbun.com/",
				Description:   "Porkbun is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsporkbun.Provider](),
			}
		},
	},
	"dreamhost": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdreamhost.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.dreamhost.com/",
				Description:   "DreamHost provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsdreamhost.Provider](),
			}
		},
	},
	"timeweb": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstimeweb.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://timeweb.com/",
				Description:   "Timeweb provides web hosting and domain services.",
				Configuration: configurationDetails[libdnstimeweb.Provider](),
			}
		},
	},
	"selectel": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsselectel.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://selectel.ru/en/",
				Description:   "Selectel provides cloud infrastructure services.",
				Configuration: configurationDetails[libdnsselectel.Provider](),
			}
		},
	},
	"luadns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsluadns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.luadns.com/",
				Description:   "LuaDNS provides DNS hosting services.",
				Configuration: configurationDetails[libdnsluadns.Provider](),
			}
		},
	},
	"dynu": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynu.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.dynu.com/",
				Description:   "Dynu provides dynamic DNS services.",
				Configuration: configurationDetails[libdnsdynu.Provider](),
			}
		},
	},
	"easydns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnseasydns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.easydns.com/",
				Description:   "EasyDNS provides DNS hosting services.",
				Configuration: configurationDetails[libdnseasydns.Provider](),
			}
		},
	},
	"transip": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstransip.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.transip.eu/",
				Description:   "TransIP provides web hosting and domain services.",
				Configuration: configurationDetails[libdnstransip.Provider](),
			}
		},
	},
	"directadmin": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdirectadmin.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.directadmin.com/",
				Description:   "DirectAdmin is a web hosting control panel.",
				Configuration: configurationDetails[libdnsdirectadmin.Provider](),
			}
		},
	},
	"infomaniak": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsinfomaniak.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.infomaniak.com/",
				Description:   "Infomaniak provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsinfomaniak.Provider](),
			}
		},
	},
	"desec": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdesec.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://desec.io/",
				Description:   "deSEC is a free DNS hosting service.",
				Configuration: configurationDetails[libdnsdesec.Provider](),
			}
		},
	},
	"dnsimple": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsimple.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dnsimple.com/",
				Description:   "DNSimple provides DNS hosting services.",
				Configuration: configurationDetails[libdnsdnsimple.Provider](),
			}
		},
	},
	"netlify": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnetlify.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.netlify.com/",
				Description:   "Netlify provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsnetlify.Provider](),
			}
		},
	},
	"nfsn": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnfsn.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.nearlyfreespeech.net/",
				Description:   "NearlyFreeSpeech.NET provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsnfsn.Provider](),
			}
		},
	},
	"namesilo": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamesilo.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.namesilo.com/",
				Description:   "NameSilo is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsnamesilo.Provider](),
			}
		},
	},
	"bunny": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsbunny.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://bunny.net/",
				Description:   "BunnyCDN provides content delivery network services.",
				Configuration: configurationDetails[libdnsbunny.Provider](),
			}
		},
	},
	"dnsupdate": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnsupdate.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dnsupdate.org/",
				Description:   "DNSUpdate provides dynamic DNS services.",
				Configuration: configurationDetails[libdnsdnsupdate.Provider](),
			}
		},
	},
	"loopia": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsloopia.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.loopia.se/",
				Description:   "Loopia provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsloopia.Provider](),
			}
		},
	},
	"mailinabox": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmailinabox.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://mailinabox.email/",
				Description:   "Mail-in-a-Box is an easy-to-deploy mail server.",
				Configuration: configurationDetails[libdnsmailinabox.Provider](),
			}
		},
	},
	"hexonet": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnshexonet.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.hexonet.net/",
				Description:   "HEXONET provides domain registration and DNS services.",
				Configuration: configurationDetails[libdnshexonet.Provider](),
			}
		},
	},
	"alidns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsalidns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.alibabacloud.com/product/dns",
				Description:   "AliDNS is a DNS hosting service by Alibaba Cloud.",
				Configuration: configurationDetails[libdnsalidns.Provider](),
			}
		},
	},
	"duckdns": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsduckdns.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.duckdns.org/",
				Description:   "DuckDNS provides free dynamic DNS services.",
				Configuration: configurationDetails[libdnsduckdns.Provider](),
			}
		},
	},
	"njalla": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnjalla.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://njal.la/",
				Description:   "Njalla provides privacy-focused domain registration and DNS services.",
				Configuration: configurationDetails[libdnsnjalla.Provider](),
			}
		},
	},
	"totaluptime": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnstotaluptime.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://totaluptime.com/",
				Description:   "Total Uptime provides cloud-based DNS services.",
				Configuration: configurationDetails[libdnstotaluptime.Provider](),
			}
		},
	},
	"ddnss": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsddnss.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://ddnss.de/",
				Description:   "DDNSS provides free dynamic DNS services.",
				Configuration: configurationDetails[libdnsddnss.Provider](),
			}
		},
	},
	"vercel": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsvercel.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://vercel.com/",
				Description:   "Vercel provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsvercel.Provider](),
			}
		},
	},
	"dinahosting": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdinahosting.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dinahosting.com/",
				Description:   "Dinahosting provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsdinahosting.Provider](),
			}
		},
	},
	"mythicbeasts": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmythicbeasts.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.mythic-beasts.com/",
				Description:   "Mythic Beasts provides web hosting and domain services.",
				Configuration: configurationDetails[libdnsmythicbeasts.Provider](),
			}
		},
	},
	"dynv6": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdynv6.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://dynv6.com/",
				Description:   "Dynv6 provides free dynamic DNS services.",
				Configuration: configurationDetails[libdnsdynv6.Provider](),
			}
		},
	},
	"namecheap": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsnamecheap.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.namecheap.com/",
				Description:   "Namecheap is a domain registrar and DNS hosting provider.",
				Configuration: configurationDetails[libdnsnamecheap.Provider](),
			}
		},
	},
	"dnspod": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsdnspod.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://www.dnspod.com/",
				Description:   "DNSPod provides DNS hosting services.",
				Configuration: configurationDetails[libdnsdnspod.Provider](),
			}
		},
	},
	"metaname": &RegistryProvider{
		Init: func(conf [][]byte) (Provider, error) {
			return initProvider[libdnsmetaname.Provider](conf)
		},
		Docs: func() RegistryProviderDocs {
			return RegistryProviderDocs{
				URL:           "https://metaname.net/",
				Description:   "Metaname provides domain registration and DNS services.",
				Configuration: configurationDetails[libdnsmetaname.Provider](),
			}
		},
	},
}
