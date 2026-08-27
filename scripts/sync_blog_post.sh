#!/bin/bash
set -euo pipefail

# Ensure required env variables are present
if [[ -z "${BLOG_TITLE:-}" || -z "${API_LINK:-}" || -z "${SYNC_BLOG_TOKEN:-}" ]]; then
  echo "Error: Required environment variables BLOG_TITLE, API_LINK, and SYNC_BLOG_TOKEN must be set." >&2
  exit 1
fi

# ==== blog title -> slug ====
SLUG=$(echo "$BLOG_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//;s/-$//')

# ==== Logs ====
echo "Syncing blog post..."
echo "Issue ID: ${ID:-}"
echo "Title: $BLOG_TITLE"
echo "Slug: $SLUG"

# ==== Fetch the blog post content from the API ====
echo "Fetching blog post content from API..."
curl -s \
  -H "Authorization: token $SYNC_BLOG_TOKEN" \
  -H "Accept: application/vnd.github.v3+json" \
  "$API_LINK" | \
  jq '.[0].body' | \
  sed -e 's/^"//' | \
  sed -e 's/"$//' | \
  sed -e 's/\\n/\n/g' | \
  sed -e 's/\\"/"/g' | \
  sed -e 's/^date: "\(.*\)"/date: \1/' > blog/"$SLUG".md

# ==== Set GITHUB_OUTPUT variables ====
if [ -n "${GITHUB_OUTPUT:-}" ]; then
  echo "slug=$SLUG" >> "$GITHUB_OUTPUT"
  echo "post_title=$BLOG_TITLE" >> "$GITHUB_OUTPUT"
fi