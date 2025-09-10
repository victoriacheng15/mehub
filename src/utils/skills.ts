export const SKILLS = [
	// Core languages
	"Python",
	"Go",
	"Ruby",
	"TypeScript",
	"JavaScript",
	// Backend & frameworks
	"Node.js",
	"Express",
	"Flask",
	// Frontend & frameworks
	"React",
	"Next.js",
	"Tailwind CSS",
	// Databases
	"PostgreSQL",
	"MongoDB",
	// Platform
	"Docker",
	"Linux",
	"Bash",
	"GitHub Actions",
	"Ansible",
	"Azure",
];

export function formatSkillPath(skill: string) {
	return skill.replace(/\./g, "dot").replace(/\s+/g, "").toLowerCase();
}
