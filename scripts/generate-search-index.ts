import fs from "fs";
import path from "path";

interface BlogPost {
	title: string;
	slug: string;
	summary: string;
	tags: string[];
	date: string;
	draft: boolean;
}

interface FrontmatterData {
	title?: string;
	description?: string;
	date?: string;
	tags?: string[];
	draft?: boolean;
}

// Simple YAML frontmatter parser
function parseFrontmatter(content: string): FrontmatterData {
	const match = content.match(/^---\n([\s\S]*?)\n---/);
	if (!match) return {};

	const frontmatter = match[1];
	const data: FrontmatterData = {};

	// Parse title
	const titleMatch = frontmatter.match(/title:\s*["'](.+?)["']/);
	if (titleMatch) data.title = titleMatch[1];

	// Parse description
	const descMatch = frontmatter.match(/description:\s*["'](.+?)["']/);
	if (descMatch) data.description = descMatch[1];

	// Parse date
	const dateMatch = frontmatter.match(/date:\s*(\d{4}-\d{2}-\d{2})/);
	if (dateMatch) data.date = dateMatch[1];

	// Parse tags (handles both array formats)
	const tagsMatch = frontmatter.match(/tags:\s*\[(.*?)\]/);
	if (tagsMatch) {
		data.tags = tagsMatch[1]
			.split(",")
			.map((tag) => tag.trim().replace(/["']/g, ""))
			.filter((tag) => tag);
	} else {
		data.tags = [];
	}

	// Parse draft
	const draftMatch = frontmatter.match(/draft:\s*(true|false)/);
	if (draftMatch) data.draft = draftMatch[1] === "true";

	return data;
}

// Read all markdown files from src/content/blog
const blogDir = "./src/content/blog";
const files = fs.readdirSync(blogDir).filter((file) => file.endsWith(".md"));

const index: BlogPost[] = files
	.map((file) => {
		const filePath = path.join(blogDir, file);
		const content = fs.readFileSync(filePath, "utf-8");
		const data = parseFrontmatter(content);

		return {
			title: data.title || file.replace(/\.md$/, ""),
			slug: file.replace(/\.md$/, ""),
			summary: data.description || "",
			tags: data.tags || [],
			date: data.date || "",
			draft: data.draft || false,
		};
	})
	.filter((post) => !post.draft)
	.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());

// Ensure public directory exists
if (!fs.existsSync("./public")) {
	fs.mkdirSync("./public", { recursive: true });
}

// Write the search index
fs.writeFileSync("./public/search-index.json", JSON.stringify(index, null, 2));
console.log(`✓ Generated search index with ${index.length} posts`);
