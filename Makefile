.PHONY: build clean test run lint docker

# Build variables
BINARY_NAME=go-lb
BUILD_DIR=build
CMD_DIR=cmd/go-lb

# Docker variables
DOCKER_IMAGE=go-loadbalancer
DOCKER_TAG=latest

# Default target
all: build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	@./$(BUILD_DIR)/$(BINARY_NAME) --config configs/config.yml

# Run linting
lint:
	@echo "Running linter..."
	@golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) -f deployments/docker/Dockerfile .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

# Run k6 load tests
k6-simple:
	@echo "Running simple k6 load test..."
	@k6 run tools/k6/simple.js

k6-stress:
	@echo "Running stress k6 load test..."
	@k6 run tools/k6/stress.js

# Help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  run          - Run the application"
	@echo "  lint         - Run linting"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  k6-simple    - Run simple k6 load test"
	@echo "  k6-stress    - Run stress k6 load test"
	@echo "  help         - Show this help"
