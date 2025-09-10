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
		title: "school-flask-api",
		shortDescription:
			"Modular Flask REST API for school records with SQLite. Containerized with Docker and integrated CI/CD via GitHub Actions, with Pytest testing, CRUD, bulk ops, and clear error handling.",
		highlights: [
			"Implemented modular MVC architecture for the backend and clean separation of routes, services, and models; no frontend UI — endpoints return JSON for viewing.",
			"Automated CI/CD with GitHub Actions (Pytest runs, coverage reports, linting).",
			"Dockerized for consistent local development",
			"Designed bulk processing utilities (handle_bulk_process, build_bulk_response) to handle single/batch requests.",
		],
		link: `${MAIN_URL}/school-flask-api${README_URL}`,
		techs: ["Python", "Flask", "SQLite", "PostgreSQL", "Docker", "GitHub Actions", "Pytest"],
	},
	{
		title: "cprg303-expense-tracker",
		shortDescription:
			"📊 Expense Tracker App – Log in with a magic link, manage transactions, view spending trends with charts, and update your profile. Built with Supabase & React Native for an academic project.",
		highlights: [
			"Developed mobile-first UI with React Native and TypeScript.",
			"Integrated Supabase Auth (passwordless magic link login).",
			"Built PostgreSQL-backed transaction and profile management with Supabase.",
			"Visualized spending patterns via chart components.",
			"Set up GitHub Actions to run linter on pull requests.",
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
			"Python application automating web scraping of articles (titles, links, dates, authors) from websites. Deploy via manual runs, cron/Docker, or GitHub Actions. Organizes data into Google Sheets.",
		highlights: [
			"Automated web scraping with Python BeautifulSoup.",
			"Exported structured data directly into Google Sheets via API.",
			"Packaged into Docker container for portable deployments.",
			"Provided flexible deployment options: run locally, schedule via cron (including Raspberry Pi), or automate with GitHub Actions.",
			"Demonstrated DevOps-style pipelines for data ingestion and reporting.",
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
		title: "hacker-news-next",
		shortDescription:
			"Using React (Next.js), Redux Toolkit, Vitest: Fetch top 5 stories from REST API on the homepage. Modal links to external sources, comments, and job postings. Browse more on other pages.",
		highlights: [
			"Built a Next.js app to retrieve and render real-time Hacker News content using the official API.",
			"Implemented state management with Redux Toolkit.",
			"Designed modal-based navigation for seamless UX.",
			"Tested components with Vitest and Testing Library.",
		],
		link: `${MAIN_URL}/hacker-news-next${README_URL}`,
		techs: ["TypeScript", "React", "Next.js", "Redux", "Testing Library"],
	},
	{
		title: "bubble-tea-api",
		shortDescription:
			"Node.js, Express, and MongoDB drive an app where users submit favorite bubble tea combos. The leaderboard ranks popularity, catering to enthusiasts for fresh drink choices.",
		highlights: [
			"Built a RESTful API with Node.js and Express following MVC architecture.",
			"Designed MongoDB schemas for flexible storage of user-submitted combos and leaderboard rankings.",
			"Implemented leaderboard logic to sort drinks by popularity based on submission volume.",
			"Enabled Read (view combos/rankings) and Create (submit new combos) operations for core user interactions."
		],
		link: `${MAIN_URL}/bubble-tea-api${README_URL}`,
		techs: ["Node.js", "Express", "MongoDB"],
	},
];
