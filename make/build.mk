# Build targets for Go Load Balancer

.PHONY: build
build:
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	$(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./main.go

.PHONY: build-all
build-all: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "Building for $@..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$@ GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$@ ./main.go

.PHONY: build-linux
build-linux:
	@echo "Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./main.go

.PHONY: build-darwin
build-darwin:
	@echo "Building for Darwin..."
	@mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOCMD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin ./main.go

.PHONY: build-windows
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