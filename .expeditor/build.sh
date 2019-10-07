#!/bin/bash

set -eou pipefail

platform="${1:-linux}"
version=$(cat VERSION)
bucket=chef-automate-artifacts
staging_dir="build/ci-utils"
asset="ci-utils-${paltform}.tar.gz"

rm -rf "$staging_dir"
mkdir -p "$staging_dir/bin"

make "build-$platform"

cp -r "build/$platform/"* "$staging_dir/bin"
cp -r buildkite-agent-hooks $staging_dir

cd build
tar -czvf "$asset" ci-utils

aws --profile chef-cd s3 cp "build/${asset}" "s3://${bucket}/files/ci-utils/${version}/${asset}" --acl public-read
aws --profile chef-cd s3 cp "build/${asset}" "s3://${bucket}/unstable/latest/ci-utils/${asset}" --acl public-read

