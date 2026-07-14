# Include modular build segments
include mk/go.mk
include mk/tailwind.mk
include mk/markdown.mk
include mk/docker.mk

.PHONY: help test-all

test-all:
	$(MAKE) test
	$(MAKE) test-bdd

help:
	@echo "Mehub SSG Build System"
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Site Generation:"
	@echo "  build           Build the SSG and generate the site"
	@echo ""
	@echo "Development:"
	@echo "  update          Update Go dependencies"
	@echo "  vet             Run Go vet and check formatting"
	@echo "  format          Format Go code"
	@echo "  test            Run all tests"
	@echo "  cov             Run tests and show coverage report"
	@echo "  test-bdd        Run BDD integration tests"
	@echo "  test-all        Run unit and BDD tests"
	@echo ""
	@echo "Markdown:"
	@echo "  md-lint         Lint Markdown files using npx"
	@echo "  md-format       Automatically format Markdown files using npx"
	@echo ""
	@echo "Setup:"
	@echo "  setup-tailwind  Download Tailwind CLI (Linux x64)"
	@echo "  setup-go        Download and setup Go locally"
	@echo "  ssg-build       Setup Go and Tailwind, then build the SSG"
	@echo ""
	@echo "Local Dev (Podman):"
	@echo "  dev-build       Build the dev container image"
	@echo "  dev-run         Start an interactive shell with repo mounted"
	@echo "  dev-clean       Remove the dev container image"
	@echo ""
	@echo "Utility:"
	@echo "  help            Show this help message"
