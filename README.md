# 👋 My Personal Website 🐧

Welcome to my personal website! This site serves as a hub for my portfolio, blog, and professional information.

## Features

- **Portfolio**: Showcasing my projects
- **Blog**: Sharing insights, tutorials, and personal thoughts
- **Tags**: Organizing content with tags for easy navigation
- **Archives**: Exploring past content by year

## Built With

- ![Astro](https://img.shields.io/badge/Astro-BC52EE.svg?style=for-the-badge&logo=Astro&logoColor=white)
- ![TypeScript](https://img.shields.io/badge/TypeScript-3178C6.svg?style=for-the-badge&logo=TypeScript&logoColor=white)
- ![Tailwind CSS](https://img.shields.io/badge/Tailwind%20CSS-06B6D4.svg?style=for-the-badge&logo=Tailwind-CSS&logoColor=white)
- ![Vercel](https://img.shields.io/badge/Vercel-000000.svg?style=for-the-badge&logo=Vercel&logoColor=white)
- ![Bash Script](https://img.shields.io/badge/GNU%20Bash-4EAA25.svg?style=for-the-badge&logo=GNU-Bash&logoColor=white)
- ![Biome](https://img.shields.io/badge/Biome-29ABE2.svg?style=for-the-badge)
- ![Markdownlint](https://img.shields.io/badge/Markdownlint-000000.svg?style=for-the-badge)

## 🚀 Tooling & DevOps Practices

- **CI/CD**: Zero-config deployment powered by Vercel — automatically deploys on push to `main`
- **Code Quality**:
  - Enforced with [Biome](https://biomejs.dev/) for formatting and linting
  - Uses [markdownlint](https://github.com/DavidAnson/markdownlint) for consistent Markdown style
**GitHub Actions Workflows**:
  - `preview.yml`: Deploys a preview build to Vercel for pull requests
  - `deploy.yml`: Deploys the production site to Vercel on push to main
  - `format.yml`: Enforces code formatting and linting using Biome
  - `label-based-merge.yml`: Automatically merges pull requests when labeled
  - `markdownlint.yml`: Checks Markdown files for style consistency
  - `sync-blog-post.yml`: Adds draft blog posts from external sources to this repository
  - `publish-blog-post.yml`: Publishes blog posts on schedule
- **Bash Workflow Scripts**:
  - `sync_blog_post.sh`: Syncs blog posts between sources or branches
  - `publish_blog_post.sh`: Finds posts with `draft: true`, checks if the publish day matches the current UTC date, and removes the draft flag

## How to Explore

Visit the live site: [https://victoriacheng15.vercel.app/](https://victoriacheng15.vercel.app/)

## Development

To run this project locally:

1. Clone the repository.
2. Install dependencies using `npm install`.
3. Run the development server with `npm run dev`.
4. Open `http://localhost:3000` in your browser.

## Credits

This website is based on the [astro-zen-blog](https://github.com/larry-xue/astro-zen-blog) template by Larry Xue.
