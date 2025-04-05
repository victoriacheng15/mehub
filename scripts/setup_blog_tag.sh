#!/bin/bash

# the blog tag is blog/yyyy-mm-dd
echo "Enter the year (default is current year): "
read -r year
echo "Enter the month (default is current month): "
read -r mm
echo "Enter the day (default is current day): "
read -r dd

year=$(date +%Y)
month=${mm:-$(date +%m)}
day=${dd:-$(date +%d)}

date=$year-$month-$day

# create the tag
tag="blog/$date"

# ask for a commit tag
echo "Enter the commit tag: "
read -r commit

# ask for the commit message
echo "Enter the commit message: "
read -r commit_message

echo "Creaating the tag..."
# create the tag file
git tag -a $tag $commit -m "$commit_message"

sleep 2

echo "Pushing the tag..."
# push the tag
git push origin $tag