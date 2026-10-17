# AGENTS.md

An [external-dns](https://github.com/kubernetes-sigs/external-dns) webhook provider (webhook API version 1, targeting external-dns v0.23) for DNS providers implemented with [libdns](https://github.com/libdns/libdns) v1. User documentation: [README.md](README.md).

## Layout

| Path | Purpose |
| --- | --- |
| `main.go` | CLI flags (urfave/cli v3). Every flag is also an env var with the `LIBDNS_` prefix. |
| `providers.json` | **Single source** of the supported libdns providers. |
| `internal/generator/` | Generates `internal/libdnsregistry/registry.go` and `PROVIDERS.md` from `providers.json`. Provider config docs come from each provider's `Provider` struct via `go/ast`. |
| `internal/libdnsregistry/` | Provider interface and the registry (`registry.go` is generated). |
| `internal/externaldns/` | Local copies of the external-dns webhook types (`Endpoint`, `Changes`, `DomainFilter`). The external-dns module is deliberately not imported, to keep dependencies small. |
| `internal/provider/` | Webhook logic. `convert.go`: endpoint ↔ libdns record conversion. `externaldns.go`: Records, AdjustEndpoints, ApplyChanges. |
| `internal/webhook/` | HTTP server: the API on `localhost:8888`, health probes on `:8080`. |
| `internal/fakememory/` | In-memory libdns provider. It is selectable at runtime and backs the tests. |
| `example/` | external-dns Helm values and DNSEndpoint examples. |

## Commands

```sh
make generate        # after changing providers.json or the generator
make check-generate  # generated files must be committed and up to date (CI checks this)
make test            # go test ./...
make lint            # go mod tidy + golangci-lint (v2 config, strict)
make test-integration  # live INWX test, needs INWX_USERNAME, INWX_PASSWORD, INWX_ZONE
make build           # goreleaser snapshot (binaries + images, no push)
```

CI (`.github/workflows/test.yml`) runs check-generate, vet, and `go test -race`. Tags `v*` release with goreleaser (`release.yml`).

## Invariants

Breaking any of these makes external-dns plan updates on every sync (endless loops) or lose records.

- **Records ⇔ AdjustEndpoints must agree.** external-dns compares desired endpoints (after `AdjustEndpoints`) with `Records` exactly: targets case-insensitively, TTL when it is set, and providerSpecific exactly. Both sides go through `canonicalTarget` (`internal/provider/convert.go`). Change it in one place only, and keep the no-op plan tests green (`TestNoOpPlan*` in `internal/provider/apply_test.go`).
- **Target formats** (external-dns side):
  - MX: `"<pref> <host>"`, no trailing dot
  - SRV: `"<prio> <weight> <port> <host>."`, trailing dot required by external-dns
  - CNAME/NS: host without dot
  - TXT: unquoted

  On the libdns side, hostnames are FQDNs with a trailing dot, and endpoints are converted to **typed** records (`RR.Parse()`), because providers type switch on them.
- **Never return providerSpecific** from `Records`. `AdjustEndpoints` strips it, and TTL 0 becomes the default TTL, raised to `--provider.min-ttl`.
- **One endpoint per name and type**, with all targets. Never one endpoint per record.
- **ApplyChanges order per zone:**
  1. deletes from `Changes.Delete`
  2. `SetRecords`, only for CNAME target or TTL changes, falling back to delete+append on error
  3. appends
  4. deletes of replaced records

  This is make-before-break: a record set must never be empty in between. Deletes always use the records returned by the provider's `GetRecords`, not constructed ones.
- **Endpoints outside the managed zones** are ignored with a warning, never failing the batch.
- **`Records` fails entirely** if a zone cannot be read. Partial results would make external-dns recreate or delete records.
- **Status codes are part of the protocol.** POST /records returns `204`, and negotiation `Content-Type` is exactly `application/external.dns.webhook+json;version=1`.

## Providers

- **Adding a provider:** append it to `providers.json`, then `make generate`. The provider must build against the libdns version in `go.mod` (v1.x); providers still on libdns v0.2 cannot be added. Do not hand-edit `registry.go` or `PROVIDERS.md`.
- **Provider bugs** (e.g. INWX `SetRecords` failing on multi-value sets, INWX clamping TTLs below 300s) are worked around in the webhook in a provider-independent way, or fixed upstream. No provider name checks in `internal/provider`.

## Tests

- Unit tests use `fakememory` through the wrapper in `internal/provider/memory_test.go`, which records the call order and zone snapshots. `assertNeverEmpty` checks the make-before-break behaviour.
- Use `planChanges` (`apply_test.go`) to assert that a sync after applying plans nothing.
- The live INWX test has build tag `integration` and creates records below a random label plus one apex TXT record, removing them on cleanup.
- End to end with real external-dns: kind + external-dns Helm chart with `LIBDNS_PROVIDER_NAME=fakememory` (see `example/values.yaml`). Use a separate kubeconfig, never the default context.

## Conventions

- Work on a feature branch. Make small logical commits in Conventional Commits style (`feat(provider): …`, `fix(webhook): …`); each commit builds and passes tests.
- Match the existing style: zerolog with context loggers (`log.Ctx(ctx)`), wrapped errors (`fmt.Errorf("…: %w", err)`), `errors.Join` to collect per-zone and per-endpoint errors.
- Tests live in `_test` packages. Internals are exposed through `export_test.go`.
- Keep dependencies minimal: no external-dns or Kubernetes modules in the main module.
