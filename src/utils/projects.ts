const MAIN_URL = "https://github.com/victoriacheng15";
const README_URL = "#readme";

export const PROJECTS = [
	{
		title: "observability-platform",
		shortDescription:
			"Self-hosted observability platform collecting system and application metrics into Postgres, visualized in Grafana with automated backups and cron-driven Go collectors.",
		link: `${MAIN_URL}/observability-platform${README_URL}`,
		techs: ["Go", "Grafana", "PostgreSQL", "Docker"],
	},
	{
		title: "cover-craft",
		shortDescription:
			"Serverless cover image generator on Azure Functions (Node.js/Canvas). Users customize text, fonts & colors; exports PNGs + logs preferences to MongoDB for analytics.",
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
		title: "personal-reading-analytics-dashboard",
		shortDescription:
			"Fully automated reading tracker—zero infra, 100% GitHub. Go + Python pipeline with interactive analytics from Shopify, Stripe, and GitHub blogs.",
		link: `${MAIN_URL}/personal-reading-analytics-dashboard${README_URL}`,
		techs: ["go", "Python", "Google Sheets", "Docker", "GitHub Actions"],
	},
];
