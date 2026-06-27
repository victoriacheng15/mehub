#!/usr/bin/env python3
"""
Fetch external GitHub contributions and merge them into projects.yaml.

Queries the GitHub REST API for external pull requests, resolves repository
metadata, and outputs the updated layout to projects.yaml.
Uses only the Python standard library (no pip dependencies required).
"""

from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.parse
import urllib.request

# Configuration settings
USERNAME = os.environ.get("GITHUB_USER", "victoriacheng15")
API_TOKEN = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
REPOS = ["chaos-mesh/chaos-mesh"]

# API base headers (preserving custom User-Agent)
HEADERS = {
    "Accept": "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
    "User-Agent": f"fetch-contribution-{USERNAME}",
}
if API_TOKEN:
    HEADERS["Authorization"] = f"Bearer {API_TOKEN}"

# File location configurations
SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
PROJECTS_YAML_PATH = os.path.join(SCRIPT_DIR, "..", "internal", "templates", "contents", "projects.yaml")


def query_github_api(url: str) -> dict:
    """Send a GET request to the GitHub API and return parsed JSON data."""
    request = urllib.request.Request(url, headers=HEADERS)
    try:
        with urllib.request.urlopen(request, timeout=30) as response:
            return json.loads(response.read().decode("utf-8"))
    except urllib.error.HTTPError as err:
        error_details = err.read().decode("utf-8", errors="replace")
        raise SystemExit(f"GitHub API HTTP Error [{err.code}]: {error_details}")
    except urllib.error.URLError as err:
        raise SystemExit(f"GitHub API Network Connection Error: {err.reason}")


def fetch_external_pull_requests() -> list[dict]:
    """Retrieve all external pull requests authored by the user from GitHub API."""
    search_query = urllib.parse.quote(f"author:{USERNAME} type:pr -user:{USERNAME}")
    pull_requests = []
    
    # Iterate through pagination pages (max 10 pages)
    for page in range(1, 11):
        url = f"https://api.github.com/search/issues?q={search_query}&per_page=100&page={page}"
        payload = query_github_api(url)
        page_items = payload.get("items", [])
        pull_requests.extend(page_items)
        if len(page_items) < 100:
            break
            
    return pull_requests


def fetch_repo_description(repo_name: str) -> str:
    """Retrieve the default description of a repository from GitHub API."""
    url = f"https://api.github.com/repos/{repo_name}"
    try:
        data = query_github_api(url)
        return data.get("description", "") or ""
    except Exception as err:
        print(f"Warning: Failed to fetch description for {repo_name}: {err}", file=sys.stderr)
        return ""


def organize_contributions(raw_prs: list[dict]) -> dict[str, list[dict]]:
    """Group pull requests by repository name and filter by featured list."""
    contributions = {}
    for item in raw_prs:
        # Extract owner/repo name from repository URL
        repo_url = item["repository_url"]
        repo_name = repo_url.split("/repos/")[1]
        
        if REPOS and repo_name not in REPOS:
            continue
            
        contributions.setdefault(repo_name, []).append({
            "number": item["number"],
            "title": item["title"],
        })
        
    # Sort pull requests in descending order by number (latest first)
    for repo in contributions:
        contributions[repo].sort(key=lambda x: x["number"], reverse=True)
        
    return contributions


def parse_existing_contributions(yaml_text: str) -> dict[str, dict[str, str]]:
    """Parse existing contributions from projects.yaml text to preserve descriptions and links."""
    existing = {}
    current_repo = None
    current_data = {}
    in_contributions = False
    
    for line in yaml_text.splitlines():
        if line.startswith("contributions:"):
            in_contributions = True
            continue
        if not in_contributions:
            continue
            
        stripped = line.strip()
        if not stripped or ":" not in stripped:
            continue
            
        key, _, val = stripped.partition(":")
        key = key.strip()
        val = val.strip()
        
        if key == "- repo":
            if current_repo and current_data:
                existing[current_repo] = current_data
            current_repo = val
            current_data = {}
        elif key == "link" and current_repo:
            current_data["link"] = val.strip('"\'')
        elif key == "description" and current_repo:
            if (val.startswith('"') and val.endswith('"')) or (val.startswith("'") and val.endswith("'")):
                val = val[1:-1]
            current_data["description"] = val.replace('\\"', '"').replace('\\\\', '\\')
                
    if current_repo and current_data:
        existing[current_repo] = current_data
        
    return existing


def strip_existing_contributions(yaml_text: str) -> str:
    """Strip the existing contributions block from projects.yaml."""
    lines = yaml_text.splitlines()
    result = []
    for line in lines:
        if line.startswith("contributions:"):
            break
        result.append(line)
    return "\n".join(result).rstrip() + "\n"


def escape_yaml_string(s: str) -> str:
    """Escape backslashes and double quotes for clean double-quoted YAML output."""
    return s.replace('\\', '\\\\').replace('"', '\\"').replace('\n', ' ')


def generate_contributions_yaml(contributions: list[dict]) -> str:
    """Convert contributions list to formatted YAML string."""
    yaml_lines = ["contributions:"]
    for contrib in contributions:
        yaml_lines.append(f"  - repo: {contrib['repo']}")
        yaml_lines.append(f"    link: {contrib['link']}")
        escaped_desc = escape_yaml_string(contrib['description'])
        yaml_lines.append(f"    description: \"{escaped_desc}\"")
        yaml_lines.append("    prs:")
        for pr in contrib["prs"]:
            yaml_lines.append(f"      - number: {pr['number']}")
            escaped_title = escape_yaml_string(pr['title'])
            yaml_lines.append(f"        title: \"{escaped_title}\"")
    return "\n".join(yaml_lines) + "\n"


def main() -> None:
    print("Fetching external contributions from GitHub...", file=sys.stderr)
    raw_prs = fetch_external_pull_requests()
    contributions_map = organize_contributions(raw_prs)

    if not os.path.exists(PROJECTS_YAML_PATH):
        print(f"Error: {PROJECTS_YAML_PATH} not found.", file=sys.stderr)
        sys.exit(1)

    with open(PROJECTS_YAML_PATH, "r", encoding="utf-8") as f:
        yaml_text = f.read()

    existing_meta = parse_existing_contributions(yaml_text)
    updated_contributions = []

    # Map fetched contributions and merge existing descriptions
    for repo_name, prs in sorted(contributions_map.items()):
        meta = existing_meta.get(repo_name, {})
        description = meta.get("description")
        
        if not description:
            print(f"Fetching metadata for {repo_name}...", file=sys.stderr)
            description = fetch_repo_description(repo_name)

        updated_contrib = {
            "repo": repo_name,
            "link": meta.get("link") or f"https://github.com/{repo_name}",
            "description": description,
            "prs": prs
        }
        updated_contributions.append(updated_contrib)

    # Maintain existing manual contributions that had no fetched PRs
    for repo_name, meta in existing_meta.items():
        if repo_name not in contributions_map:
            updated_contrib = {
                "repo": repo_name,
                "link": meta.get("link") or f"https://github.com/{repo_name}",
                "description": meta.get("description") or "",
                "prs": []
            }
            updated_contributions.append(updated_contrib)

    # Assemble updated projects.yaml
    clean_yaml = strip_existing_contributions(yaml_text)
    if updated_contributions:
        contributions_yaml = generate_contributions_yaml(updated_contributions)
        final_yaml = f"{clean_yaml}\n{contributions_yaml}"
    else:
        final_yaml = clean_yaml

    with open(PROJECTS_YAML_PATH, "w", encoding="utf-8") as f:
        f.write(final_yaml)

    print("Successfully updated projects.yaml with latest contributions.", file=sys.stderr)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(130)
