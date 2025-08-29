const MAIN_URL = "https://github.com/victoriacheng15";
const README_URL = "#readme";

export const PROJECTS = [
		{
			title: "school-flask-api",
			description: "Modular Flask REST API for school records with SQLite. Containerized with Docker and integrated CI/CD via GitHub Actions, with Pytest testing, CRUD, bulk ops, and clear error handling.",
			link: `${MAIN_URL}/school-flask-api${README_URL}`,
			techs: [
				"Python",
				"Flask",
				"SQLite",
				"Docker",
				"GitHub Actions",
				"Pytest",
			],
		},
		{
			title: "cprg303-expense-tracker",
			description: "📊 Expense Tracker App – Log in with a magic link, manage transactions, view spending trends with charts, and update your profile. Built with Supabase & React Native for an academic project.",
			link: `${MAIN_URL}/cprg303-expense-tracker${README_URL}`,
			techs: ["TypeScript", "React Native", "GitHub Actions", "Supabase (PostgreSQL)", "Supabase Auth"],
		},
		{
			title: "articles-extractor",
			description:
				"Python application automating web scraping of articles (titles, links, dates, authors) from websites. Deploy via manual runs, cron/Docker, or GitHub Actions. Organizes data into Google Sheets.",
			link: `${MAIN_URL}/articles-extractor${README_URL}`,
		techs: [
			"Python",
			"Google Sheets",
			"Docker",
			"Bash",
			"Raspberry Pi",
			"GitHub Actions",
		],
	},
	{
		title: "hacker-news-next",
		description:
			"Using React (Next.js), Redux Toolkit, Vitest: Fetch top 5 stories from REST API on the homepage. Modal links to external sources, comments, and job postings. Browse more on other pages.",
		link: `${MAIN_URL}/hacker-news-next${README_URL}`,
		techs: ["TypeScript", "React", "Next.js", "Redux", "Testing Library"],
	},
	{
		title: "bubble-tea-api",
		description:
			"Node.js, Express, and MongoDB drive an app where users submit favorite bubble tea combos. The leaderboard ranks popularity, catering to enthusiasts for fresh drink choices.",
		link: `${MAIN_URL}/bubble-tea-api${README_URL}`,
		techs: ["Node.js", "Express", "MongoDB"],
	},
];
