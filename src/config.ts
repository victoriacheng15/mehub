import type { PostFilter } from "./utils/posts";

export interface SiteConfig {
	title: string;
	name: string;
	slogan: string;
	description?: string;
	site: string;
	rss?: boolean;
	homepage: PostFilter;
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
};
