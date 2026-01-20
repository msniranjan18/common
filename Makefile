


# ==============================================================================
# Help
# ==============================================================================

.PHONY: help
help: ## Display this help message
	@echo "ChitChat - Messaging Application"
	@echo ""
	@echo "Usage:"
	@echo "  make <target>"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

# ==============================================================================
# Dependencies
# ==============================================================================

.PHONY: deps
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	@$(GO) mod download
	@$(GO) mod verify
	@echo "Dependencies downloaded."

.PHONY: deps-update
deps-update: ## Update all dependencies
	@echo "Updating dependencies..."
	@$(GO) get -u ./...
	@$(GO) mod tidy
	@echo "Dependencies updated."

.PHONY: deps-clean
deps-clean: ## Clean dependencies cache
	@echo "Cleaning dependencies cache..."
	@$(GO) clean -modcache
	@echo "Dependencies cache cleaned."

.PHONY: tidy
tidy: ## Tidy go.mod
	@echo "Tidying go.mod..."
	@$(GO) mod tidy
	@echo "go.mod tidied."

# ==============================================================================
# Build
# ==============================================================================

.PHONY: build
build: clean deps deps-update  ## Build the repo
	@echo "Building ..."
	@$(GO) build $(GO_BUILD_FLAGS) ./...
	@echo "Build complete"

# ==============================================================================
# Testing
# ==============================================================================

.PHONY: test
test: ## Run unit tests
	@echo "Running tests..."
	@$(GO) test $(GO_TEST_FLAGS) ./...
	@echo "Tests passed."

# ==============================================================================
# Linting & Formatting
# ==============================================================================

.PHONY: lint
lint: ## Run linter
	@echo "Running linter..."
	@if [ -x "$$(command -v golangci-lint)" ]; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
		golangci-lint run ./...; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo "Formatting Go code..."
	@$(GO) fmt ./...
	@echo "Code formatted."

.PHONY: vet
vet: ## Run go vet
	@echo "Running go vet..."
	@$(GO) vet ./...
	@echo "Vet completed."

# ==============================================================================
# Cleanup
# ==============================================================================

.PHONY: clean
clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(BIN_DIR) $(DIST_DIR) $(COVERAGE_DIR) $(LOGS_DIR)
	@$(GO) clean
	@echo "Cleanup complete."

.PHONY: clean-all
clean-all: clean deps-clean ## Clean everything including dependencies
	@echo "Cleaning everything..."
	@$(DOCKER) system prune -f
	@echo "Complete cleanup done."

# ==============================================================================
# Security
# ==============================================================================

.PHONY: security-scan
security-scan: ## Run security scan
	@echo "Running security scan..."
	@if [ -x "$$(command -v gosec)" ]; then \
		gosec ./...; \
	else \
		echo "gosec not found, installing..."; \
		go install github.com/securego/gosec/v2/cmd/gosec@latest; \
		gosec ./...; \
	fi

.PHONY: audit
audit: ## Audit dependencies
	@echo "Auditing dependencies..."
	@$(GO) list -m all | tail -n +2 | awk '{print $$1}' | xargs -n 1 $(GO) mod why

# ==============================================================================
# Default target
# ==============================================================================

.DEFAULT_GOAL := help
