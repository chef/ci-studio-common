#!/bin/bash
#
# Update the version of Habitat in .hab-version
#
set -evx

if [[ -z $EXPEDITOR_VERSION ]]; then
  echo "ERROR: Version undefined, cannot continue, exiting"
  exit 1
fi

branch="expeditor/upgrade-hab-${EXPEDITOR_VERSION}"
git checkout -b "$branch"

echo "$EXPEDITOR_VERSION" > .hab-version

git add .
git commit --message "Update to habitat $EXPEDITOR_VERSION" --message "This pull request was triggered automatically via Expeditor when Habitat $EXPEDITOR_VERSION was promoted to stable." --message "This change falls under the obvious fix policy so no Developer Certificate of Origin (DCO) sign-off is required."

open_pull_request

git checkout -
git branch -D "$branch"
