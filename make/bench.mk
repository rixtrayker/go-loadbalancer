# Benchmark targets for Go Load Balancer

.PHONY: bench
bench:
	@echo "Running all benchmarks..."
	@mkdir -p $(BENCH_DIR)
	$(GOTEST) $(BENCH_FLAGS) ./... > $(BENCH_DIR)/bench.out
	@echo "Benchmark results saved to $(BENCH_DIR)/bench.out"

.PHONY: bench-cpu
bench-cpu:
	@echo "Running CPU benchmarks..."
	@mkdir -p $(BENCH_DIR)
	$(GOTEST) $(BENCH_FLAGS) -cpuprofile=$(BENCH_DIR)/cpu.prof ./... > $(BENCH_DIR)/cpu.out
	@echo "CPU profile saved to $(BENCH_DIR)/cpu.prof"
	@echo "CPU benchmark results saved to $(BENCH_DIR)/cpu.out"

.PHONY: bench-mem
bench-mem:
	@echo "Running memory benchmarks..."
	@mkdir -p $(BENCH_DIR)
	$(GOTEST) $(BENCH_FLAGS) -memprofile=$(BENCH_DIR)/mem.prof ./... > $(BENCH_DIR)/mem.out
	@echo "Memory profile saved to $(BENCH_DIR)/mem.prof"
	@echo "Memory benchmark results saved to $(BENCH_DIR)/mem.out"

.PHONY: bench-alloc
bench-alloc:
	@echo "Running allocation benchmarks..."
	@mkdir -p $(BENCH_DIR)
	$(GOTEST) $(BENCH_FLAGS) -benchmem -memprofile=$(BENCH_DIR)/alloc.prof ./... > $(BENCH_DIR)/alloc.out
	@echo "Allocation profile saved to $(BENCH_DIR)/alloc.prof"
	@echo "Allocation benchmark results saved to $(BENCH_DIR)/alloc.out"

.PHONY: bench-compare
bench-compare:
	@echo "Comparing benchmark results..."
	@command -v benchcmp >/dev/null 2>&1 || { echo "Installing benchcmp..."; go install golang.org/x/tools/cmd/benchcmp@latest; }
	benchcmp $(BENCH_DIR)/old.out $(BENCH_DIR)/new.out

.PHONY: bench-all
bench-all: bench bench-cpu bench-mem bench-alloc

.PHONY: bench-clean
bench-clean:
	@echo "Cleaning benchmark artifacts..."
	rm -rf $(BENCH_DIR)/*.prof
	rm -rf $(BENCH_DIR)/*.out 