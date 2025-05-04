# Justfile for container-digest
# This file contains recipes for common development tasks

# List all available recipes with descriptions
default:
    @just --list

# Build the container-digest application
build:
    @echo "Building container-digest..."
    go build -o target/digest ./cmd/digest

# Run all unit tests
test:
    @echo "Running tests..."
    go test -v ./...
