.PHONY: build test generate lint

build:
	goreleaser release --snapshot  --clean

test:
	go test -v -cover ./...

generate:
	go generate ./...

lint:
	go mod tidy
	golangci-lint cache clean
	golangci-lint run --fix ./...