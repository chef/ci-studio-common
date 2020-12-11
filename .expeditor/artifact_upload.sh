#!/bin/bash

set -eou pipefail

version=$(cat VERSION)

# this will ideally be a directory we export as an env var when processing the config
cd /workdir/workdir/dist
workdir="."

echo "--- Uploading binaries to Artifactory"

# these paths could/would change and we can read those via Expeditor so we arent hardcoding array values below
for file in ${workdir}/**/*
do
  IFS='/' read -r -a path <<< "${file}"
  dir_path=${path[1]}
  binary=${path[2]}

  IFS='_' read -r -a parts <<< "${dir_path}"
  util_name=${parts[0]}
  os=${parts[1]}
  arch=${parts[2]}

  jfrog rt u \
  --apikey="${ARTIFACTORY_TOKEN}" \
  --url=https://artifactory.chef.co/artifactory \
  --props "project=ci-studio-common;version=${version};os=${os};arch=${arch}" \
  ${file} \
  "go-binaries-local/${util_name}/${version}/${os}/${arch}/${binary}"
done