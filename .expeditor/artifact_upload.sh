#!/bin/bash

set -eou pipefail

version=$(cat VERSION)
workdir="workdir/workdir/dist"

echo "--- Uploading binaries to Artifactory"
for file in ${workdir}/ci-utils_*/*
do
  IFS='/' read -r -a path <<< "${file}"
  IFS='_' read -r -a parts <<< "${path[1]}"

  jfrog rt u \
  --apikey="${ARTIFACTORY_TOKEN}" \
  --url=https://artifactory.chef.co/artifactory \
  --props "project=ci-studio-common;version=${version};os=${parts[1]};arch=${parts[2]}" \
  /workdir/${path} \
  "go-binaries-local/${parts[0]}/${version}/${parts[1]}/${parts[2]}/${path[2]}"
done