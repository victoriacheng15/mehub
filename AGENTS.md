# Agent Guide for Mehub

This document provides context and instructions for AI agents working on the **Mehub** project, a custom Static Site Generator (SSG) written in Go.

## 1. Project Overview

**Mehub** is a personal website and blog platform built with a custom Go-based SSG. It emphasizes performance, zero external runtime dependencies, and AI-first discoverability.

- **Core Tech**: Go (Golang) 1.23.4
- **Styling**: Tailwind CSS (standalone CLI)
- **Content**: Markdown with YAML frontmatter
- **Architecture**: Single-binary generator that renders templates into a `dist/` directory.

## 2. Project Architecture

| Component | Path | Description |
| :--- | :--- | :--- |
| **Entry Point** | `cmd/ssg/main.go` | Orchestrates the build process and clean/setup tasks. |
| **Generator** | `internal/generator/` | Core logic for rendering HTML, RSS, Sitemaps, and API registries. |
| **Content** | `internal/content/` | Markdown parsing (goldmark), tag processing, and post data structures. |
| **Config** | `internal/config/` | YAML configuration loading for site metadata, projects, and skills. |
| **Templates** | `internal/templates/` | Go HTML templates and Tailwind input CSS. |

## 3. Build and Development

The project uses a `Makefile` that automatically wraps commands in `nix-shell` if available.

| Command | Description |
| :--- | :--- |
| `make build` | **Primary Command**. Downloads Tailwind, builds the SSG, and generates the full site in `dist/`. |
| `make test` | Runs all Go unit tests. |
| `make check` | Verifies code formatting (`gofmt`) and static analysis (`go vet`). |
| `make format` | Automatically formats all Go code. |
| `make lint` | Lints all Markdown files using `markdownlint`. |
| `make add-hr` | Helper script to insert horizontal rules (`---`) between H2 headings in blog posts. |
| `make update` | Updates Go dependencies and tidies `go.mod`. |

## 4. AI Discoverability & API Registries

Mehub is designed to be easily consumed by AI agents and recruitment systems via machine-readable endpoints generated in `dist/api/`.

### API Registry Map
- **Catalog**: `/api/catalog-registry.json` — Entry point for discovery.
- **Profile**: `/api/profile-registry.json` — Site owner metadata.
- **Skills**: `/api/skills-registry.json` — Structured skills and specialties.
- **Projects**: `/api/projects-registry.json` — Project portoflio with tech stacks.
- **Blog**: `/api/blog-registry.json` — Metadata for all blog posts for NLP processing.

### Discoverability Assets
- **robots.txt**: Explicitly allows `/api/` for all crawlers.
- **sitemap.xml**: Includes both HTML pages and the API registry JSON files to ensure thorough indexing.

## 5. Guidelines

### Go
- Use idiomatic Go and prefer the standard library.
- Code **must** pass `make check` before completion.
- Handle all errors explicitly.

### Markdown
- Maintain valid YAML frontmatter in `blog/*.md`.
- Ensure all posts have a `title`, `date` (YYYY-MM-DD), and `tags`.
- Avoid em dashes; use commas or parentheses for clarity.

### Style
- Use Tailwind utility classes in templates.
- Avoid inline styles or custom CSS blocks outside of `internal/templates/input.css`.
