#!/bin/bash
set -eou pipefail

version=$(cat VERSION)
art_token=$(vault kv get -field token account/static/artifactory/buildkite)
files=()

function download_artifacts {
  os=$1
  arch=$2
  
  echo "--- Artifactory download ci-studio-common binaries for ${os} ${arch}"
  jfrog rt dl \
  --apikey="${art_token}" \
  --url="https://artifactory.chef.co/artifactory" \
  --flat \
  --detailed-summary \
  --props "project=ci-studio-common;version=${version};os=${os};arch=${arch}" \
  "go-binaries-local/*" "go-binaries/${os}/${arch}/"

  for file in go-binaries/${os}/${arch}/*
  do
    if [[ ${os} == "windows" ]]; then
      IFS='.' read -r -a parts <<< "${file}"
      util_name=${parts[0]}

      zip -j "${util_name}_${os}_${arch}.zip" "${file}"
      files+=("${util_name}_${os}_${arch}.zip")
    else
      tar -czf "${file}_${os}_${arch}.tar.gz" -C "go-binaries/${os}/${arch}" .
      files+=("${file}_${os}_${arch}.tar.gz")
    fi
  done
}

download_artifacts linux amd64
download_artifacts darwin amd64
download_artifacts windows amd64

notes=$(sed -n -E '/<!-- latest_release (.+) -->|<!-- latest_release -->/,/<!-- latest_release -->/p' CHANGELOG.md)

echo "--- GitHub publish release ${version}"
gh release create "${version}" "${files[@]}" README.md -n "${notes}" -t "${version}"