#!/bin/bash
if [[ "$BUILDKITE_PULL_REQUEST" != "false" ]]; then
  echo "Switching to refspec 'refs/pull/$BUILDKITE_PULL_REQUEST/merge'"
  git fetch origin "+refs/pull/$BUILDKITE_PULL_REQUEST/merge"
  git checkout -qf FETCH_HEAD
fi
