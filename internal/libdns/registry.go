package libdns

// This file is auto generated, do not modify directly
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

var Registry = map[string]func(conf [][]byte) (LibdnsProvider, error){
	"dnsexit": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdnsexit.Provider](conf)
	},
	"rfc2136": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsrfc2136.Provider](conf)
	},
	"inwx": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsinwx.Provider](conf)
	},
	"katapult": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnskatapult.Provider](conf)
	},
	"autodns": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsautodns.Provider](conf)
	},
	"hosttech": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnshosttech.Provider](conf)
	},
	"he": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnshe.Provider](conf)
	},
	"dnsmadeeasy": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdnsmadeeasy.Provider](conf)
	},
	"namedotcom": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnamedotcom.Provider](conf)
	},
	"leaseweb": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsleaseweb.Provider](conf)
	},
	"porkbun": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsporkbun.Provider](conf)
	},
	"dreamhost": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdreamhost.Provider](conf)
	},
	"timeweb": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnstimeweb.Provider](conf)
	},
	"selectel": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsselectel.Provider](conf)
	},
	"luadns": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsluadns.Provider](conf)
	},
	"dynu": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdynu.Provider](conf)
	},
	"easydns": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnseasydns.Provider](conf)
	},
	"acmeproxy": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsacmeproxy.Provider](conf)
	},
	"transip": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnstransip.Provider](conf)
	},
	"directadmin": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdirectadmin.Provider](conf)
	},
	"infomaniak": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsinfomaniak.Provider](conf)
	},
	"desec": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdesec.Provider](conf)
	},
	"dnsimple": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdnsimple.Provider](conf)
	},
	"netlify": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnetlify.Provider](conf)
	},
	"nfsn": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnfsn.Provider](conf)
	},
	"namesilo": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnamesilo.Provider](conf)
	},
	"bunny": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsbunny.Provider](conf)
	},
	"dnsupdate": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdnsupdate.Provider](conf)
	},
	"loopia": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsloopia.Provider](conf)
	},
	"mailinabox": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsmailinabox.Provider](conf)
	},
	"hexonet": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnshexonet.Provider](conf)
	},
	"alidns": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsalidns.Provider](conf)
	},
	"duckdns": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsduckdns.Provider](conf)
	},
	"njalla": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnjalla.Provider](conf)
	},
	"totaluptime": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnstotaluptime.Provider](conf)
	},
	"ddnss": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsddnss.Provider](conf)
	},
	"vercel": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsvercel.Provider](conf)
	},
	"dinahosting": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdinahosting.Provider](conf)
	},
	"mythicbeasts": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsmythicbeasts.Provider](conf)
	},
	"dynv6": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdynv6.Provider](conf)
	},
	"namecheap": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsnamecheap.Provider](conf)
	},
	"dnspod": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsdnspod.Provider](conf)
	},
	"metaname": func(conf [][]byte) (LibdnsProvider, error) {
		return initProvider[libdnsmetaname.Provider](conf)
	},
}
