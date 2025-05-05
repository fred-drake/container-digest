# Justfile for container-digest
# This file contains recipes for common development tasks

# List all available recipes with descriptions
default:
    @just --list

# Define application name variable
app := "container-digest"

# Define target directory variable
target-dir := "target"

clean:
    @echo "Cleaning up..."
    rm -rf {{target-dir}}

# Build the container-digest application
build:
    @echo "Building {{app}} for darwin/arm64..."
    GOOS=darwin GOARCH=arm64 go build -o {{target-dir}}/darwin-arm64/{{app}} ./cmd/{{app}}
    @echo "Building {{app}} for darwin/amd64..."
    GOOS=darwin GOARCH=amd64 go build -o {{target-dir}}/darwin-amd64/{{app}} ./cmd/{{app}}
    @echo "Building {{app}} for linux/amd64..."
    GOOS=linux GOARCH=amd64 go build -o {{target-dir}}/linux-amd64/{{app}} ./cmd/{{app}}
    @echo "Building {{app}} for linux/arm64..."
    GOOS=linux GOARCH=arm64 go build -o {{target-dir}}/linux-arm64/{{app}} ./cmd/{{app}}

# Run all unit tests
test:
    @echo "Running tests..."
    go test -v ./...
