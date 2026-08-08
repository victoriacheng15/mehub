# Mehub

Mehub is a personal website, portfolio, and blog platform built on a custom, zero-runtime-dependency Go-based Static Site Generator (SSG).

## Design Philosophy

- **Simplified Toolchain**: Replaced JS framework toolchains to eliminate NPM runtime dependencies and continuous package maintenance overhead.
- **Fast Compilation**: Compiles and renders the entire site (HTML pages, XML sitemaps, RSS feeds, and JSON APIs) in under 10 seconds.

## Built With

![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4.svg?style=for-the-badge&logo=Tailwind-CSS&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF.svg?style=for-the-badge&logo=GitHub-Actions&logoColor=white)
![Bash Script](https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=for-the-badge&logo=GNU-Bash&logoColor=white)

---

## System Architecture

```text
  ┌───────────────────────────┐
  │   Markdown Content        │────────┐
  │   (blog/)                 │        │
  └───────────────────────────┘        │
  ┌───────────────────────────┐        │       ┌───────────────────┐
  │   YAML Configurations     │────────┼──────>│   Go SSG Engine   │───(Generate HTML)───┐
  │   (templates/contents)    │        │       │   (cmd/ssg)       │                     │
  └───────────────────────────┘        │       └───────────────────┘                     │
  ┌───────────────────────────┐        │                                                 │
  │   HTML Templates          │────────┘                                                 ▼
  │   (templates)             │                                                      ┌───────┐
  └───────────────────────────┘                                                      │ dist/ │
                                                                                     └───────┘
  ┌───────────────────────────┐                ┌───────────────────┐                     ▲
  │   Tailwind CSS Input      │───────────────>│   Tailwind CLI    │───(Compile Styles)──┘
  │   (input.css)             │                └───────────────────┘
  └───────────────────────────┘
```

### Key Components

- **SSG Entrypoint (`cmd/ssg/main.go`)**: Orchestrates content parsing, site compilation, and distribution directory generation.
- **Core Generator (`internal/generator.go`)**: Renders HTML layouts, RSS feeds, sitemaps, and JSON API registries.
- **Content Engine (`internal/content.go`)**: Parses YAML configuration and Markdown posts with Goldmark.
- **Templates & Styling**: Standard Go `html/template` layouts paired with standalone Tailwind CSS CLI compilation.

---

## Local Development & Build Commands

Build targets are automated through the root Makefile.

### Environment Setup

| Command | Action |
| :--- | :--- |
| `make update` | Updates Go dependencies and tidies `go.mod`. |

### Local Dev Container

| Command | Action |
| :--- | :--- |
| `make dev-build` | Builds the development container image (Podman/Docker). |
| `make dev-run` | Starts an interactive container shell with the repository mounted. |
| `make dev-clean` | Removes the local development container image. |

### Quality Checks & Formatting

| Command | Action |
| :--- | :--- |
| `make vet` | Validates Go code formatting (`gofmt`) and static analysis (`go vet`). |
| `make format` | Formats all Go codebase files using `go fmt` and `goimports`. |
| `make format-all` | Formats all Go and Markdown files across the codebase. |
| `make lint-md` | Lints Markdown files using `markdownlint-cli`. |
| `make format-md` | Automatically formats and fixes Markdown files. |

### Testing

| Command | Action |
| :--- | :--- |
| `make test` | Runs Go unit tests under `internal/`. |
| `make cov` | Runs unit tests with test coverage reporting. |
| `make test-bdd` | Executes Cucumber/Godog E2E BDD integration tests under `e2e/`. |
| `make test-all` | Executes both Go unit tests and BDD integration tests. |

### Site Generation & Build

| Command | Action |
| :--- | :--- |
| `make build` | Primary build command. Downloads Tailwind CSS, executes the SSG, and generates the site in `dist/`. |
| `make ssg-build` | Prepares local Go and Tailwind tooling, then builds the SSG. |

### Helper Scripts

| Command | Action |
| :--- | :--- |
| `python3 scripts/audit_tags.py` | Audits and validates tags across blog posts. |
| `python3 scripts/update_fork_cache.py` | Queries GitHub for fork parent repositories and updates `scripts/fork_cache.json`. |
| `python3 scripts/fetch_contributions.py` | Updates `projects.yaml` with latest pull requests and issues. |

---

## Workflows & Automation

> **Philosophy**: *Automate repetition. Preserve judgment.*

```text
  ┌────────────────────────────┐
  │ sync-blog-post.yml         │────(Import Drafts)──────────────┐
  └────────────────────────────┘                                 ▼
  ┌────────────────────────────┐                          ┌─────────────┐
  │ publish-blog-post.yml      │──────(Publish PR)───────>│    blog/    │
  └────────────────────────────┘                          └─────────────┘

  ┌────────────────────────────┐                          ┌───────────────────────────────┐
  │ update-contributions.yml   │──(Fetch Contributions)──>│ templates/contents/           │
  │ (fetch_contributions.py)   │                          │ projects.yaml                 │
  └────────────────────────────┘                          └───────────────────────────────┘
  ┌────────────────────────────┐                          ┌───────────────────────────────┐
  │ update-fork-cache.yml      │──────(Cache Forks)──────>│ scripts/fork_cache.json       │
  │ (update_fork_cache.py)     │                          │                               │
  └────────────────────────────┘                          └───────────────────────────────┘
```

Automated pipelines handle repetitive tasks while keeping merge decisions manual:

- `sync-blog-post.yml`: Imports new blog drafts from remote APIs.
- `publish-blog-post.yml`: Publishes scheduled blog drafts by creating pull requests.
- `update-contributions.yml`: Updates open-source contribution metrics from GitHub into `projects.yaml`.
- `update-fork-cache.yml`: Queries GitHub API to cache parent repository metadata for forks.
- `ci.yml`: Runs tests, static analysis (`go vet`), formatting checks, and markdown linting.

All automated pull requests require manual review and merging to preserve final content judgment.
