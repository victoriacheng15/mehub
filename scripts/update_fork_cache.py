#!/usr/bin/env python3
"""
Update the scripts/fork_cache.json file by querying the GitHub API.
Resolves active (non-archived) forks owned by the user and finds their upstream parents.
Uses only the Python standard library.
"""

import datetime
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


def save_fork_cache(cache: dict[str, str]) -> None:
    """Save the fork cache to CACHE_PATH, updating the date only if contents changed."""
    # Sort the cache by fork name to ensure deterministic comparison and output format
    cache = {k: cache[k] for k in sorted(cache.keys())}

    existing_forks = {}
    existing_date = None
    if os.path.exists(CACHE_PATH):
        try:
            with open(CACHE_PATH, "r", encoding="utf-8") as f:
                data = json.load(f)
            if isinstance(data, dict):
                if "forks" in data and "updated_date" in data:
                    existing_forks = data["forks"]
                    existing_date = data["updated_date"]
                else:
                    existing_forks = data
        except Exception as err:
            print(f"Warning: Failed to read existing cache: {err}", file=sys.stderr)

    if cache == existing_forks and existing_date is not None:
        updated_date = existing_date
        is_updated = False
    else:
        updated_date = datetime.date.today().isoformat()
        is_updated = True

    output_data = {
        "updated_date": updated_date,
        "forks": cache
    }

    with open(CACHE_PATH, "w", encoding="utf-8") as f:
        json.dump(output_data, f, indent=2)

    if is_updated:
        print(f"Successfully updated cache with {len(cache)} active parents in {CACHE_PATH} (updated_date: {updated_date})", file=sys.stderr)
    else:
        print(f"No changes detected in active fork repositories. Cache not updated (kept updated_date: {updated_date})", file=sys.stderr)


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

    save_fork_cache(cache)


if __name__ == "__main__":
    try:
        main()
    except KeyboardInterrupt:
        sys.exit(130)
