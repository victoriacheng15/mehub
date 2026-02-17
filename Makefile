# === Variables ===
BINARY_NAME=ssg
TAILWIND_BIN=./tailwindcss
GO_VERSION=1.23.4
GO_TAR=go$(GO_VERSION).linux-amd64.tar.gz
GO_DIR=./go-dist
LINT_IMAGE = ghcr.io/igorshubovych/markdownlint-cli:v0.44.0

# Nix wrapper logic: Use nix-shell if available and not already inside one
# Also check if we are in a CI environment where we usually want to use system tools
USE_NIX = $(shell if command -v nix-shell >/dev/null 2>&1 && [ -z "$$IN_NIX_SHELL" ] && [ "$$GITHUB_ACTIONS" != "true" ]; then echo "yes"; else echo "no"; fi)

ifeq ($(USE_NIX),yes)
    NIX_RUN = nix-shell --run
else
    NIX_RUN = bash -c
endif

.PHONY: help build format update vet check test cov-log setup-tailwind setup-go lint vercel-build add-hr

help:
	@echo "Mehub SSG Build System"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Site Generation:"
	@echo "  build           Build the SSG and generate the site"
	@echo ""
	@echo "Development:"
	@echo "  format          Format Go code"
	@echo "  update          Update Go dependencies"
	@echo "  check           Check Go code formatting and static analysis (vet)"
	@echo "  test            Run all tests"
	@echo "  cov-log         Run tests and show coverage report"
	@echo ""
	@echo "Markdown:"
	@echo "  lint            Lint Markdown files using Docker"
	@echo "  add-hr          Add '---' separators between H2 headings in blog posts"
	@echo ""
	@echo "Setup:"
	@echo "  setup-tailwind  Download Tailwind CLI (Linux x64)"
	@echo "  setup-go        Download and setup Go $(GO_VERSION) locally"
	@echo "  vercel-build    Setup Go and Tailwind, then build (for Vercel deployment)"
	@echo ""
	@echo "Utility:"
	@echo "  help            Show this help message"

build: setup-tailwind
	@$(NIX_RUN) "go build -o $(BINARY_NAME) ./cmd/ssg && \
	./$(BINARY_NAME) && \
	if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi && \
	rm $(BINARY_NAME)"

format:
	@$(NIX_RUN) "go fmt ./..."

update:
	@$(NIX_RUN) "go get -u ./... && go mod tidy"

vet:
	@$(NIX_RUN) "go vet ./..."

check: vet
	@$(NIX_RUN) "if [ -n \"\$$(gofmt -l .)\" ]; then \
		echo \"Go code is not formatted. Please run 'make format':\"; \
		gofmt -l .; \
		exit 1; \
	fi && \
	echo \"✅ Go code is formatted correctly and vetted.\""

test:
	@$(NIX_RUN) "go test -v ./..."

cov-log:
	@$(NIX_RUN) "go test -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out && \
	rm coverage.out"

lint:
	@echo "Linting Markdown files..."
	@docker run --rm -v "$(PWD):/data" -w /data $(LINT_IMAGE) --fix "**/*.md"

add-hr:
	@$(NIX_RUN) "go run scripts/add-hr.go"

setup-tailwind:
	@curl -sLO https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.17/tailwindcss-linux-x64
	@mv tailwindcss-linux-x64 $(TAILWIND_BIN)
	@chmod +x $(TAILWIND_BIN)

setup-go:
	@echo "Setting up Go $(GO_VERSION)..."
	@curl -sLO https://go.dev/dl/$(GO_TAR)
	@mkdir -p $(GO_DIR)
	@tar -xzf $(GO_TAR) -C $(GO_DIR)
	@rm $(GO_TAR)
	@echo "Go setup complete."

vercel-build: setup-go setup-tailwind
	@export PATH=$(PWD)/$(GO_DIR)/go/bin:$$PATH; \
	go build -o $(BINARY_NAME) ./cmd/ssg && \
	./$(BINARY_NAME) && \
	if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi && \
	rm $(BINARY_NAME)
