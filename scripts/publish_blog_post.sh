#!/bin/bash
set -euo pipefail

ROOT_DIR=$(pwd)

# Enable test mode if --test flag is passed
if [[ "${1:-}" == "--test" ]]; then
  echo "Test mode enabled: creating temporary draft post..."
  cat <<EOF > blog/temp-draft-test.md
---
title: "Temp Draft Test"
date: 2020-01-01
draft: true
---
This is a test post.
EOF
  trap 'rm -f "$ROOT_DIR/blog/temp-draft-test.md"' EXIT INT TERM
fi

cd "blog" || exit 1

today=$(date -u +%Y-%m-%d)

echo "UTC Time: $(date -u +"%Y-%m-%d %T")"
echo "===================="
echo "Script started"
echo -e "====================\n"

# Find all files with the draft line
draft_files=$(grep -l '^draft:' *.md || true)

if [[ -z "$draft_files" ]]; then
  echo "No draft files found."
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "changed=false" >> "$GITHUB_OUTPUT"
  fi
  exit 0
fi

ready_publish=false

for file in $draft_files; do
  file_date=$(grep '^date:' "$file" | awk '{print $2}' | tr -d '"')

  if [[ "$file_date" < "$today" || "$file_date" == "$today" ]]; then
    echo "Publishing $file..."
    post_title=$(sed -n 's/^title:[[:space:]]*"\(.*\)"/\1/p' "$file")
    sed -i '/^draft:/d' "$file"
    ready_publish=true
    
    if [ -n "${GITHUB_OUTPUT:-}" ]; then
      echo "changed=true" >> "$GITHUB_OUTPUT"
      echo "post_title=$post_title" >> "$GITHUB_OUTPUT"
    fi
    break
  fi
done

if [[ "$ready_publish" = false ]]; then
  echo "No blog post was ready to publish."
  if [ -n "${GITHUB_OUTPUT:-}" ]; then
    echo "changed=false" >> "$GITHUB_OUTPUT"
  fi
  exit 0
fi

echo -e "\n===================="
echo "Script completed."
echo "===================="