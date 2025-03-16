export const SKILLS = [
	"TypeScript",
	"JavaScript",
	"Python",
	".NET",
	"React",
	"Next.js",
	"Redux",
	"React Native",
	"Tailwind CSS",
	"SASS",
	"HTML",
	"CSS",
	"Node.js",
	"Express",
	"Flask",
	"MongoDB",
	"PostgreSQL",
	"Cypress",
	"Jest",
	"Testing Library",
	"Github Actions",
	"Docker",
	"Linux",
	"Bash",
	"Ansible",
	"Azure",
];

export function formatSkillPath(skill: string) {
	return skill.replace(/\./g, "dot").replace(/\s+/g, "").toLowerCase();
}
