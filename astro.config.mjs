import { defineConfig } from "astro/config";
import { siteConfig } from "./src/config";
import rehypeMermaid from 'rehype-mermaid';

import tailwindcss from "@tailwindcss/vite";

import expressiveCode from "astro-expressive-code";

export default defineConfig({
	site: siteConfig.site,
	integrations: [expressiveCode()],
	markdown: {
		syntaxHighlight: {
			type: "shiki",
		},
		rehypePlugins: [
			[
				rehypeMermaid,
				{
					mermaidConfig: {
						theme: "default",
					},
				},
			],
		],
	},
	vite: {
		plugins: [tailwindcss()],
	},
});