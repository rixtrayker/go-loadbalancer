# Cleanup targets for Go Load Balancer

.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)
	rm -rf $(COVERAGE_DIR)
	rm -rf $(BENCH_DIR)
	rm -rf $(DIST_DIR)
	find . -name "*.test" -type f -delete
	find . -name "*.out" -type f -delete
	find . -name "*.prof" -type f -delete

.PHONY: clean-all
clean-all: clean docker-clean
	@echo "Cleaning all generated files..."
	rm -rf vendor/
	go clean -cache -testcache -modcache
	find . -name ".DS_Store" -type f -delete
	find . -name "*.swp" -type f -delete
	find . -name "*.swo" -type f -delete
	find . -name "*.tmp" -type f -delete

.PHONY: clean-deps
clean-deps:
	@echo "Cleaning dependencies..."
	rm -rf vendor/
	go clean -modcache

.PHONY: clean-cache
clean-cache:
	@echo "Cleaning Go cache..."
	go clean -cache -testcache

.PHONY: clean-docs
clean-docs:
	@echo "Cleaning documentation..."
	rm -rf $(DOCS_DIR)/*

.PHONY: clean-tmp
clean-tmp:
	@echo "Cleaning temporary files..."
	find . -name "*.tmp" -type f -delete
	find . -name "*.temp" -type f -delete
	find . -name "*.swp" -type f -delete
	find . -name "*.swo" -type f -delete
	find . -name ".DS_Store" -type f -delete 