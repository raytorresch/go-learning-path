#!/bin/bash
# Script de construcción profesional

echo "Construyendo aplicación..."
go mod tidy
go vet ./...
go build -o bin/app ./cmd/app

echo "Build completado: bin/app"