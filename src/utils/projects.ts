const MAIN_URL = "https://github.com/victoriacheng15";
const README_URL = "#readme";

export const PROJECTS = [
	{
		title: "platform-actions",
		shortDescription:
			"My personal platform tooling — reusable GitHub Actions for automation, testing, and deployment.",
		link: `${MAIN_URL}hub/platform-actions${README_URL}`,
		techs: ["GitHub Actions", "CI/CD", "Automation", "YAML"],
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
		title: "school-management-api",
		shortDescription:
			"Modular Flask REST API for school records using Azure DB PostgreSQL. Dockerized, tested with Pytest, CI/CD integrated, featuring CRUD, bulk ops, and structured error handling.",
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
		title: "article-extractor",
		shortDescription:
			"Serverless ETL pipeline: extracts, transforms, and deduplicates article metadata from multiple sources, exporting to Google Sheets with GitHub Actions and Jenkins.",
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
