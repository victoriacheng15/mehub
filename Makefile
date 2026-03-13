# === Variables ===
TAILWIND_BIN=./tailwindcss
LINT_IMAGE = ghcr.io/igorshubovych/markdownlint-cli:v0.44.0
docker ?= podman

.PHONY: help build format update vet test cov-log setup-tailwind setup-go lint vercel-build add-hr

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
	@echo "  vet             Run Go vet and check formatting"
	@echo "  test            Run all tests"
	@echo "  test-cov        Run tests and show coverage report"
	@echo ""
	@echo "Markdown:"
	@echo "  lint            Lint Markdown files using Docker"
	@echo ""
	@echo "Setup:"
	@echo "  setup-tailwind  Download Tailwind CLI (Linux x64)"
	@echo "  setup-go        Download and setup Go $(GO_VERSION) locally"
	@echo "  vercel-build    Setup Go and Tailwind, then build (for Vercel deployment)"
	@echo ""
	@echo "Utility:"
	@echo "  help            Show this help message"

format:
	go fmt ./...

update:
	go get -u ./... && go mod tidy

vet:
	@go vet ./...
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted. Please run 'make format':"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "✅ Go code is formatted correctly and vetted."

test:
	go test -v ./...

test-cov:
	go test -coverprofile=coverage.out ./... && \
	go tool cover -func=coverage.out && \
	rm coverage.out

lint:
	@echo "Linting Markdown files..."
	@$(docker) run --rm -v "$(PWD):/data:Z" -w /data $(LINT_IMAGE) --fix "**/*.md"

setup-tailwind:
	@echo "Downloading tailwind css cli..."
	@curl -sL https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 -o $(TAILWIND_BIN)
	@chmod +x $(TAILWIND_BIN)

build: setup-tailwind
	rm -rf dist && \
	go run ./cmd/ssg && \
	$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
	rm $(TAILWIND_BIN);

setup-go:
	@echo "Setting up Go $(GO_VERSION)..."
	@curl -sLO https://go.dev/dl/$(GO_TAR)
	@mkdir -p $(GO_DIR)
	@tar -xzf $(GO_TAR) -C $(GO_DIR)
	@rm $(GO_TAR)
	@echo "Go setup complete."

GO_VERSION=1.25.0
GO_TAR=go$(GO_VERSION).linux-amd64.tar.gz
GO_DIR=./go-dist

vercel-build: setup-go setup-tailwind
	@export PATH=$(PWD)/$(GO_DIR)/go/bin:$$PATH; \
	go run ./cmd/ssg && \
	if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi
