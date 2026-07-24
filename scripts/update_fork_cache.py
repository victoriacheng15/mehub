#!/usr/bin/env python3
"""
Update the scripts/fork_cache.json file by querying the GitHub API.
Resolves active (non-archived) forks owned by the user and finds their upstream parents.
Uses only the Python standard library.
"""

import json
import os
import sys
import urllib.error
import urllib.request

USERNAME = os.environ.get("GITHUB_USER", "victoriacheng15")
API_TOKEN = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")

HEADERS = {
    "Accept": "application/vnd.github+json",
    "X-GitHub-Api-Version": "2022-11-28",
    "User-Agent": f"update-fork-cache-{USERNAME}",
}
if API_TOKEN:
    HEADERS["Authorization"] = f"Bearer {API_TOKEN}"

SCRIPT_DIR = os.path.dirname(os.path.abspath(__file__))
CACHE_PATH = os.path.join(SCRIPT_DIR, "fork_cache.json")


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


def main() -> None:
    print("Fetching active fork repositories from GitHub...", file=sys.stderr)
    url = f"https://api.github.com/users/{USERNAME}/repos?per_page=100"
    payload = query_github_api(url)
    if not isinstance(payload, list):
        print("Error: Expected a list of repositories from GitHub API.", file=sys.stderr)
        sys.exit(1)

    cache = {}
    for repo in payload:
        if repo.get("fork") and not repo.get("archived"):
            repo_name = repo["full_name"]
            print(f"Resolving parent for {repo_name}...", file=sys.stderr)
            detail_url = f"https://api.github.com/repos/{repo_name}"
            try:
                detail = query_github_api(detail_url)
                parent = detail.get("parent", {}).get("full_name")
                if parent:
                    cache[repo_name] = parent
            except Exception as err:
                print(f"Warning: Failed to fetch parent for {repo_name}: {err}", file=sys.stderr)

    with open(CACHE_PATH, "w", encoding="utf-8") as f:
        json.dump(cache, f, indent=2)
    print(f"Successfully updated cache with {len(cache)} active parents in {CACHE_PATH}", file=sys.stderr)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(130)
