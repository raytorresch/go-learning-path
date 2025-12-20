.PHONY: test lint build test-coverage lint-strict coverage-view

test:
	go test ./... -v -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./...

lint-strict:
	golangci-lint run ./... --timeout=5m

build:
	CGO_ENABLED=0 go build -o bin/api ./cmd/api

test-coverage:
	go test ./... -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html

coverage-view: test-coverage
	@echo "Abriendo reporte de cobertura..."
	open coverage.html || xdg-open coverage.html || echo "Abre manualmente: coverage.html"