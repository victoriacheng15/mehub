# Agent Guide for Mehub

This document provides context and instructions for AI agents working on the **Mehub** project, a custom Static Site Generator (SSG) written in Go.

## 1. Project Overview

**Mehub** is a personal website and blog platform built with a custom Go-based SSG. It emphasizes performance, zero external runtime dependencies, and AI-first discoverability.

- **Core Tech**: Go (Golang) 1.26
- **Styling**: Tailwind CSS (standalone CLI)
- **Content**: Markdown with YAML frontmatter
- **Architecture**: Single-binary generator that renders templates into a `dist/` directory.

## 2. Project Architecture

| Component | Path | Description |
| :--- | :--- | :--- |
| **Entry Point** | `cmd/ssg/main.go` | Orchestrates the build process and clean/setup tasks. |
| **Generator** | `internal/generator.go` | Core logic for rendering HTML, RSS, Sitemaps, file copy helpers, and API registries. |
| **Schemas** | `internal/schema.go` | Unified definitions of data structures, site configurations, and post models. |
| **Content Parsing** | `internal/content.go` | Config loaders and Markdown parsing (goldmark), tag grouping, and post processing. |
| **E2E Tests** | `e2e/` | BDD (Cucumber/Godog) end-to-end integration feature tests. |
| **Config & Metadata** | `internal/templates/contents/` | YAML configuration for site metadata, projects, and skills. |
| **Templates & Assets** | `internal/templates/` | Go HTML templates, Tailwind CSS, and static assets (icons, images). |

## 3. Build and Development

The project uses a `Makefile` to orchestrate build, test, and formatting tasks.

### Site Generation

| Command | Description |
| :--- | :--- |
| `make build` | **Primary Command**. Downloads Tailwind, builds the SSG, and generates the full site in `dist/`. |

### Development

| Command | Description |
| :--- | :--- |
| `make update` | Updates Go dependencies and tidies `go.mod`. |
| `make vet` | Verifies code formatting (`gofmt`) and static analysis (`go vet`) under `cmd/` and `internal/`. |
| `make format` | Formats all Go code under `cmd/` and `internal/` using `go fmt` and `goimports`. |
| `make test` | Runs Go unit tests under `internal/`. |
| `make cov` | Runs Go unit tests with coverage report under `internal/`. |
| `make test-bdd` | Runs Cucumber/Godog E2E BDD feature integration tests under `e2e/`. |
| `make test-all` | Executes both Go unit tests and E2E BDD integration tests. |

### Markdown Files

| Command | Description |
| :--- | :--- |
| `make md-lint` | Lints all Markdown files using `markdownlint-cli`. |
| `make md-format` | Automatically formats all Markdown files using `markdownlint-cli`. |

### Setup

| Command | Description |
| :--- | :--- |
| `make setup-tailwind` | Downloads and sets up the Tailwind CSS CLI locally. |
| `make setup-go` | Downloads and sets up Go locally. |
| `make ssg-build` | Sets up Go and Tailwind CLI, then builds the SSG. |

## 4. AI Discoverability & API Registries

Mehub is designed to be easily consumed by AI agents and recruitment systems via machine-readable endpoints generated in `dist/api/`.

### API Registry Map

- **Unified Context**: `/api/manifest.json` — Unified Context Hub for AI agents (Model Context Protocol). It consolidates profile, skills, projects, and blog metadata into a single queryable resource.

### Discoverability Assets

- **robots.txt**: Explicitly allows `/api/` for all crawlers.
- **sitemap.xml**: Includes both HTML pages and the API registry JSON files to ensure thorough indexing.

## 5. Guidelines

### Go

- Use idiomatic Go and prefer the standard library.
- Code **must** pass `make vet` before completion.
- Handle all errors explicitly.

### Markdown

- Maintain valid YAML frontmatter in `blog/*.md`.
- Ensure all posts have a `title`, `date` (YYYY-MM-DD), and `tags`.
- Avoid em dashes; use commas or parentheses for clarity.

### Style

- Use Tailwind utility classes in templates.
- Avoid inline styles or custom CSS blocks outside of `internal/templates/input.css`.
