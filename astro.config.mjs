import { defineConfig } from "astro/config";
import { siteConfig } from "./src/config";
import expressiveCode from "astro-expressive-code";
import rehypeMermaid from 'rehype-mermaid';
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
	site: siteConfig.site,
	integrations: [expressiveCode({
		theme: ["catppuccin-mocha", "catppuccin-latte"],
		defaultProps: {
			wrap: true,
		}
	})],
	markdown: {
		syntaxHighlight: {
      excludeLangs: ["mermaid"]
    },
		rehypePlugins: [[rehypeMermaid, { strategy: "img-svg", dark: true, colorScheme: "forest" }]],
	},
	vite: {
		plugins: [tailwindcss()],
	},
});
