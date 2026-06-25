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
