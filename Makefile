# Include modular build segments
include mk/go.mk
include mk/tailwind.mk
include mk/markdown.mk

.PHONY: help

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
	@echo "  test-cov        Run tests and show coverage report"
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
	@echo "Utility:"
	@echo "  help            Show this help message"
