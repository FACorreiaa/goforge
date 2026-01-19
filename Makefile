# GoForge CLI Makefile
# =========================================================================

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BINARY := goforge

.PHONY: all build test clean release tag push help

all: build

# =========================================================================
# Build
# =========================================================================

build: ## Build the CLI binary
	go build -ldflags="-s -w -X github.com/FACorreiaa/goforge/cmd.version=$(VERSION)" -o $(BINARY) .

install: build ## Install to $GOPATH/bin
	go install .

# =========================================================================
# Testing
# =========================================================================

test: ## Run tests
	go test -v ./...

test-coverage: ## Run tests with coverage
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# =========================================================================
# Release
# =========================================================================

tag: ## Create a new version tag (usage: make tag v=v1.0.0)
	@if [ -z "$(v)" ]; then \
		echo "Error: Version not specified. Usage: make tag v=v1.0.0"; \
		exit 1; \
	fi
	@echo "Creating tag $(v)..."
	git tag $(v)
	@echo "✅ Tag $(v) created locally"
	@echo "Run 'make push-tag v=$(v)' to push to remote"

push-tag: ## Push a tag to remote (usage: make push-tag v=v1.0.0)
	@if [ -z "$(v)" ]; then \
		echo "Error: Version not specified. Usage: make push-tag v=v1.0.0"; \
		exit 1; \
	fi
	@echo "Pushing tag $(v) to origin..."
	git push origin $(v)
	@echo "✅ Tag $(v) pushed to remote"
	@echo ""
	@echo "Users can now install with:"
	@echo "  go install github.com/FACorreiaa/goforge@$(v)"

release: ## Create and push a new version tag (usage: make release v=v1.0.0)
	@if [ -z "$(v)" ]; then \
		echo "Error: Version not specified. Usage: make release v=v1.0.0"; \
		exit 1; \
	fi
	@echo "🚀 Creating release $(v)..."
	git tag $(v)
	git push origin $(v)
	@echo ""
	@echo "✅ Release $(v) complete!"
	@echo ""
	@echo "Users can now install with:"
	@echo "  go install github.com/FACorreiaa/goforge@$(v)"
	@echo "  go install github.com/FACorreiaa/goforge@latest"

# =========================================================================
# Utilities
# =========================================================================

clean: ## Remove build artifacts
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -rf test-app/

tidy: ## Tidy Go modules
	go mod tidy

lint: ## Run linter
	golangci-lint run

# =========================================================================
# Development
# =========================================================================

demo: build ## Build and run a demo generation
	rm -rf test-app
	./$(BINARY) new test-app github.com/demo/test-app --frontend htmx
	@echo ""
	@echo "Generated test-app/ - explore the files!"

# =========================================================================
# Help
# =========================================================================

help: ## Show this help message
	@echo "GoForge CLI - Build & Release"
	@echo ""
	@echo "Usage: make [target] [v=version]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""
	@echo "Examples:"
	@echo "  make build           # Build binary"
	@echo "  make test            # Run tests"
	@echo "  make release v=v1.0.0  # Tag and push release"
