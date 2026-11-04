#!/usr/bin/env python3
"""
Audit blog post tags.

- Extracts tags dynamically from all blog markdown files.
- Validates tag naming format (lowercase kebab-case).
- Detects spelling typos and near-duplicate tag collisions across the corpus.
- Validates that every engineering log post carries the 'retrospective' tag.
- Validates that no non-engineering-log post carries the 'retrospective' tag.
- Validates that no post exceeds MAX_TAGS tags.
- Exits with code 1 if any violations are found.
"""

from __future__ import annotations

from collections import Counter
import difflib
import os
import re
import sys

BLOG_DIR = os.path.join(os.path.dirname(__file__), "..", "blog")
ENGINEERING_LOG_PREFIX = "engineering-log-"
REQUIRED_TAG = "retrospective"
MAX_TAGS = 3
TYPO_CUTOFF = 0.8

TAG_PATTERN = re.compile(r'^tags:\s*\[([^\]]+)\]', re.MULTILINE)
SINGLE_TAG = re.compile(r'"([^"]+)"')
KEBAB_CASE = re.compile(r'^[a-z0-9]+(-[a-z0-9]+)*$')


def parse_tags(content: str) -> list[str]:
    match = TAG_PATTERN.search(content)
    if not match:
        return []
    return SINGLE_TAG.findall(match.group(1))


def check_tag_spelling(
    tag: str,
    tag_counts: Counter[str],
    all_tags: set[str],
) -> list[str]:
    """Validate tag format and check for near-duplicate typos across all blog tags."""
    issues = []

    # Format check: must be lowercase kebab-case.
    if not KEBAB_CASE.match(tag):
        suggested = re.sub(r'[^a-z0-9]+', '-', tag.lower()).strip('-')
        if suggested:
            issues.append(f"tag '{tag}' must be lowercase kebab-case (did you mean '{suggested}'?)")
        else:
            issues.append(f"tag '{tag}' is invalid: must be lowercase alphanumeric with hyphens")
        return issues

    # Fuzzy typo check: compare against other tags in the blog corpus.
    candidates = [t for t in all_tags if t != tag and tag_counts[t] >= tag_counts[tag]]
    matches = difflib.get_close_matches(tag, candidates, n=1, cutoff=TYPO_CUTOFF)
    if matches:
        issues.append(f"possible typo in tag '{tag}' (did you mean '{matches[0]}'?)")

    return issues


def run_checks(
    filename: str,
    tags: list[str],
    tag_counts: Counter[str],
    all_tags: set[str],
) -> list[str]:
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

    # Spelling and format checks.
    for tag in tags:
        issues.extend(check_tag_spelling(tag, tag_counts, all_tags))

    return issues


def main() -> int:
    all_tags_list: list[str] = []
    parsed_posts: list[tuple[str, list[str]]] = []
    violations: dict[str, list[str]] = {}

    posts = sorted(f for f in os.listdir(BLOG_DIR) if f.endswith(".md"))

    # Pass 1: Dynamically extract all tags from the blog markdown files.
    for filename in posts:
        path = os.path.join(BLOG_DIR, filename)
        with open(path, encoding="utf-8") as fh:
            content = fh.read()

        tags = parse_tags(content)
        all_tags_list.extend(tags)
        parsed_posts.append((filename, tags))

    tag_counts = Counter(all_tags_list)
    all_tags = set(all_tags_list)

    # Pass 2: Run validation checks across posts and tag corpus.
    for filename, tags in parsed_posts:
        issues = run_checks(filename, tags, tag_counts, all_tags)
        if issues:
            violations[filename] = issues

    print("=== All tags ===")
    print(sorted(all_tags))
    print(f"\nTotal unique tags: {len(all_tags)}")

    checks = [
        f"engineering log posts carry '{REQUIRED_TAG}'",
        f"no non-engineering-log posts carry '{REQUIRED_TAG}'",
        f"no post exceeds {MAX_TAGS} tags",
        "all tags follow lowercase kebab-case format",
        "no spelling typos or near-duplicate tag collisions detected",
    ]

    print("\n=== Tag audit ===")
    if violations:
        print(f"FAIL: {len(violations)} post(s) with violations:\n")
        for filename, issues in violations.items():
            print(f"  {filename}")
            for issue in issues:
                print(f"    - {issue}")
        return 1

    for check in checks:
        print(f"OK: {check}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
