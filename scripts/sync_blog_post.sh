#!/bin/bash

set -e

# ==== blog title -> slug ====
SLUG=$(echo "$BLOG_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | sed 's/^-//;s/-$//')

# ==== Git config ====
git config --local user.name "$AUTHOR_NAME"
git config --local user.email "$AUTHOR_EMAIL"

# ==== Logs ====
echo "Syncing blog post..."
echo "Issue ID: $ID"
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
  sed -e 's/^date: "\(.*\)"/date: \1/' > src/content/blog/$SLUG.md

# ==== Branch, commit, and push ====
git switch -c "blog/$SLUG"
git add src/content/blog/$SLUG.md
git commit -m "Sync blog post: $BLOG_TITLE"
git push --force-with-lease origin "blog/$SLUG"

gh pr create \
  --base main \
  --head "blog/$SLUG" \
  --title "Sync Blog Post: $BLOG_TITLE" \
  --body "This PR syncs the blog post titled **$BLOG_TITLE** with the content from the API. 🎉" \
  --label "blog"