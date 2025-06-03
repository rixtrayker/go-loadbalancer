# Build targets for Go Load Balancer

.PHONY: build build-linux build-darwin build-windows help-build

# Build for current platform
build:
	@echo "Building for current platform..."
	@mkdir -p $(BUILD_DIR)
	$(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go

# Build for Linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./main.go

# Build for Darwin
build-darwin:
	@echo "Building for Darwin..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin ./main.go

# Build for Windows
build-windows:
	@echo "Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	GOOS=windows GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows.exe ./main.go

.PHONY: build-arm64
build-arm64:
	@echo "Building for ARM64..."
	@mkdir -p $(BUILD_DIR)
	GOARCH=arm64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-arm64 ./main.go

.PHONY: build-release
build-release: clean
	@echo "Building release binaries..."
	@mkdir -p $(DIST_DIR)
	for os in $(PLATFORMS); do \
		for arch in $(ARCHITECTURES); do \
			GOOS=$$os GOARCH=$$arch $(GOCMD) $(BUILD_FLAGS) -o $(DIST_DIR)/$(BINARY_NAME)-$$os-$$arch ./main.go; \
		done \
	done
	cd $(DIST_DIR) && sha256sum * > checksums.txt

# Build help
help-build:
	@echo "Build targets:"
	@echo "  build           - Build for current platform"
	@echo "  build-linux     - Build for Linux"
	@echo "  build-darwin    - Build for Darwin"
	@echo "  build-windows   - Build for Windows"
	@echo "  build-arm64     - Build for ARM64"
	@echo "  build-release   - Build release binaries" 