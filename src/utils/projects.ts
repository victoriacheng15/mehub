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
			"Serverless cover image generator powered by Azure Function App and Node Canvas. Users can customize text, colors, and fonts, then export PNGs for posts or social media.",
		highlights: [
			"Generate custom cover images in post or square formats with editable text, colors, and fonts.",
			"Utilizes Node Canvas for dynamic server-side image rendering and PNG generation.",
			"Responsive and reusable frontend components built with Next.js.",
			"Serverless backend powered by Azure Function App.",
			"Integrated CI/CD pipelines with GitHub Actions for testing, build, and deployment automation.",
		],
		link: `${MAIN_URL}/cover-craft${README_URL}`,
		techs: [
			"TypeScript",
			"Azure Function App",
			"Next.js",
			"Github Actions",
			"Vitest",
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
			"Includes automated testing with Pytest to ensure code quality and reliability."
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
			"📊 Expense Tracker App – Log in with a magic link, manage transactions, view spending trends with charts, and update your profile. Built with Supabase & React Native for an academic project.",
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
		title: "articles-extractor",
		shortDescription:
			"Automates collection of article metadata (title, link, date) from multiple blogs and exports to Google Sheets. Supports scheduled runs via GitHub Actions or manual execution.",
		highlights: [
			"Supports scraping articles from FreeCodeCamp, GitHub Engineering blog, Substack feeds, and Shopify Engineering blog.",
			"Dockerized for local development; can run manually for debugging or on a Raspberry Pi via cron job.",
			"Integrated CI/CD workflow with GitHub Actions to schedule automated scraping (e.g., daily at 06:00 UTC).",
			"- Logs of scraped data and errors are uploaded as GitHub artifacts, ensuring visibility and reliability.",
		],
		link: `${MAIN_URL}/articles-extractor${README_URL}`,
		techs: [
			"Python",
			"Google Sheets",
			"Docker",
			"Bash Scripting",
			"Raspberry Pi",
			"GitHub Actions",
		],
	},
	{
		title: "bubble-tea-api",
		shortDescription:
			"Node.js, Express, and MongoDB drive an app where users submit favorite bubble tea combos. The leaderboard ranks popularity, catering to enthusiasts for fresh drink choices.",
		highlights: [
			"Built a RESTful API with Node.js and Express following MVC architecture.",
			"Designed MongoDB schemas for flexible storage of user-submitted combos and leaderboard rankings.",
			"Implemented leaderboard logic to sort drinks by popularity based on submission volume.",
			"Enabled Read (view combos/rankings) and Create (submit new combos) operations for core user interactions.",
		],
		link: `${MAIN_URL}/bubble-tea-api${README_URL}`,
		techs: ["Node.js", "Express", "MongoDB"],
	},
];
