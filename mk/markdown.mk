.PHONY: md-lint md-format

# Lint all markdown files in the repository
md-lint:
	@echo "Linting Markdown files..."
	npx markdownlint-cli "**/*.md"

# Automatically fix markdown lint errors
md-format:
	@echo "Formatting Markdown files..."
	npx markdownlint-cli --fix "**/*.md"
