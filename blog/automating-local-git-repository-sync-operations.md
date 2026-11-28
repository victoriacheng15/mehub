---
title: "Automating Local Git Repository Sync Operations"
description: "How a custom Python script automates fetching and rebasing local git repositories to eliminate manual folder switching and handle branch naming conventions automatically."
date: 2026-12-01
tags: ["python", "automation"]
draft: true
---

## The Workspace Sync Problem

Managing dozens of local git repositories across an active workspace makes it easy to forget which folders are actually up to date. Keeping forks in sync required manually changing into each directory, checking branch names, and running fetch and rebase commands one by one. Outdated local branches go unnoticed until opening an older project leads to avoidable merge conflicts.

I wrote a standalone Python automation script to handle workspace updates in a single pass. The tool iterates through local repositories, resolves remotes dynamically, and synchronizes default branches without manual intervention.

---

## Sync Engine Design

The sync engine automates remote checks and local rebases by executing git commands in target directories. The logic queries remote definitions and prioritizes upstream sources over origin destinations to keep fork repositories aligned with original upstreams. Additionally, the process dynamically detects default branch configurations like `main` or `master` via symbolic remote HEAD references. This architecture resolves distinct repository conventions without requiring hardcoded configuration mappings.

```text
           +--------------------+
           |  Target Directory  |
           +--------------------+
                     |
                     v
           +--------------------+
           |  Clean Tree Check  |
           +--------------------+
            /                  \\
      (No) /                    \\ (Yes)
          v                      v
  +---------------+      +--------------------+
  | Skip Repo     |      |   Resolve Remote   |
  +---------------+      +--------------------+
                          /                  \\
                    (No) /                    \\ (Yes: Upstream/Origin)
                        v                      v
                +---------------+      +--------------------+
                | Fail Repo     |      |   Fetch & Rebase   |
                +---------------+      +--------------------+
                                                 |
                                                 v
                                       +--------------------+
                                       |   Success Report   |
                                       +--------------------+
```

---

## The Sync Script

The complete Python automation script is presented below to serve as a reference implementation. It requires no external library dependencies and executes standard git commands via the subprocess module. This utility can be saved directly into a local scripts directory and run from the command line.

<details>
<summary>Click to view sync_repos.py</summary>

```python
#!/usr/bin/env python3
"""
Syncs local git repositories with their upstream or origin remotes using rebase.
Supports dry-run simulation and targeting single repositories or directories.

Usage Examples:
    # 1. Dry run across all repositories in the parent workspace (safe simulation):
    python3 scripts/sync_repos.py --dry-run
    python3 scripts/sync_repos.py -n

    # 2. Sync all repositories in the parent workspace:
    python3 scripts/sync_repos.py

    # 3. Sync a single repository by folder name:
    python3 scripts/sync_repos.py -p [repo-name]
    python3 scripts/sync_repos.py --path [repo-name]

    # 4. Sync a single repository by relative or absolute path:
    python3 scripts/sync_repos.py -p ../[repo-name]
    python3 scripts/sync_repos.py -p /[your-path]/[repo-name]

    # 5. Dry run against a single repository:
    python3 scripts/sync_repos.py -p [repo-name] -n
"""

import argparse
from pathlib import Path
import subprocess
import sys


def run_git(repo_path: Path, *args: str) -> subprocess.CompletedProcess:
    """Executes a git command inside the target repository directory."""
    return subprocess.run(
        ["git", "-C", str(repo_path), *args],
        capture_output=True,
        text=True,
        check=False,
    )


def get_remote(repo_path: Path) -> str | None:
    """Determines preferred remote, prioritizing upstream over origin."""
    res = run_git(repo_path, "remote")
    remotes = res.stdout.split()
    if "upstream" in remotes:
        return "upstream"
    if "origin" in remotes:
        return "origin"
    return None


def get_default_branch(repo_path: Path, remote: str) -> str | None:
    """Detects default branch name via remote HEAD or standard branch names."""
    res = run_git(repo_path, "symbolic-ref", f"refs/remotes/{remote}/HEAD")
    if res.returncode == 0:
        return res.stdout.strip().split("/")[-1]

    branches = run_git(repo_path, "branch", "--list", "main", "master").stdout.split()
    for candidate in ("main", "master"):
        if candidate in branches or f"*{candidate}" in branches:
            return candidate
    return None


def is_clean(repo_path: Path) -> bool:
    """
    Returns True if the working tree has no staged or modified tracked files.
    Untracked files (e.g. local scratch files/notes) are ignored.
    """
    res = run_git(repo_path, "status", "--porcelain", "--untracked-files=no")
    return res.returncode == 0 and not res.stdout.strip()


def sync_repo(repo_path: Path, dry_run: bool = False) -> tuple[str, str]:
    """Fetches and rebases the default branch of a single repository."""
    if not (repo_path / ".git").is_dir():
        return "SKIPPED", "Not a git repository"

    if not is_clean(repo_path):
        return "SKIPPED", "Tracked file modifications or staged changes present"

    remote = get_remote(repo_path)
    if not remote:
        return "FAILED", "No origin or upstream remote found"

    branch = get_default_branch(repo_path, remote)
    if not branch:
        return "FAILED", "Could not determine default branch (main/master)"

    if dry_run:
        return "DRY-RUN", f"Would fetch {remote} and rebase {branch} on {remote}/{branch}"

    fetch_res = run_git(repo_path, "fetch", remote)
    if fetch_res.returncode != 0:
        return "FAILED", f"Fetch failed: {fetch_res.stderr.strip()}"

    checkout_res = run_git(repo_path, "checkout", branch)
    if checkout_res.returncode != 0:
        return "FAILED", f"Checkout failed: {checkout_res.stderr.strip()}"

    rebase_res = run_git(repo_path, "rebase", f"{remote}/{branch}")
    if rebase_res.returncode != 0:
        run_git(repo_path, "rebase", "--abort")
        return "FAILED", f"Rebase failed (aborted): {rebase_res.stderr.strip()}"

    return "UPDATED", f"Rebased {branch} on {remote}/{branch}"


def resolve_targets(target_input: str | None, default_parent: Path) -> list[Path]:
    """Resolves target input into a list of repository paths."""
    if not target_input:
        target_dir = default_parent
    else:
        path_candidate = Path(target_input).expanduser().resolve()
        if path_candidate.exists():
            target_dir = path_candidate
        else:
            named_candidate = default_parent / target_input
            if named_candidate.exists():
                target_dir = named_candidate
            else:
                print(f"Error: Target path '{target_input}' does not exist.")
                sys.exit(1)

    # If the target itself is a git repository
    if (target_dir / ".git").is_dir():
        return [target_dir]

    # If the target is a parent directory containing repos
    if target_dir.is_dir():
        repos = [p for p in sorted(target_dir.iterdir()) if p.is_dir() and (p / ".git").is_dir()]
        return repos

    print(f"Error: '{target_dir}' is not a valid git repository or directory.")
    sys.exit(1)


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Sync local git repositories with upstream or origin using rebase."
    )
    parser.add_argument(
        "-p",
        "--path",
        dest="target",
        help="Path to a single repository, repository folder name, or directory of repositories.",
        default=None,
    )
    parser.add_argument(
        "-n",
        "--dry-run",
        action="store_true",
        help="Simulate updates without executing fetch, checkout, or rebase.",
    )
    return parser.parse_args()


def find_default_workspace(script_path: Path) -> Path:
    """
    Dynamically finds the workspace directory containing multiple git repositories:
    1. If the current working directory contains git repos, use CWD.
    2. Traverse up from the script location to find the enclosing git repo root,
       then return that repo's parent directory (regardless of nesting depth).
    3. Fallback to CWD.
    """
    cwd = Path.cwd().resolve()
    if any(p.is_dir() and (p / ".git").is_dir() for p in cwd.iterdir()):
        return cwd

    current = script_path.resolve().parent
    for parent in [current, *current.parents]:
        if (parent / ".git").is_dir():
            return parent.parent

    return cwd


def main():
    args = parse_args()
    default_dir = find_default_workspace(Path(__file__))
    repos = resolve_targets(args.target, default_dir)

    if not repos:
        print("No git repositories found to process.")
        sys.exit(0)

    mode_label = " (DRY-RUN)" if args.dry_run else ""
    print(f"Found {len(repos)} repository(ies) to process{mode_label}.\
")

    results: dict[str, list[tuple[str, str]]] = {
        "UPDATED": [],
        "DRY-RUN": [],
        "SKIPPED": [],
        "FAILED": [],
    }

    for repo in repos:
        status, message = sync_repo(repo, dry_run=args.dry_run)
        results[status].append((repo.name, message))
        print(f"[{status:<7}] {repo.name}: {message}")

    print("\
--- Summary ---")
    summary_parts = [f"Total: {len(repos)}"]
    for status_key in ("UPDATED", "DRY-RUN", "SKIPPED", "FAILED"):
        if results[status_key]:
            summary_parts.append(f"{status_key.capitalize()}: {len(results[status_key])}")
    print(" | ".join(summary_parts))


if __name__ == "__main__":
    main()
```
</details>

---

## Designing for Operational Safety

Automating rebases can corrupt branch histories or lose uncommitted work if safety boundaries are omitted. To ensure non-destructive execution across dozens of folders, the script enforces three strict operational guardrails:

- Verifying that a target folder contains a valid `.git` directory before invoking commands.
- Skipping repositories with uncommitted tracked changes or staged index entries.
- Aborting failed rebases immediately to prevent repositories from getting stuck in detached states.

The cleanliness check runs `git status --porcelain --untracked-files=no`. This allows untracked local scratch files and temporary notes to remain while halting if tracked source files are modified. When a rebase fails due to upstream conflicts, the script executes `git rebase --abort` immediately, returning the working tree to its exact starting state without human intervention.

---

## Conclusion

Remembering whether a repository uses `main` or `master` was one of the most tedious parts of manual syncing. Checking each project's branch name before repeating the same fetch and rebase commands created constant friction. Automating that detection removes the guesswork completely, letting the script handle branch resolution in the background.

Keeping the script in a designated scripts folder allows it to sweep through all repositories grouped under a parent workspace in a single run. With built-in guardrails to skip uncommitted changes and abort conflicted rebases, running the batch sync at any time keeps every local project updated without manual maintenance.
