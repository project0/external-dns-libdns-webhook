.PHONY: build test test-integration generate check-generate lint

build:
	goreleaser release --snapshot  --clean

test:
	go test -v -cover ./...

# requires INWX_USERNAME, INWX_PASSWORD and INWX_ZONE
test-integration:
	go test -v -tags integration -run INWX ./internal/provider

generate:
	go generate ./...

check-generate: generate
	git diff --exit-code -- internal/libdnsregistry/registry.go PROVIDERS.md

lint:
	go mod tidy
	golangci-lint cache clean
	golangci-lint run --fix --build-tags integration ./...