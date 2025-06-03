# Variables for Go Load Balancer

# Project info
BINARY_NAME=go-lb
VERSION=0.1.0-beta
MODULE_NAME=github.com/amr/go-loadbalancer

# Directories
BUILD_DIR=build
COVERAGE_DIR=coverage
BENCH_DIR=bench
DIST_DIR=dist
DOCS_DIR=docs

# Go specific
GO=go
GOPATH=$(shell go env GOPATH)
GOBIN=$(GOPATH)/bin
GOCMD=$(GO) build
GOTEST=$(GO) test
GOVET=$(GO) vet
GOFMT=$(GO) fmt
GOMOD=$(GO) mod

# Tools
GOLINT=golangci-lint
GOSEC=gosec
GODOC=godoc
GOCOVER=go tool cover
GORACE=go test -race

# Docker
DOCKER=docker
DOCKER_REGISTRY?=your-registry
DOCKER_IMAGE=$(DOCKER_REGISTRY)/$(BINARY_NAME)
DOCKER_TAG=$(VERSION)

# Build flags
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.BuildTime=$(shell date -u '+%Y-%m-%d_%H:%M:%S')"
BUILD_FLAGS=-v $(LDFLAGS)

# Test flags
TEST_FLAGS=-v -race -cover
BENCH_FLAGS=-bench=. -benchmem

# Platform specific
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

# Environment
export GO111MODULE=on
export CGO_ENABLED=0
export GOOS=$(shell go env GOOS)
export GOARCH=$(shell go env GOARCH) 