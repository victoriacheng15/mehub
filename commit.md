# Git Commit Info

## PR Description

### Summary

Grouped GitHub Actions updates, introduced a Go dependency update command, refined the Makefile structure, and added an AI agent guide with Nix-first instructions.

### List of Changes

- **Dependabot**: Grouped all `github-actions` updates into a single PR in `.github/dependabot.yml` to minimize noise.
- **Makefile**:
  - Added `update` target for automated `go get -u` and `go mod tidy`.
  - Reordered `.PHONY` targets to align with their definition order for better maintainability.
  - Updated `help` documentation.
- **Documentation**: Created `AGENTS.md` to provide architectural context and operational guidelines for AI agents, emphasizing reproducible builds via `nix-shell`.

### Verification

- **Dependabot**: Verified YAML syntax.
- **Makefile**: Verified `.PHONY` consistency and `help` output.
- **AGENTS.md**: Verified build instructions (`make setup-tailwind && make nix-build`) match the project's Nix-based workflow.

## Execution Commands

```bash
git checkout -b chore/improve-dx-and-deps
git add .github/dependabot.yml Makefile AGENTS.md commit.md
git commit -m "chore: group actions deps, add AGENTS.md, and refine Makefile"
```
