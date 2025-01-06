module github.com/project0/external-dns-webhook-libdns

go 1.23.4

require (
	github.com/libdns/acmeproxy v0.0.0-20240622122018-d329e1aa0fc9
	github.com/libdns/alidns v1.0.3
	github.com/libdns/autodns v0.0.0-20241218101955-48afd1966c75
	github.com/libdns/azure v0.4.0
	github.com/libdns/bunny v0.1.0
	github.com/libdns/civo v0.1.27
	github.com/libdns/cloudflare v0.1.1
	github.com/libdns/ddnss v0.1.0
	github.com/libdns/desec v0.2.4
	github.com/libdns/digitalocean v0.0.0-20230728223659-4f9064657aea
	github.com/libdns/dinahosting v1.0.0
	github.com/libdns/directadmin v0.3.1
	github.com/libdns/dnsexit v1.0.0
	github.com/libdns/dnsimple v0.1.3
	github.com/libdns/dnsmadeeasy v1.1.3
	github.com/libdns/dnspod v0.0.3
	github.com/libdns/dnsupdate v0.0.0-20230728193621-2e79c50ea2ee
	github.com/libdns/dreamhost v0.1.1
	github.com/libdns/duckdns v0.2.0
	github.com/libdns/dynu v0.1.1
	github.com/libdns/dynv6 v1.0.0
	github.com/libdns/easydns v0.2.1
	github.com/libdns/exoscale v1.0.0
	github.com/libdns/gandi v1.0.3
	github.com/libdns/glesys v0.0.2
	github.com/libdns/godaddy v1.0.3
	github.com/libdns/googleclouddns v1.1.0
	github.com/libdns/he v1.0.2
	github.com/libdns/hetzner v0.0.1
	github.com/libdns/hexonet v0.1.0
	github.com/libdns/hosttech v1.1.0
	github.com/libdns/huaweicloud v0.3.1
	github.com/libdns/infomaniak v0.1.3
	github.com/libdns/inwx v0.2.1
	github.com/libdns/ionos v1.1.0
	github.com/libdns/katapult v1.0.0
	github.com/libdns/leaseweb v0.4.0
	github.com/libdns/libdns v0.2.2
	github.com/libdns/linode v0.4.1
	github.com/libdns/loopia v0.0.3
	github.com/libdns/luadns v0.1.0
	github.com/libdns/mailinabox v0.0.1
	github.com/libdns/metaname v0.3.0
	github.com/libdns/mythicbeasts v1.0.1
	github.com/libdns/namecheap v0.0.0-20211109042440-fc7440785c8e
	github.com/libdns/namedotcom v0.3.3
	github.com/libdns/namesilo v0.1.1
	github.com/libdns/netcup v0.1.1-0.20240604141625-bdf109361f52
	github.com/libdns/netlify v1.1.0
	github.com/libdns/nfsn v0.1.1
	github.com/libdns/njalla v0.0.0-20230106195713-96e29ec4149e
	github.com/libdns/openstack-designate v0.1.0
	github.com/libdns/ovh v0.0.3
	github.com/libdns/porkbun v0.2.0
	github.com/libdns/powerdns v0.1.3
	github.com/libdns/rfc2136 v0.1.1
	github.com/libdns/route53 v1.5.1
	github.com/libdns/scaleway v0.1.11
	github.com/libdns/selectel v1.0.0
	github.com/libdns/tencentcloud v1.1.0
	github.com/libdns/timeweb v1.0.1
	github.com/libdns/totaluptime v1.0.1
	github.com/libdns/transip v0.0.0-20240619142000-fc072056ed2e
	github.com/libdns/vercel v0.0.2
	github.com/libdns/vultr v1.0.0
	sigs.k8s.io/external-dns v0.15.1
)

require (
	golang.org/x/tools v0.28.0 // indirect
	k8s.io/apimachinery v0.32.0 // indirect
)

require (
	cloud.google.com/go/auth v0.13.0 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.6 // indirect
	cloud.google.com/go/compute/metadata v0.6.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.16.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.8.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns v1.2.0 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.2.2 // indirect
	github.com/adamantal/go-dreamhost v0.1.1 // indirect
	github.com/antchfx/htmlquery v1.2.5 // indirect
	github.com/antchfx/xpath v1.2.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/aws/aws-sdk-go-v2 v1.32.6 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.28.6 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.47 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/route53 v1.46.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.6 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.33.2 // indirect
	github.com/aws/smithy-go v1.22.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/boombuler/barcode v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/civo/civogo v0.3.89 // indirect
	github.com/digitalocean/godo v1.132.0 // indirect
	github.com/dnsimple/dnsimple-go v1.7.0 // indirect
	github.com/exoscale/egoscale/v3 v3.1.7 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/analysis v0.21.4 // indirect
	github.com/go-openapi/errors v0.21.0 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/loads v0.21.2 // indirect
	github.com/go-openapi/runtime v0.19.24 // indirect
	github.com/go-openapi/spec v0.20.8 // indirect
	github.com/go-openapi/strfmt v0.22.1 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/go-openapi/validate v0.20.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.19.0 // indirect
	github.com/go-resty/resty/v2 v2.16.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/s2a-go v0.1.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.4 // indirect
	github.com/googleapis/gax-go/v2 v2.14.0 // indirect
	github.com/gophercloud/gophercloud v1.14.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.7 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hexonet/go-sdk v3.5.1+incompatible // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/john-k/dnsmadeeasy v1.1.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kolo/xmlrpc v0.0.0-20201022064351-38db28db192b // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/linode/linodego v1.44.0 // indirect
	github.com/luadns/luadns-go v0.2.0 // indirect
	github.com/luv2code/gomiabdns v1.0.0 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/miekg/dns v1.1.62 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mittwald/go-powerdns v0.6.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/netlify/open-api/v2 v2.12.2 // indirect
	github.com/nrdcg/dnspod-go v0.4.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/ovh/go-ovh v1.6.0 // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pquerna/otp v1.4.0 // indirect
	github.com/prometheus/client_golang v1.20.5 // indirect
	github.com/prometheus/client_model v0.6.1 // indirect
	github.com/prometheus/common v0.55.0 // indirect
	github.com/prometheus/procfs v0.15.1 // indirect
	github.com/sagikazarmark/locafero v0.4.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.30 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.11.0 // indirect
	github.com/spf13/cast v1.6.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.18.2 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.1064 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod v1.0.1064 // indirect
	github.com/transip/gotransip/v6 v6.26.0 // indirect
	github.com/urfave/cli/v3 v3.0.0-beta1
	github.com/vultr/govultr/v3 v3.0.2 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.54.0 // indirect
	go.opentelemetry.io/otel v1.29.0 // indirect
	go.opentelemetry.io/otel/metric v1.29.0 // indirect
	go.opentelemetry.io/otel/trace v1.29.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.31.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.22.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/oauth2 v0.24.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.8.0 // indirect
	google.golang.org/api v0.213.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241209162323-e6fa225c2576 // indirect
	google.golang.org/grpc v1.67.1 // indirect
	google.golang.org/protobuf v1.35.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.32.0 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/utils v0.0.0-20241104100929-3ea5e8cea738 // indirect
	sigs.k8s.io/json v0.0.0-20241010143419-9aa6b5e7a4b3 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.2 // indirect
	sigs.k8s.io/yaml v1.4.0 // indirect
)
