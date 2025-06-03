# Docker targets for Go Load Balancer

.PHONY: docker-build docker-run docker-push help-docker

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	$(DOCKER) build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	$(DOCKER) tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

.PHONY: docker-build-multi
docker-build-multi:
	@echo "Building multi-arch Docker images..."
	$(DOCKER) buildx create --use
	$(DOCKER) buildx build --platform linux/amd64,linux/arm64 -t $(DOCKER_IMAGE):$(DOCKER_TAG) --push .

# Run Docker container
docker-run:
	@echo "Running Docker container..."
	$(DOCKER) run -p 8080:8080 $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: docker-run-dev
docker-run-dev:
	@echo "Running Docker container in development mode..."
	$(DOCKER) run -p 8080:8080 -v $(PWD):/app $(DOCKER_IMAGE):$(DOCKER_TAG)

# Push Docker image
docker-push:
	@echo "Pushing Docker image..."
	$(DOCKER) push $(DOCKER_IMAGE):$(DOCKER_TAG)
	$(DOCKER) push $(DOCKER_IMAGE):latest

.PHONY: docker-clean
docker-clean:
	@echo "Cleaning Docker images..."
	$(DOCKER) rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	$(DOCKER) rmi $(DOCKER_IMAGE):latest || true

.PHONY: docker-logs
docker-logs:
	@echo "Showing Docker container logs..."
	$(DOCKER) logs -f $$($(DOCKER) ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG))

.PHONY: docker-shell
docker-shell:
	@echo "Opening shell in Docker container..."
	$(DOCKER) exec -it $$($(DOCKER) ps -q --filter ancestor=$(DOCKER_IMAGE):$(DOCKER_TAG)) /bin/sh

.PHONY: docker-test
docker-test:
	@echo "Running tests in Docker container..."
	$(DOCKER) run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) make test

.PHONY: docker-bench
docker-bench:
	@echo "Running benchmarks in Docker container..."
	$(DOCKER) run --rm $(DOCKER_IMAGE):$(DOCKER_TAG) make bench

# Docker help
help-docker:
	@echo "Docker targets:"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Run Docker container"
	@echo "  docker-push     - Push Docker image" 