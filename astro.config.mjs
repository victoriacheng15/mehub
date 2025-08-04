import { defineConfig } from "astro/config";
import { siteConfig } from "./src/config";
import { visit } from 'unist-util-visit';
import expressiveCode from "astro-expressive-code";
import rehypeMermaid from 'rehype-mermaid';
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
	site: siteConfig.site,
	integrations: [expressiveCode()],
	markdown: {
		rehypePlugins: [rehypeMermaid],
	},
	vite: {
		plugins: [tailwindcss()],
	},
});
