export const SKILLS = [
	// Core languages
	"Python",
	"Go",
	"TypeScript",
	"JavaScript",
	// Backend & frameworks
	"Node.js",
	"Flask",
	// Frontend & frameworks
	"React",
	"Next.js",
	"Tailwind CSS",
	// Databases
	"PostgreSQL",
	"MongoDB",
	// Platform & DevOps & cloud
	"Grafana",
	"Docker",
	"Linux",
	"Bash",
	"GitHub Actions",
	"Azure",
];

export function formatSkillPath(skill: string) {
	return skill.replace(/\./g, "dot").replace(/\s+/g, "").toLowerCase();
}
