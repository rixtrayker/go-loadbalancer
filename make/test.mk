# Test targets for Go Load Balancer

.PHONY: test test-coverage test-race help-test

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) $(TEST_FLAGS) ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p $(COVERAGE_DIR)
	$(GOTEST) $(TEST_FLAGS) -coverprofile=$(COVERAGE_DIR)/coverage.out ./...
	$(GOCOVER) -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "Coverage report generated at $(COVERAGE_DIR)/coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	$(GOTEST) -race $(TEST_FLAGS) ./...

# Test help
help-test:
	@echo "Testing targets:"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage"
	@echo "  test-race       - Run tests with race detection"

.PHONY: test-verbose
test-verbose:
	@echo "Running tests with verbose output..."
	$(GOTEST) -v $(TEST_FLAGS) ./...

.PHONY: test-short
test-short:
	@echo "Running short tests..."
	$(GOTEST) -short $(TEST_FLAGS) ./...

.PHONY: test-integration
test-integration:
	@echo "Running integration tests..."
	$(GOTEST) -tags=integration $(TEST_FLAGS) ./...

.PHONY: test-fuzz
test-fuzz:
	@echo "Running fuzz tests..."
	$(GOTEST) -fuzz=. -fuzztime=30s ./...

.PHONY: test-all
test-all: test test-race test-coverage test-integration

.PHONY: test-watch
test-watch:
	@echo "Watching for changes and running tests..."
	@command -v reflex >/dev/null 2>&1 || { echo "Installing reflex..."; go install github.com/cespare/reflex@latest; }
	reflex -r '\.go$$' -s -- sh -c 'make test' 