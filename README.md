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
- **Custom GitHub Actions**:
  - Scheduled via cron to automate blog publishing workflow:
    - Creates a new branch
    - Runs `publish_post.sh` to prepare changes
    - Commits updates (if any)
    - Opens a pull request titled with the blog post title
- **Bash Workflow Scripts**:
  - `setup_new_post`: Scaffold a new blog post with frontmatter, filename, and folder based on title and date
  - `publish_post`: Finds posts with `draft: true`, checks if the publish day matches the current UTC date, and removes the draft flag
  - `setup_blog_tag`: Creates a Git tag based on the current date (`blog/yyyy-mm-dd`) and pushes it to the remote repository with a custom commit message
  - `check_pr_status`: Checks the status of pull requests via GitHub CLI

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
