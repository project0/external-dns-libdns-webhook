# ExternalDNS - libdns provider webhook

An [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) [webhook provider](https://kubernetes-sigs.github.io/external-dns/latest/docs/tutorials/webhook-provider/) for DNS providers implemented with [libdns](https://github.com/libdns/libdns).

- Supported providers and their configuration: [PROVIDERS.md](PROVIDERS.md)
- Record types: `A`, `AAAA`, `CNAME`, `TXT`, `MX`, `SRV`, `NS`
- Built for external-dns v0.23.0 (webhook API version 1) and libdns v1.1.1

## Configuration

All flags can be set as environment variables (`--provider.name` → `LIBDNS_PROVIDER_NAME`).

| Variable | Default | Description |
| --- | --- | --- |
| `LIBDNS_PROVIDER_NAME` | | libdns provider, see [PROVIDERS.md](PROVIDERS.md) |
| `LIBDNS_PROVIDER_CONFIG` | `{}` | provider configuration as JSON |
| `LIBDNS_PROVIDER_ZONES` | | comma separated zones to manage, discovered from the provider if empty and supported |
| `LIBDNS_PROVIDER_DEFAULT_TTL` | `1h` | TTL of records without TTL |
| `LIBDNS_PROVIDER_MIN_TTL` | `0s` | minimum TTL of the DNS provider, lower TTLs are raised |
| `LIBDNS_WEBHOOK_LISTEN` | `localhost:8888` | webhook API address, keep it on localhost |
| `LIBDNS_WEBHOOK_HEALTH_LISTEN` | `:8080` | `/healthz` and `/readyz` for probes, empty to disable |
| `LIBDNS_WEBHOOK_READ_TIMEOUT` / `LIBDNS_WEBHOOK_WRITE_TIMEOUT` | `1m` | HTTP server timeouts |
| `LIBDNS_LOG_LEVEL` | `info` | `trace` logs request and response bodies |
| `LIBDNS_LOG_FORMAT` | `text` | `text` or `json` |

Set `LIBDNS_PROVIDER_MIN_TTL` when the DNS provider raises low TTLs, otherwise external-dns sees a different TTL and updates the record on every sync.

### INWX

```sh
LIBDNS_PROVIDER_NAME=inwx
LIBDNS_PROVIDER_CONFIG='{"username":"...","password":"..."}'
LIBDNS_PROVIDER_ZONES=example.com
LIBDNS_PROVIDER_MIN_TTL=300s
```

The optional `shared_secret` enables 2FA and `endpoint_url` selects another API endpoint, e.g. the OTE test system `https://api.ote.domrobot.com/jsonrpc/`.

### Testing without a DNS provider

The `fakememory` provider keeps records in memory:

```sh
LIBDNS_PROVIDER_NAME=fakememory LIBDNS_PROVIDER_CONFIG='{"zones":["example.com"]}' external-dns-libdns-webhook
```

## external-dns

Run the webhook as sidecar with the external-dns Helm chart, see [example/values.yaml](example/values.yaml):

- external-dns manages only `A`, `AAAA` and `CNAME` records by default, add the other types to `managedRecordTypes` (`--managed-record-types`).
- Set `txtPrefix` (`--txt-prefix`) containing `%{record_type}`, e.g. `edns-%{record_type}.`. By default the TXT registry names the ownership record of a zone apex record `mx-example.com`, which is outside of the zone and ignored by the webhook. Changing the prefix of an existing installation changes the names of all ownership records.
- external-dns v0.22.0 changed the annotation prefix to `external-dns.kubernetes.io/`, annotations with the old `external-dns.alpha.kubernetes.io/` prefix are ignored unless `--enable-legacy-annotation-prefix` is set.

### MX, SRV and NS records

Services and ingresses only produce `A`, `AAAA` and `CNAME` records. Create other records with the `DNSEndpoint` CRD (source `crd`), priority, weight and port are part of the target like in a zone file, see [example/dnsendpoint.yaml](example/dnsendpoint.yaml):

| Type | Target |
| --- | --- |
| `MX` | `10 mail.example.com` |
| `SRV` | `10 5 5060 sip.example.com.` (trailing dot required by external-dns) |
| `NS` | `ns1.example.net` |

## Limitations

- Routing policies (`set-identifier`) and provider specific properties (`webhook-*` annotations) are not supported and ignored.
- Changed targets are created before the old ones are deleted, so a record set is never empty. CNAME target and TTL changes use the provider's `SetRecords`, which is atomic if the provider implements it. If `SetRecords` fails (e.g. INWX with multiple records), the webhook falls back to delete and create, which leaves the record set briefly empty.
- Zones are discovered once on startup.

## Development

```sh
make generate          # regenerate the registry and PROVIDERS.md from providers.json
make test lint
INWX_USERNAME=... INWX_PASSWORD=... INWX_ZONE=example.com make test-integration
```

A provider is added by appending it to [providers.json](providers.json) and running `make generate`, it must support libdns v1.
