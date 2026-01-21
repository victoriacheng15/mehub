.PHONY: help build format update vet check test cov-log setup-tailwind setup-go

BINARY_NAME=ssg.exe
TAILWIND_BIN=./tailwindcss
GO_VERSION=1.23.4
GO_TAR=go$(GO_VERSION).linux-amd64.tar.gz
GO_DIR=./go-dist

help:
	@echo "Mehub SSG Build System"
	@echo ""
	@echo "Usage:"
	@echo "  make build           Build and run the generator (requires Go)"
	@echo "  make format          Format Go code"
	@echo "  make update          Update Go dependencies"
	@echo "  make check           Check if Go code is formatted and vetted"
	@echo "  make test            Run all tests"
	@echo "  make cov-log         Run tests and show coverage report"
	@echo "  make setup-tailwind  Download Tailwind CLI (Linux x64)"
	@echo "  make setup-go        Download and setup Go $(GO_VERSION) locally"
	@echo "  make nix-<target>    Run any target inside nix-shell (e.g., nix-build)"
	@echo "  make help            Show this help message"

build:
	@export PATH=$(PWD)/$(GO_DIR)/go/bin:$$PATH; \
	go build -o $(BINARY_NAME) ./cmd/ssg
	@./$(BINARY_NAME)
	@if [ -f $(TAILWIND_BIN) ]; then \
		$(TAILWIND_BIN) -i internal/templates/input.css -o dist/styles.css --minify; \
		rm $(TAILWIND_BIN); \
	fi
	@rm $(BINARY_NAME)

nix-%:
	@nix-shell --run "make $*"

format:
	@go fmt ./...

update:
	@go get -u ./...
	@go mod tidy

vet:
	@go vet ./...

check: vet
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "Go code is not formatted. Please run 'make format':"; \
		gofmt -l .; \
		exit 1; \
	fi
	@echo "✅ Go code is formatted correctly and vetted."

test:
	@go test -v ./...

cov-log:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -func=coverage.out
	@rm coverage.out

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


