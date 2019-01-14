#!/bin/bash
#
# Update the version of Habitat in .hab-version
#
set -evx

if [[ -z $VERSION ]]; then
  echo "ERROR: Version undefined, cannot continue, exiting"
  exit 1
fi

branch="expeditor/upgrade-hab"
git checkout -b ${branch}

echo "$VERSION" > .hab-version

git add .
git commit --message "Update to habitat $VERSION" --message "This pull request was triggered automatically via Expeditor when Habitat $VERSION was promoted to stable." --message "This change falls under the obvious fix policy so no Developer Certificate of Origin (DCO) sign-off is required."

open_pull_request

git checkout -
git branch -D ${branch}
