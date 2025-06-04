#!/bin/bash

# Setup script for go-loadbalancer
# This script initializes the development environment

set -e

echo "Setting up Go Load Balancer development environment..."

# Check Go installation
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed. Please install Go 1.23 or higher."
    exit 1
fi

# Install development dependencies
echo "Installing development dependencies..."
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install golang.org/x/tools/cmd/goimports@latest

# Build the project
echo "Building the project..."
go build -o build/go-lb ./cmd/go-lb

# Setup test backends
echo "Setting up test backends..."
mkdir -p test/backend1 test/backend2

echo "Setup complete! You can now run the load balancer with:"
echo "./build/go-lb --config configs/config.yml"
