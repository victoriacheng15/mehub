# Agent Guide for Mehub

This document provides context and instructions for AI agents working on the **Mehub** project, a custom Static Site Generator (SSG) written in Go.

## 1. Project Overview

**Mehub** is a personal website and blog platform built with a custom Go-based SSG. It was migrated from Astro to achieve zero external dependencies (aside from Go modules) and faster build times.

- **Core Tech**: Go (Golang) 1.23+
- **Styling**: Tailwind CSS (processed via standalone CLI)
- **Content**: Markdown files (parsed with `goldmark`)
- **Goal**: Single-binary simplicity, high performance, and full ownership of the rendering pipeline.

## 2. Build and Test Commands

The project relies on **Nix** for a reproducible environment. **Always use the `make nix-<target>` variants** (e.g., `make nix-test`, `make nix-update`) for all Go-related tasks to ensure the toolchain is correctly loaded.

| Command | Description |
| :--- | :--- |
| `make nix-build` | **Primary Build Command**. Builds the `ssg.exe` binary and generates the site in `dist/` inside `nix-shell`. Requires `make setup-tailwind` first. |
| `make nix-test` | Runs all Go unit tests inside `nix-shell`. |
| `make nix-check` | Verifies code formatting and runs static analysis inside `nix-shell`. |
| `make nix-format` | Automatically formats all Go code inside `nix-shell`. |
| `make nix-update` | Updates Go dependencies and runs `go mod tidy` inside `nix-shell`. |
| `make setup-tailwind` | Downloads the standalone Tailwind CSS CLI. Required before building. |

**Important:** To build the full static site with styles, you must run:

```bash
make setup-tailwind && make nix-build
```

## 3. Code Style Guidelines

### Go

- **Strict Adherence**: Code **must** pass `go fmt` and `go vet`.
- **Idiomatic Go**: Prefer standard library solutions where possible. Keep functions small and focused.
- **Error Handling**: Handle errors explicitly. Do not ignore returned errors.
- **Imports**: Group standard library imports separately from third-party imports.

### Markdown

- **Linting**: All markdown files (blog posts, documentation) must comply with `markdownlint` rules defined in `.markdownlint.json`.
- **Frontmatter**: Blog posts require valid YAML frontmatter (title, date, tags, etc.).

### HTML/CSS

- **Tailwind**: Use utility classes for styling. Avoid inline styles or custom CSS files unless absolutely necessary.
- **Templates**: HTML templates are located in `internal/templates/`. Maintain clean, semantic HTML structure.

## 4. Testing Instructions

- **Unit Tests**: Run `make nix-test` to execute the standard Go test suite.
- **Coverage**: Run `make nix-cov-log` to see a coverage report in the terminal.
- **New Features**: Any new logic in the SSG (e.g., new markdown parsing features, template functions) **must** include accompanying unit tests.

## 5. Security & Automation

- **CI/CD**: GitHub Actions handle linting, testing, and deployment.
- **File System**: The SSG reads from the local file system (`blog/`, `public/`) and writes to `dist/`. Ensure file path joins are safe to prevent directory traversal.
- **External Inputs**: The site is static, but build-time inputs (markdown files) should be treated as untrusted content to prevent build failures or injection issues.
