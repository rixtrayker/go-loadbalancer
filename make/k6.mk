# K6 Load Testing Makefile
# =====================================================

# Default values for k6 tests
K6_TARGET ?= http://loadbalancer:8080
K6_DURATION ?= 1m
K6_VUS ?= 20
K6_OUT_DIR ?= ./results
K6_SCENARIO ?= simple

# K6 Load Testing targets
.PHONY: k6-simple k6-stress k6-spike k6-endurance k6-monitor k6-custom
.PHONY: k6-validate k6-setup k6-cleanup k6-stop k6-logs k6-results
.PHONY: k6-smoke k6-load k6-breakpoint k6-volume k6-report k6-compare

# Validation and setup
k6-validate:
	@echo "Validating k6 setup..."
	@if ! docker-compose config > /dev/null 2>&1; then \
		echo "Error: docker-compose.yml is invalid"; \
		exit 1; \
	fi
	@if [ ! -d "./k6" ]; then \
		echo "Error: k6 scripts directory not found"; \
		exit 1; \
	fi
	@echo "K6 setup validation passed"

k6-setup: k6-validate
	@echo "Setting up k6 testing environment..."
	@mkdir -p $(K6_OUT_DIR)
	@docker-compose pull k6 influxdb grafana
	@echo "K6 environment ready"

# Basic load tests
k6-simple: k6-setup
	@echo "Running simple load test..."
	@echo "Target: $(K6_TARGET) | Duration: $(K6_DURATION) | VUs: $(K6_VUS)"
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/simple-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/simple.js

k6-smoke: k6-setup
	@echo "Running smoke test (minimal load)..."
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=30s \
		-e VUS=1 \
		--out json=./results/smoke-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/simple.js

k6-load: k6-setup
	@echo "Running load test..."
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/load-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/simple.js

k6-stress: k6-setup
	@echo "Running stress test..."
	@echo "Target: $(K6_TARGET) | Duration: $(K6_DURATION) | VUs: $(K6_VUS)"
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/stress-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/stress.js

k6-spike: k6-setup
	@echo "Running spike test..."
	@echo "Target: $(K6_TARGET) | Duration: $(K6_DURATION) | VUs: $(K6_VUS)"
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/spike-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/spike.js

k6-endurance: k6-setup
	@echo "Running endurance test..."
	@echo "Target: $(K6_TARGET) | Duration: $(K6_DURATION) | VUs: $(K6_VUS)"
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/endurance-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/endurance.js

k6-breakpoint: k6-setup
	@echo "Running breakpoint test (finding breaking point)..."
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=10m \
		-e VUS=1 \
		--out json=./results/breakpoint-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/stress.js

k6-volume: k6-setup
	@echo "Running volume test (high load, extended duration)..."
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=30m \
		-e VUS=50 \
		--out json=./results/volume-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/endurance.js

k6-custom: k6-setup
	@echo "Running custom k6 test..."
	@if [ -z "$(TEST_TYPE)" ]; then \
		echo "Error: TEST_TYPE is required. Use: make k6-custom TEST_TYPE=<type>"; \
		echo "Available types: simple, stress, spike, endurance"; \
		exit 1; \
	fi
	@if [ ! -f "./k6/$(TEST_TYPE).js" ]; then \
		echo "Error: Script ./k6/$(TEST_TYPE).js not found"; \
		exit 1; \
	fi
	docker-compose run --rm k6 run \
		-e TARGET_URL=$(K6_TARGET) \
		-e TEST_DURATION=$(K6_DURATION) \
		-e VUS=$(K6_VUS) \
		--out json=./results/$(TEST_TYPE)-$(shell date +%Y%m%d-%H%M%S).json \
		/scripts/$(TEST_TYPE).js

# Monitoring and observability
k6-monitor: k6-setup
	@echo "Starting monitoring stack..."
	docker-compose up -d influxdb grafana
	@echo "Waiting for services to be ready..."
	@sleep 10
	@echo "✅ Grafana is available at http://localhost:3000"
	@echo "✅ InfluxDB is available at http://localhost:8086"
	@echo "   Username: admin | Password: admin (Grafana)"

k6-logs:
	@echo "Showing k6 container logs..."
	docker-compose logs -f k6

k6-stop:
	@echo "Stopping k6 and monitoring services..."
	docker-compose stop k6 influxdb grafana

k6-cleanup:
	@echo "Cleaning up k6 containers and volumes..."
	docker-compose down -v
	docker-compose rm -f

# Results and reporting
k6-results:
	@echo "K6 Test Results:"
	@echo "==============="
	@if [ -d "$(K6_OUT_DIR)" ] && [ "$$(ls -A $(K6_OUT_DIR) 2>/dev/null)" ]; then \
		echo "Recent test results:"; \
		ls -la $(K6_OUT_DIR)/*.json 2>/dev/null | tail -10 || echo "No JSON results found"; \
	else \
		echo "No results directory or results found"; \
	fi

k6-report:
	@echo "Generating k6 test report..."
	@if [ -z "$(RESULT_FILE)" ]; then \
		echo "Error: RESULT_FILE is required. Use: make k6-report RESULT_FILE=<path>"; \
		echo "Available files:"; \
		ls -1 $(K6_OUT_DIR)/*.json 2>/dev/null || echo "No JSON results found"; \
		exit 1; \
	fi
	@if [ ! -f "$(RESULT_FILE)" ]; then \
		echo "Error: Result file $(RESULT_FILE) not found"; \
		exit 1; \
	fi
	@echo "Report for: $(RESULT_FILE)"
	@echo "=========================="
	@docker run --rm -v $(PWD)/$(K6_OUT_DIR):/results \
		loadimpact/k6 inspect $(RESULT_FILE) || \
		echo "Install jq for better JSON parsing: brew install jq"

k6-compare:
	@echo "Comparing k6 test results..."
	@if [ -z "$(FILE1)" ] || [ -z "$(FILE2)" ]; then \
		echo "Error: FILE1 and FILE2 are required"; \
		echo "Usage: make k6-compare FILE1=<path1> FILE2=<path2>"; \
		exit 1; \
	fi
	@echo "Comparing $(FILE1) vs $(FILE2)"
	@echo "=============================="
	@echo "This requires k6 compare tool or custom analysis"

# Quick test scenarios
k6-quick: k6-smoke
	@echo "Quick smoke test completed"

k6-full: k6-simple k6-stress k6-spike
	@echo "Full test suite completed"

k6-ci: k6-smoke k6-simple
	@echo "CI test suite completed"

# Development helpers
k6-dev:
	@echo "Starting development environment..."
	@echo "Load balancer + backends + monitoring"
	docker-compose up -d loadbalancer backend1 backend2 influxdb grafana
	@echo "Environment ready for testing"

k6-list-scripts:
	@echo "Available k6 test scripts:"
	@echo "========================="
	@ls -1 ./k6/*.js 2>/dev/null || echo "No scripts found in ./k6/"

# Help section
help-k6:
	@echo "K6 Load Testing Commands:"
	@echo "========================"
	@echo ""
	@echo "Setup & Validation:"
	@echo "  k6-setup        - Set up k6 testing environment"
	@echo "  k6-validate     - Validate k6 configuration"
	@echo "  k6-dev          - Start development environment"
	@echo ""
	@echo "Basic Tests:"
	@echo "  k6-smoke        - Smoke test (minimal load)"
	@echo "  k6-simple       - Simple load test"
	@echo "  k6-load         - Standard load test"
	@echo "  k6-stress       - Stress test (high load)"
	@echo "  k6-spike        - Spike test (sudden load increase)"
	@echo "  k6-endurance    - Endurance test (extended duration)"
	@echo ""
	@echo "Advanced Tests:"
	@echo "  k6-breakpoint   - Find breaking point"
	@echo "  k6-volume       - High volume test"
	@echo "  k6-custom       - Custom test (TEST_TYPE=<type>)"
	@echo ""
	@echo "Test Suites:"
	@echo "  k6-quick        - Quick smoke test"
	@echo "  k6-full         - Full test suite"
	@echo "  k6-ci           - CI test suite"
	@echo ""
	@echo "Monitoring & Observability:"
	@echo "  k6-monitor      - Start monitoring stack"
	@echo "  k6-logs         - Show k6 logs"
	@echo "  k6-stop         - Stop services"
	@echo "  k6-cleanup      - Clean up containers and volumes"
	@echo ""
	@echo "Results & Reporting:"
	@echo "  k6-results      - List test results"
	@echo "  k6-report       - Generate report (RESULT_FILE=<path>)"
	@echo "  k6-compare      - Compare results (FILE1=<path> FILE2=<path>)"
	@echo ""
	@echo "Utilities:"
	@echo "  k6-list-scripts - List available test scripts"
	@echo ""
	@echo "Configuration Variables:"
	@echo "  K6_TARGET       - Target URL (default: http://loadbalancer:8080)"
	@echo "  K6_DURATION     - Test duration (default: 1m)"
	@echo "  K6_VUS          - Virtual Users (default: 20)"
	@echo "  K6_OUT_DIR      - Output directory (default: ./results)"
	@echo ""
	@echo "Examples:"
	@echo "  make k6-stress K6_TARGET=http://localhost:8080 K6_DURATION=5m K6_VUS=100"
	@echo "  make k6-custom TEST_TYPE=spike K6_VUS=50"
	@echo "  make k6-report RESULT_FILE=./results/stress-20231201-120000.json"
	@echo "  make k6-compare FILE1=./results/before.json FILE2=./results/after.json" 