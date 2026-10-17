.PHONY: build test generate check-generate lint

build:
	goreleaser release --snapshot  --clean

test:
	go test -v -cover ./...

generate:
	go generate ./...

check-generate: generate
	git diff --exit-code -- internal/libdnsregistry/registry.go PROVIDERS.md

lint:
	go mod tidy
	golangci-lint cache clean
	golangci-lint run --fix ./...