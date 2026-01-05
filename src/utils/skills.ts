export const SKILLS = [
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
];

export function formatSkillPath(skill: string) {
	return skill.replace(/\./g, "dot").replace(/\s+/g, "").toLowerCase();
}
