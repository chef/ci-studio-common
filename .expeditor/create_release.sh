#!/bin/bash
set -eou pipefail

version=$(cat VERSION)
art_token=$(vault kv get -field token account/static/artifactory/buildkite)
files=""

function download_artifacts {
  os=$1
  arch=$2
  
  echo "--- Artifactory downloading ci-studio-common binaries for ${os} ${arch}"
  jfrog rt dl \
  --apikey="${art_token}" \
  --url="https://artifactory.chef.co/artifactory" \
  --flat \
  --detailed-summary \
  --props "project=ci-studio-common;version=${version};os=${os};arch=${arch}" \
  "go-binaries-local/*"
  "go-binaries/${os}/${arch}/"

  zip -r "ci_stuido_common_${os}_${arch}.zip" go-binaries/${os}/${arch}/
  files="${files} ci_stuido_common_${os}_${arch}.zip"

  for file in go-binaries/${os}/${arch}/*
  do
    if [[ ${os} == "windows" ]]; then
      IFS='.' read -r -a parts <<< "${file}"
      mv ${file} "${parts[0]}_${os}_${arch}.${parts[1]}"
      files="${files} ${parts[0]}_${os}_${arch}.${parts[1]}"
    else
      mv ${file} "${file}_${os}_${arch}"
      files="${files} ${file}_${os}_${arch}"
    fi
  done
}

download_artifacts linux amd64
download_artifacts darwin amd64
download_artifacts windows amd64

notes=$(sed -n -E '/<!-- latest_release (.+) -->|<!-- latest_release -->/,/<!-- latest_release -->/p' CHANGELOG.md)

echo "--- GitHub publish release ${version}"
gh release create ${version} ${files} README.md -n "${notes}" -t "${version}"