# 👋 My Personal Website 🐧

Welcome to my personal website! This site serves as a hub for my portfolio, blog, and professional information.

## Features

- **Portfolio**: Showcasing my projects
- **Blog**: Sharing insights, tutorials, and personal thoughts
- **Tags**: Organizing content with tags for easy navigation
- **Archives**: Exploring past content by year

## Built With

![Astro](https://img.shields.io/badge/Astro-BC52EE.svg?style=for-the-badge&logo=Astro&logoColor=white)
![TypeScript](https://img.shields.io/badge/TypeScript-3178C6.svg?style=for-the-badge&logo=TypeScript&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4.svg?style=for-the-badge&logo=Tailwind-CSS&logoColor=white)
![Vercel](https://img.shields.io/badge/Vercel-000000.svg?style=for-the-badge&logo=Vercel&logoColor=white)
![Bash Script](https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=for-the-badge&logo=GNU-Bash&logoColor=white)

## 🤖 Automation-First Publishing Workflow

This site is maintained through a **Git-native, mostly automated content pipeline** — designed for consistency, safety, and minimal manual work. All changes are validated and previewed, then merged with **human oversight** before going live.

### ✅ Code & Content Quality  

- **Formatting & linting**: Enforced automatically via [Biome](https://biomejs.dev/) (TypeScript, JavaScript, CSS)  
- **Markdown consistency**: Validated with [markdownlint](https://github.com/DavidAnson/markdownlint)

### 🔄 GitHub Actions Workflows  

**Content Automation** (blog-specific):  

- `sync-blog-post.yml` → Pulls draft posts from private sources  
- `publish-blog-post.yml` → Auto-publishes on scheduled UTC date  

**CI/CD & General Automation**:  

- `preview.yml` → Vercel preview for every pull request  
- `deploy.yml` → Production deploy on merge to `main`  
- `format.yml` → Runs Biome formatting and linting  
- [`markdownlint.yml`](https://github.com/victoriacheng15hub/platform-actions) → Checks Markdown files for style violations  
- [`label-based-merge.yml`](https://github.com/victoriacheng15hub/platform-actions) → Auto-merges PRs when labeled (e.g., for Dependabot updates — reduces manual clicks while preserving control)

### 🧠 Human-in-the-Loop  

I **manually review all blog content** before merging — via GitHub UI, CLI, or label — to verify formatting, clarity, and intent. Only low-risk changes (like dependency updates) are auto-merged after passing all checks.

### ⚙️ Supporting Bash Scripts  

- `sync_blog_post.sh` → Syncs content across repositories or branches  
- `publish_blog_post.sh` → Finds posts with `draft: true` and removes the flag **only if the current UTC date matches `publishDate`**

> **Philosophy**: *Automate repetition. Preserve judgment.*

## Credits

This website is based on the [astro-zen-blog](https://github.com/larry-xue/astro-zen-blog) template by Larry Xue.
