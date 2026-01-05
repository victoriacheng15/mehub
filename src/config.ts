import type { PostFilter } from "./utils/posts";

export interface NavItem {
	href: string;
	text: string;
	position?: "top" | null;
}

export interface Social {
	name: string;
	href: string;
}

export interface Project {
	title: string;
	shortDescription: string;
	link: string;
	techs: string[];
}

export interface SiteConfig {
	title: string;
	name: string;
	slogan: string;
	description?: string;
	site: string;
	rss?: boolean;
	homepage: PostFilter;
	about?: {
		paragraphs: string[];
	};
	navigation: NavItem[];
	socials: Social[];
	projects: Project[];
	skills: string[];
}

export const siteConfig: SiteConfig = {
	site: "https://victoriacheng15.vercel.app/",
	title: "Victoria's Tech Hub",
	name: "Victoria Cheng",
	slogan: "Navigating the endless world of technology.",
	description:
		"Welcome to my tech hub! I’m a Software Development student at SAIT, sharing blogs, projects, and lessons along the way. After completing my internship at Shopify, I’m continuing to build skills and prepare for the next chapter of my journey in technology.",
	rss: true,
	homepage: {
		maxPosts: 5,
		tags: [],
		excludeTags: [],
	},
	about: {
		paragraphs: [
			"I didn't start in tech. I started in <span>Packaging Engineering</span>, designing physical solutions that survive real-world chaos. Now I trade CAD drawings for CI/CD pipelines and cardboard for Cloud infrastructure.",
			"At <span>Shopify</span>, I learned how a \"measure twice, cut once\" mindset prevents regressions in production code and keeps systems resilient under the weight of millions of users.",
			"These days, you'll find me tinkering in my <span>Observability Hub</span> (my homelab), reading on my Kobo (<span>Atomic Habits</span>, <span>Software Engineering at Google</span>), or thinking about how to make systems a little more resilient.",
			"This site is where the rigid logic of systems meets the messy reality of learning. I'm here to document the why behind decisions, lessons from failures, and the curiosity that keeps me digging into the details."
		]
	},
	navigation: [
		{ href: "/", text: "Home", position: "top" },
		{ href: "/about", text: "About", position: "top" },
		{ href: "/blog", text: "Blog", position: "top" },
		{ href: "/tags", text: "Tags", position: "top" },
		{ href: "/archive", text: "Archive", position: null },
	],
	socials: [
		{ name: "GitHub", href: "https://github.com/victoriacheng15" },
		{ name: "LinkedIn", href: "https://www.linkedin.com/in/victoriacheng15" },
		{ name: "YouTube", href: "https://www.youtube.com/@victoriacheng15" },
		{ name: "Buy Me a Coffee", href: "https://www.buymeacoffee.com/victoriacheng15" },
		{ name: "rss", href: "/rss.xml" },
	],
	projects: [
		{
			title: "observability-hub",
			shortDescription:
				"Self-hosted observability hub collecting system and application metrics into Postgres, visualized in Grafana with automated backups and cron-driven Go collectors.",
			link: "https://github.com/victoriacheng15/observability-hub#readme",
			techs: ["Go", "Grafana", "PostgreSQL (TimescaleDB)", "Docker"],
		},
		{
			title: "cover-craft",
			shortDescription:
				"Serverless cover image generator on Azure Functions (Node.js/Canvas). Users customize text, fonts & colors; exports PNGs + logs preferences to MongoDB for analytics.",
			link: "https://github.com/victoriacheng15/cover-craft#readme",
			techs: ["TypeScript", "Azure Functions", "Next.js", "GitHub Actions", "MongoDB Atlas"],
		},
		{
			title: "personal-reading-analytics-dashboard",
			shortDescription:
				"Fully automated reading tracker—zero infra, 100% GitHub. Go + Python pipeline with interactive analytics from Shopify, Stripe, and GitHub blogs.",
			link: "https://github.com/victoriacheng15/personal-reading-analytics-dashboard#readme",
			techs: ["Go", "Python", "Google Sheets", "Docker", "GitHub Actions", "MongoDB Atlas"],
		},
	],
	skills: [
		// Core languages
		"Go",
		"Python",
		"TypeScript",
		"JavaScript",
		// Platform & DevOps & cloud
		"Grafana",
		"Docker",
		"Linux",
		"GitHub Actions",
		"Azure",
		// Backend & frameworks
		"Flask",
		"Node.js",
		// Databases
		"PostgreSQL",
		"MongoDB",
		// Frontend & frameworks
		"React",
		"Next.js",
		"Tailwind CSS",
	],
};

// Helper functions
export function getNavHeader() {
	return siteConfig.navigation.filter((item) => item.position === "top");
}

export function formatSocialName(name: string) {
	return name.replace(/\W/g, "").toLowerCase();
}

export function formatSkillPath(skill: string) {
	return skill.replace(/\./g, "dot").replace(/\s+/g, "").toLowerCase();
}