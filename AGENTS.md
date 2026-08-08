# AGENTS.md

This file provides guidance to AI coding agents when working with code in this repository.

## Project Overview

Mehub is a personal website and blog platform built with a custom Go-based Static Site Generator (SSG). It emphasizes performance, zero external runtime dependencies, and AI-first discoverability.

- **Core Tech**: Go (Golang) 1.26
- **Styling**: Tailwind CSS (standalone CLI)
- **Content**: Markdown with YAML frontmatter
- **Architecture**: Single-binary generator that renders templates into a `dist/` directory.

## Development Commands

### Environment Setup

```bash
make update          # Update Go dependencies and tidy go.mod
```

### Local Dev Container

```bash
make dev-build   # Build development container image (Podman/Docker)
make dev-run     # Start interactive container shell with repository mounted
make dev-clean   # Remove development container image
```

### Quality Checks & Formatting

```bash
make vet         # Run go vet and verify formatting with gofmt
make format      # Format Go code with go fmt and goimports
make lint-md     # Lint Markdown files using markdownlint-cli
make format-md   # Format Markdown files using markdownlint-cli
make format-all  # Format all codebase files (Go and Markdown)
```

### Testing

```bash
make test        # Run Go unit tests under internal/
make cov         # Run Go unit tests with coverage report
make test-bdd    # Run Cucumber/Godog E2E BDD integration tests under e2e/
make test-all    # Execute both Go unit tests and E2E BDD integration tests
```

### Site Generation & Build

```bash
make build       # Primary build: download Tailwind, run SSG, and generate site in dist/
make ssg-build   # Setup Go and Tailwind CLI, then build the SSG
```

### Helper Scripts

```bash
python3 scripts/audit_tags.py            # Audit and validate tags across blog posts
python3 scripts/update_fork_cache.py     # Query GitHub for fork parent repositories and update scripts/fork_cache.json
python3 scripts/fetch_contributions.py   # Update projects.yaml with latest pull requests and issues
```

## Architecture

### Key Components

- **Entry Point (`cmd/ssg/main.go`)**: Orchestrates the build process and clean/setup tasks.
- **Generator (`internal/generator.go`)**: Core logic for rendering HTML, RSS feeds, sitemaps, file copying, and API registries.
- **Content Parsing (`internal/content.go`)**: Configuration loaders, Markdown parsing (Goldmark), tag grouping, and post processing.
- **Schemas (`internal/schema.go`)**: Unified definitions of data structures, site configurations, and post models.

### Key Directories

- `cmd/ssg/`: Application entry point.
- `internal/`: Core SSG engine, templates, and content parsing logic.
- `internal/templates/`: Go HTML templates, Tailwind CSS input (`input.css`), and static assets.
- `internal/templates/contents/`: YAML configuration files for site metadata, projects, and skills.
- `blog/`: Markdown posts with YAML frontmatter.
- `e2e/`: BDD (Cucumber/Godog) end-to-end integration feature tests.
- `scripts/`: Python helper scripts for data fetching and validation.

## Development Workflow

1. **Before Making Changes**: Ensure dependencies are installed and test runs pass (`make test-all`).
2. **Code Standards**: Run `make vet` before committing to ensure Go formatting and static checks pass.
3. **Markdown Updates**: Maintain valid YAML frontmatter in `blog/*.md` (all posts require `title`, `date` in `YYYY-MM-DD`, and `tags`). Run `make format-md` after editing Markdown.

## Common Pitfalls

1. **Go Best Practices**: Use idiomatic Go and prefer the standard library. Handle all errors explicitly.
2. **Styling Constraints**: Use Tailwind utility classes in templates. Avoid inline styles or custom CSS blocks outside of `internal/templates/input.css`.
3. **Python Standard Library**: Helper scripts must use only the Python standard library with no external pip dependencies, and must compile cleanly under Python 3 (`python3 -m py_compile`).
4. **Punctuation Rules**: Do not use em dashes in documentation or templates. Use commas, parentheses, or periods instead.
