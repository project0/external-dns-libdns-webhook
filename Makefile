.PHONY: build generate

build:
	goreleaser release --snapshot  --clean

generate:
	go generate ./...

lint:
	golangci-lint run --fix