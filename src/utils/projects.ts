const MAIN_URL = "https://github.com/victoriacheng15";
const README_URL = "#readme";

export const PROJECTS = [
	{
		title: "platform-actions",
		shortDescription:
			"My personal platform tooling — reusable GitHub Actions for automation, testing, and deployment.",
		highlights: [
			"Developed a centralized repository of reusable workflows for CI/CD.",
			"Implemented Markdown linting and label-based merge automation, already in use.",
			"Improves maintainability by standardizing pipelines and removing duplicate YAML across repositories.",
		],
		link: `${MAIN_URL}hub/platform-actions${README_URL}`,
		techs: ["GitHub Actions", "CI/CD", "Automation", "YAML"],
	},
	{
		title: "cover-craft",
		shortDescription:
			"Serverless cover image generator with Azure Functions and Node Canvas. Customize text, colors, fonts, export PNGs, and track user analytics in MongoDB.",
		highlights: [
			"Node Canvas for dynamic server-side image rendering and PNG export.",
			"Responsive, reusable frontend components with Next.js.",
			"Serverless backend powered by Azure Function App.",
			"MongoDB-based user analytics tracking to visualize engagement.",
			"CI/CD pipelines via GitHub Actions for testing, build, and deployment automation.",
		],
		link: `${MAIN_URL}/cover-craft${README_URL}`,
		techs: [
			"TypeScript",
			"Azure Function App",
			"Next.js",
			"Github Actions",
			"Vitest",
			"MongoDB",
		],
	},
	{
		title: "school-management-api",
		shortDescription:
			"Modular Flask REST API for school records using Azure DB PostgreSQL. Dockerized, tested with Pytest, CI/CD integrated, featuring CRUD, bulk ops, and structured error handling.",
		highlights: [
			"Supports create, read, update, and archive (soft delete) operations.",
			"Handles both individual and bulk record updates through unified API routes.",
			"Dockerized for local development and deployed to Azure Web App with Azure Database for PostgreSQL for production.",
			"Integrated CI/CD pipelines using GitHub Actions for linting, testing, and validation.",
			"Includes automated testing with Pytest to ensure code quality and reliability.",
		],
		link: `${MAIN_URL}/school-management-api${README_URL}`,
		techs: [
			"Python",
			"Flask",
			"PostgreSQL",
			"Azure",
			"Docker",
			"GitHub Actions",
			"Pytest",
		],
	},
	{
		title: "cprg303-expense-tracker",
		shortDescription:
			"Expense Tracker App – Log in with a magic link, manage transactions, view spending trends with charts, and update your profile. Built with Supabase & React Native for an academic project.",
		highlights: [
			"Secure authentication with magic link login via Supabase Auth.",
			"Visualize spending trends with interactive charts.",
			"Implemented React Context and custom hooks for state management and component reusability.",
			"Optimized performance using useMemo for filtering and chart rendering (e.g., by year or month).",
			"PostgreSQL for reliable and consistent data storage.",
		],
		link: `${MAIN_URL}/cprg303-expense-tracker${README_URL}`,
		techs: [
			"TypeScript",
			"React Native",
			"GitHub Actions",
			"Supabase (PostgreSQL)",
			"Supabase Auth",
		],
	},
	{
		title: "article-extractor",
		shortDescription:
			"Serverless ETL pipeline: extracts, transforms, and deduplicates article metadata from multiple sources, exporting to Google Sheets with GitHub Actions and Jenkins.",
		highlights: [
			"Orchestrates a serverless ETL pipeline to extract, transform, and deduplicate articles from FreeCodeCamp, GitHub Engineering, Substack, and Shopify Engineering.",
			"Exports structured article metadata (title, link, date) to Google Sheets for centralized access and analysis.",
			"Automated scheduling via GitHub Actions; supports manual runs, Docker, and Raspberry Pi (cron).",
			"Comprehensive logging and error handling; run logs and results uploaded as GitHub artifacts for transparency and reliability.",
		],
		link: `${MAIN_URL}/article-extractor${README_URL}`,
		techs: [
			"Python",
			"Google Sheets",
			"Docker",
			"Bash Scripting",
			"Raspberry Pi",
			"Jenkins",
			"GitHub Actions",
		],
	},
];
