#!/bin/bash

# the blog tag is blog/yyyy-mm-dd

# get the current year
date=$(date +%Y)

# month and day input
echo "Enter the month: "
read month
echo "Enter the day: "
read day

# create the tag
tag="blog/$date-$month-$day"

# ask for commit
echo "Enter the commit: "
read commit

# ask for  the commit message
echo "Enter the commit message: "
read commit_message

# create the tag file
git tag -a $tag $commit -m "$commit_message"

# push the tag
git push origin $tag