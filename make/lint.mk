# Linting and code quality targets for Go Load Balancer

.PHONY: lint
lint:
	@echo "Running linter..."
	@command -v $(GOLINT) >/dev/null 2>&1 || { echo "Installing golangci-lint..."; go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; }
	$(GOLINT) run ./...

.PHONY: lint-fix
lint-fix:
	@echo "Running linter with auto-fix..."
	$(GOLINT) run --fix ./...

.PHONY: vet
vet:
	@echo "Running go vet..."
	$(GOVET) ./...

.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) ./...

.PHONY: fmt-check
fmt-check:
	@echo "Checking code formatting..."
	@test -z $$($(GOFMT) -l .)

.PHONY: security
security:
	@echo "Running security checks..."
	@command -v $(GOSEC) >/dev/null 2>&1 || { echo "Installing gosec..."; go install github.com/securego/gosec/v2/cmd/gosec@latest; }
	$(GOSEC) ./...

.PHONY: staticcheck
staticcheck:
	@echo "Running staticcheck..."
	@command -v staticcheck >/dev/null 2>&1 || { echo "Installing staticcheck..."; go install honnef.co/go/tools/cmd/staticcheck@latest; }
	staticcheck ./...

.PHONY: ineffassign
ineffassign:
	@echo "Checking for ineffective assignments..."
	@command -v ineffassign >/dev/null 2>&1 || { echo "Installing ineffassign..."; go install github.com/gordonklaus/ineffassign@latest; }
	ineffassign ./...

.PHONY: misspell
misspell:
	@echo "Checking for misspellings..."
	@command -v misspell >/dev/null 2>&1 || { echo "Installing misspell..."; go install github.com/client9/misspell/cmd/misspell@latest; }
	misspell -error ./...

.PHONY: code-quality
code-quality: lint vet fmt-check security staticcheck ineffassign misspell

.PHONY: install-linters
install-linters:
	@echo "Installing linters..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest
	go install github.com/gordonklaus/ineffassign@latest
	go install github.com/client9/misspell/cmd/misspell@latest 