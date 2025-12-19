.PHONY: test lint build

test:
    go test ./... -v -coverprofile=coverage.out
    go tool cover -html=coverage.out -o coverage.html

lint:
    golangci-lint run ./... --timeout=5m

build:
    CGO_ENABLED=0 go build -o bin/api ./cmd/api