# Main Makefile for Go Load Balancer

# Include all makefiles
include make/vars.mk
include make/build.mk
include make/test.mk
include make/bench.mk
include make/lint.mk
include make/docker.mk
include make/clean.mk
include make/help.mk

# Default target
.PHONY: all
all: clean test build

# Development environment setup
.PHONY: setup
setup: install-tools install-deps
	@echo "Development environment setup complete"

# Install dependencies
.PHONY: install-deps
install-deps:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Show help
.PHONY: help
help:
	@echo "Go Load Balancer - Available targets:"
	@echo ""
	@echo "Environment:"
	@echo "  setup           - Set up development environment"
	@echo "  install-deps    - Install dependencies"
	@echo "  install-tools   - Install development tools"
	@echo ""
	@echo "Build:"
	@echo "  build           - Build for current platform"
	@echo "  build-linux     - Build for Linux"
	@echo "  build-darwin    - Build for Darwin"
	@echo "  build-windows   - Build for Windows"
	@echo ""
	@echo "Testing:"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-race       - Run tests with race detection"
	@echo "  bench           - Run benchmarks"
	@echo "  bench-cpu       - Run CPU benchmarks"
	@echo "  bench-mem       - Run memory benchmarks"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint            - Run linter"
	@echo "  vet             - Run go vet"
	@echo "  security        - Run security checks"
	@echo "  fmt             - Format code"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-push     - Push Docker image"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean           - Clean build artifacts"
	@echo "  clean-all       - Clean all generated files"
	@echo ""
	@echo "Help:"
	@echo "  help            - Show this help message" 