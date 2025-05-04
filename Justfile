# Justfile for container-digest
# This file contains recipes for common development tasks

# List all available recipes with descriptions
default:
    @just --list

# Build the container-digest application
build:
    @echo "Building container-digest for darwin/arm64..."
    GOOS=darwin GOARCH=arm64 go build -o target/darwin-arm64/digest ./cmd/digest
    @echo "Building container-digest for darwin/amd64..."
    GOOS=darwin GOARCH=amd64 go build -o target/darwin-amd64/digest ./cmd/digest
    @echo "Building container-digest for linux/amd64..."
    GOOS=linux GOARCH=amd64 go build -o target/linux-amd64/digest ./cmd/digest
    @echo "Building container-digest for linux/arm64..."
    GOOS=linux GOARCH=arm64 go build -o target/linux-arm64/digest ./cmd/digest

# Run all unit tests
test:
    @echo "Running tests..."
    go test -v ./...
