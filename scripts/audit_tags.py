#!/usr/bin/env python3
"""
Audit blog post tags.

- Prints all unique tags across the blog directory.
- Validates that every engineering log post carries the 'retrospective' tag.
- Validates that no non-engineering-log post carries the 'retrospective' tag.
- Validates that no post exceeds MAX_TAGS tags.
- Validates that all tags belong to the known tag list.
- Exits with code 1 if any violations are found.
"""

from __future__ import annotations

import os
import re
import sys

BLOG_DIR = os.path.join(os.path.dirname(__file__), "..", "blog")
ENGINEERING_LOG_PREFIX = "engineering-log-"
REQUIRED_TAG = "retrospective"
MAX_TAGS = 3

KNOWN_TAGS = {
    "backend",
    "cloud",
    "cncf",
    "data-structure",
    "docker",
    "frontend",
    "go",
    "growth",
    "javascript",
    "kubernetes",
    "linux",
    "mcp",
    "monthly-log",
    "observability",
    "platform",
    "python",
    "retrospective",
    "sre",
    "system-design",
    "terraform",
    "typescript",
}

TAG_PATTERN = re.compile(r'^tags:\s*\[([^\]]+)\]', re.MULTILINE)
SINGLE_TAG = re.compile(r'"([^"]+)"')


def parse_tags(content: str) -> list[str]:
    match = TAG_PATTERN.search(content)
    if not match:
        return []
    return SINGLE_TAG.findall(match.group(1))


def run_checks(filename: str, tags: list[str]) -> list[str]:
    """Return a list of violation strings for a single post."""
    issues = []

    # Engineering log must carry retrospective; others must not.
    if filename.startswith(ENGINEERING_LOG_PREFIX):
        if REQUIRED_TAG not in tags:
            issues.append(f"missing '{REQUIRED_TAG}' tag")
    else:
        if REQUIRED_TAG in tags:
            issues.append(f"carries '{REQUIRED_TAG}' but is not an engineering log")

    # Tag count ceiling.
    if len(tags) > MAX_TAGS:
        issues.append(f"exceeds {MAX_TAGS} tags ({len(tags)} found: {tags})")

    # Unknown tags.
    unknown = sorted(set(tags) - KNOWN_TAGS)
    if unknown:
        issues.append(f"unknown tag(s): {unknown}")

    return issues


def main() -> int:
    all_tags: set[str] = set()
    violations: dict[str, list[str]] = {}

    posts = sorted(f for f in os.listdir(BLOG_DIR) if f.endswith(".md"))

    for filename in posts:
        path = os.path.join(BLOG_DIR, filename)
        with open(path, encoding="utf-8") as fh:
            content = fh.read()

        tags = parse_tags(content)
        all_tags.update(tags)

        issues = run_checks(filename, tags)
        if issues:
            violations[filename] = issues

    print("=== All tags ===")
    print(sorted(all_tags))
    print(f"\nTotal unique tags: {len(all_tags)}")

    checks = [
        f"engineering log posts carry '{REQUIRED_TAG}'",
        f"no non-engineering-log posts carry '{REQUIRED_TAG}'",
        f"no post exceeds {MAX_TAGS} tags",
        "all tags are known",
    ]

    print("\n=== Tag audit ===")
    if violations:
        print(f"FAIL — {len(violations)} post(s) with violations:\n")
        for filename, issues in violations.items():
            print(f"  {filename}")
            for issue in issues:
                print(f"    - {issue}")
        return 1

    for check in checks:
        print(f"OK — {check}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
