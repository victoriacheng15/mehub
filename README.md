# 👋 My Personal Website & Custom Go SSG 🐧

Welcome! This repository contains the source for my personal website, which is built and rendered by **`mehub`**, a custom Static Site Generator I wrote in Go.

This project demonstrates a complete, end-to-end ownership of a personal platform, from the core rendering engine to the automated CI/CD pipeline that publishes it.

## The 'Why': From Astro to a Custom Go SSG

This site was previously built with Astro. I migrated to a custom Go SSG to solve two key engineering challenges:

1. **Dependency Churn:** The `npm` ecosystem, even with minimal packages, required constant updates. The CI/CD install cycle alone took over 30 seconds.
2. **Performance & Simplicity:** I wanted a zero-dependency, single-binary solution. My Go SSG compiles and builds the entire site in under 2 seconds, eliminating `node_modules` and simplifying the entire toolchain.

## Features

- **Portfolio**: Showcasing my projects.
- **Blog**: Sharing insights, tutorials, and personal thoughts.
- **Tags**: Organizing content with tags for easy navigation.
- **Archives**: Exploring past content by year.
- **Custom Go SSG**: The engine that powers it all.

## Built With

![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=Go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4.svg?style=for-the-badge&logo=Tailwind-CSS&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/GitHub%20Actions-2088FF.svg?style=for-the-badge&logo=GitHub-Actions&logoColor=white)
![Bash Script](https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=for-the-badge&logo=GNU-Bash&logoColor=white)

## 🤖 Automation-First Publishing Workflow

This site is maintained through a **Git-native, mostly automated content pipeline** — designed for consistency, safety, and minimal manual work. All changes are validated, then merged with **human oversight** before going live.

### ✅ Code & Content Quality  

- **Go Formatting & Analysis**: Enforced automatically via `go fmt` and `go vet`.
- **Markdown Consistency**: Validated with [markdownlint](https://github.com/DavidAnson/markdownlint).

### 🔄 GitHub Actions Workflows  

**Content Automation** (blog-specific):  

- `sync-blog-post.yml` → Pulls draft posts from private sources.
- `publish-blog-post.yml` → Auto-publishes on scheduled UTC date.

**CI/CD & General Automation**:  

- `ci.yml` → Runs Go formatter checks (`go fmt`) on every pull request to ensure code quality.
- **Vercel Deployment** → Automatically builds and deploys the Go SSG on every push to `main` (Zero-Config deployment).
- [`markdownlint.yml`](https://github.com/victoriacheng15hub/platform-actions) → Checks Markdown files for style violations.
- [`label-based-merge.yml`](https://github.com/victoriacheng15hub/platform-actions) → Auto-merges PRs when labeled (e.g., for `go.mod` updates) after passing all checks.

### 🧠 Human-in-the-Loop  

I **manually review all blog content** before merging — via GitHub UI, CLI, or label — to verify formatting, clarity, and intent. Only low-risk changes (like dependency updates) are auto-merged.

### ⚙️ Supporting Bash Scripts  

- `sync_blog_post.sh` → Syncs content across repositories or branches.
- `publish_blog_post.sh` → Finds posts with `draft: true` and removes the flag **only if the current UTC date matches `publishDate`**.

> **Philosophy**: *Automate repetition. Preserve judgment.*
