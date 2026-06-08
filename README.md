# Mehub

Mehub is a personal website, portfolio, and blog platform built on a custom, zero-runtime-dependency Go-based Static Site Generator (SSG).

## Design Philosophy & "Why"

- **Zero-Dependency Core**: Replaced Astro/JS framework toolchains to eliminate dependency churn and NPM package updates.
- **High-Performance Compilation**: Compiles and renders the entire site (HTML pages, XML sitemaps, RSS feeds, and JSON APIs) in under 10 seconds.

## Built With

![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4.svg?style=for-the-badge&logo=Tailwind-CSS&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF.svg?style=for-the-badge&logo=GitHub-Actions&logoColor=white)
![Bash Script](https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=for-the-badge&logo=GNU-Bash&logoColor=white)

---

## System Architecture

```mermaid
graph TD
    A[Markdown Content /blog] --> G[Go SSG Engine]
    B[YAML Configs /templates/contents] --> G
    C[HTML Templates /templates] --> G
    D[Tailwind CSS input.css] --> E[Tailwind CLI]
    G -->|Generate HTML| F[dist/]
    E -->|Compile Styles| F
```

### Key Components

- **SSG Entrypoint**: `cmd/ssg/main.go` orchestrates parsing, content compiling, and distribution directory generation.
- **Core Generator**: `internal/generator.go` renders HTML layouts, RSS feeds, and JSON endpoints.
- **Templates & Styling**: Handled via Go's standard `html/template` packages and a standalone Tailwind CSS CLI pipeline.

## Local Development & Build Commands

Build targets are automated through the root Makefile.

| Command | Action |
| :--- | :--- |
| `make build` | Primary build task. Downloads Tailwind CSS, compiles Go SSG, and builds static site into `dist/`. |
| `make ssg-build` | Continuous integration task. Prepares Go, Tailwind, compiles, and packages production assets. |
| `make test` | Executes Go unit tests. |
| `make lint` | Validates Go code styling and execution safety via `go vet`. |
| `make format` | Formats all Go codebase files using `go fmt`. |
| `make md-lint` | Analyzes Markdown consistency using `markdownlint-cli`. |
| `make md-format` | Corrects style inconsistencies in Markdown files. |

---

## Workflows & Automation

> **Philosophy**: *Automate repetition. Preserve judgment.*

Following this philosophy, the automation pipelines handle repetitive tasks while keeping integration decisions manual:

- `ci.yml`: Unifies Go formatting/vetting, Markdown linting, and unit tests under optimized path filters.
- `sync-blog-post.yml`: Automates pulling draft content from APIs and staging it as a local branch/PR.
- `publish-blog-post.yml`: Automatically processes due blog drafts by stripping draft flags and opening release PRs.

All automated PRs require manual review and merging to preserve final content judgment.
