#!/bin/bash

YELLOW='\033[1;33m'
NC='\033[0m' 

debug_mode=false

debug_log() {
  if $debug_mode; then
    echo "DEBUG: $1"
  fi
}

run_if_not_debug() {
  if [ "$debug_mode" != true ]; then
    eval "$1"
  else
    echo -e "${YELLOW}Skipping -> $1${NC}"
  fi
}

ROOT_DIR=$(pwd)
cd "src/content/blog" || exit

today=$(date -u +%Y-%m-%d)

echo "UTC Time: $(date -u +"%Y-%m-%d %T")"
echo "===================="
echo "Script started"
echo -e "====================\n"

# Find all files with the draft line
draft_files=$(grep -l '^draft:' *.md)

if [[ -z "$draft_files" ]]; then
  echo "No draft files found."
  exit 0
fi

ready_publish=false

for file in $draft_files; do
  file_date=$(grep '^date:' "$file" | awk '{print $2}' | tr -d '"')

  if [[ "$file_date" == "$today" ]]; then
    echo "Publishing $file..."
    sed -i '/^draft:/d' "$file"
    ready_publish=true
  else
    echo -e "\nNot time yet, Skipping $file..."
  fi
done

if [[ "$ready_publish" = false ]]; then
  echo "No blog post was ready to publish."
  exit 0
fi

cd "$ROOT_DIR"

branch_name="publish-post-$(date +%s)"

debug_log "Creating branch: $branch_name"
run_if_not_debug "git switch -c \"$branch_name\""
run_if_not_debug "git config --local user.email \"\$AUTHOR_EMAIL\""
run_if_not_debug "git config --local user.name \"\$AUTHOR_NAME\""

run_if_not_debug "git add src/content/blog"

post_file=$(git diff --name-only HEAD~1 src/content/ | head -n 1)
post_title=$(sed -n 's/^title:[[:space:]]*"\(.*\)"/\1/p' "$post_file")

debug_log "Post file: $post_file"
debug_log "Post title: $post_title"

run_if_not_debug "git commit -m \"chore: publish post $post_title ! 🎉\""
run_if_not_debug "git push --force-with-lease origin $branch_name"

run_if_not_debug "gh pr create \
  --base main \
  --head \"$branch_name\" \
  --title \"Publish Post: $post_title\" \
  --body \"This PR publishes a post titled **$post_title** to the blog. 🎉\""

echo -e "\n===================="
echo "Script completed."
echo "===================="