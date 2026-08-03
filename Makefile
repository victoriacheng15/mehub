# ==============================================================================
# ANSI Color Codes & General Settings
# ==============================================================================
BLUE  := \033[1;34m
CYAN  := \033[1;36m
GREEN := \033[1;32m
RESET := \033[0m

.PHONY: help test-all format-all

test-all:
	$(MAKE) test
	$(MAKE) test-bdd

format-all:
	$(MAKE) format
	$(MAKE) format-md

help:
	@printf "$(BLUE)Mehub SSG Build System$(RESET)\n\n"
	@printf "Usage: make $(GREEN)<target>$(RESET)\n\n"
	@printf "$(CYAN)Site Generation:$(RESET)\n"
	@printf "  $(GREEN)build$(RESET)                  Build the SSG and generate the site\n"
	@printf "  $(GREEN)ssg-build$(RESET)              Setup Go and Tailwind, then build the SSG\n\n"
	@printf "$(CYAN)Development & Testing:$(RESET)\n"
	@printf "  $(GREEN)update$(RESET)                 Update Go dependencies\n"
	@printf "  $(GREEN)vet$(RESET)                    Run Go vet and check formatting\n"
	@printf "  $(GREEN)format$(RESET)                 Format Go code\n"
	@printf "  $(GREEN)format-all$(RESET)             Format all codebase files (Go, Markdown)\n"
	@printf "  $(GREEN)test$(RESET)                   Run Go unit tests\n"
	@printf "  $(GREEN)cov$(RESET)                    Run unit tests with coverage report\n"
	@printf "  $(GREEN)test-bdd$(RESET)               Run BDD integration tests\n"
	@printf "  $(GREEN)test-all$(RESET)               Run unit and BDD tests\n\n"
	@printf "$(CYAN)Markdown:$(RESET)\n"
	@printf "  $(GREEN)lint-md$(RESET)                Lint Markdown files using npx\n"
	@printf "  $(GREEN)format-md$(RESET)              Format Markdown files using npx\n\n"
	@printf "$(CYAN)Local Dev (Podman):$(RESET)\n"
	@printf "  $(GREEN)dev-build$(RESET)              Build the dev container image\n"
	@printf "  $(GREEN)dev-run$(RESET)                Start interactive shell with repo mounted\n"
	@printf "  $(GREEN)dev-clean$(RESET)              Remove the dev container image\n\n"
	@printf "$(CYAN)Utility:$(RESET)\n"
	@printf "  $(GREEN)help$(RESET)                   Show this help message\n"

# ==============================================================================
# Go Build & Toolchain Targets
# ==============================================================================
GO_VERSION=1.26.0
GO_TAR=go$(GO_VERSION).linux-amd64.tar.gz
GO_DIR=./go-dist

.PHONY: format update vet test cov test-bdd setup-go lint

update:
	go get -u ./... && go mod tidy

vet:
	@go vet ./cmd/... ./internal/...
	@if [ -n "$$(gofmt -l cmd/ internal/)" ]; then \
		echo "Go code is not formatted. Please run 'make format':"; \
		gofmt -l cmd/ internal/; \
		exit 1; \
	fi
	@echo "✅ Go code is formatted correctly and vetted."

format:
	go fmt ./cmd/... ./internal/...
	~/go/bin/goimports -local mehub -w cmd/ internal/

test:
	go test -v ./internal/...

cov:
	go test -coverprofile=coverage.out ./internal/... && \
	go tool cover -func=coverage.out && \
	rm coverage.out

test-bdd:
	go test -v ./e2e/...

setup-go:
	@echo "Setting up Go $(GO_VERSION)..."
	@curl -sLO https://go.dev/dl/$(GO_TAR)
	@mkdir -p $(GO_DIR)
	@tar -xzf $(GO_TAR) -C $(GO_DIR)
	@rm $(GO_TAR)
	@echo "Go setup complete."

# ==============================================================================
# Tailwind CSS Build & Setup Targets
# ==============================================================================
TAILWIND_BIN=./tailwindcss
TAILWIND_VERSION=v4.3.3

.PHONY: check-tailwind-version setup-tailwind build ssg-build

check-tailwind-version:
	@LATEST_VERSION=$$(curl -sI https://github.com/tailwindlabs/tailwindcss/releases/latest | grep -i "^location:" | awk -F/ '{print $$NF}' | tr -d '\r'); \
	if [ -n "$$LATEST_VERSION" ] && [ "$$LATEST_VERSION" != "$(TAILWIND_VERSION)" ]; then \
		echo "⚠️ A newer version of Tailwind CSS is available: $$LATEST_VERSION (Current pinned: $(TAILWIND_VERSION))"; \
	fi

setup-tailwind: check-tailwind-version
	@echo "Downloading tailwind css cli $(TAILWIND_VERSION)..."
	@curl -sL https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-linux-x64 -o $(TAILWIND_BIN)
	@chmod +x $(TAILWIND_BIN)

build: setup-tailwind
	python3 scripts/audit_tags.py && \
	rm -rf dist && \
	go run ./cmd/ssg && \
	$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
	rm $(TAILWIND_BIN);

ssg-build: setup-go setup-tailwind
	python3 scripts/audit_tags.py && \
	@export PATH=$(PWD)/$(GO_DIR)/go/bin:$$PATH; \
	go run ./cmd/ssg && \
	if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi

# ==============================================================================
# Markdown Linting & Formatting
# ==============================================================================

.PHONY: lint-md format-md

# Lint all markdown files in the repository
lint-md:
	@echo "Linting Markdown files..."
	npx markdownlint-cli "**/*.md"

# Automatically fix markdown lint errors
format-md:
	@echo "Formatting Markdown files..."
	npx markdownlint-cli --fix "**/*.md"

# ==============================================================================
# Docker / Development Container
# ==============================================================================

DEV_IMAGE=mehub-dev
DEV_CONTAINER=mehub-dev

.PHONY: dev-build dev-run dev-clean

## dev-build: Build the local development container image.
dev-build:
	podman build -f Dockerfile.dev -t $(DEV_IMAGE) .

## dev-run: Run an interactive shell inside the dev container with the repo mounted.
dev-run:
	podman run --rm -it \
		-v "$(PWD)":/workspace:Z \
		-w /workspace \
		-p 8080:8080 \
		--name $(DEV_CONTAINER) \
		$(DEV_IMAGE)

## dev-clean: Remove the local development container image.
dev-clean:
	podman rmi --force $(DEV_IMAGE)
	@echo "Image '$(DEV_IMAGE)' removed."
