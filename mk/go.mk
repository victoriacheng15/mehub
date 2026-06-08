GO_VERSION=1.25.0
GO_TAR=go$(GO_VERSION).linux-amd64.tar.gz
GO_DIR=./go-dist

.PHONY: format update vet test test-cov setup-go lint

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

setup-go:
	@echo "Setting up Go $(GO_VERSION)..."
	@curl -sLO https://go.dev/dl/$(GO_TAR)
	@mkdir -p $(GO_DIR)
	@tar -xzf $(GO_TAR) -C $(GO_DIR)
	@rm $(GO_TAR)
	@echo "Go setup complete."
