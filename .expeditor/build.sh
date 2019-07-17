#!/bin/bash

set -eou pipefail

function usage() {
  cat <<DOC
Usage: ${0##*/}

Build and release a version of ci-studio-common.

DOC
}

version=$(cat VERSION)
upload=s3 # s3, github
suffix="rc"

while getopts ":u:v:s:t:h" opt; do
  case "${opt}" in
    u)
      upload="$OPTARG"
      ;;
    v)
      version="$OPTARG"
      ;;
    s)
      suffix="$OPTARG"
      ;;
    h)
      usage
      exit 0
      ;;
    \?)
      usage
      exit 1
      ;;
    *)
      ;;
  esac
done
shift $((OPTIND -1))

function clean_all() {
  make clean-all
}

function build_tarball() {
  platform=${1:-linux}
  staging_dir="build/ci-studio-common"

  rm -rf "$staging_dir"
  mkdir -p "$staging_dir/bin"

  make "build-$platform"

  cp -r "build/$platform/"* "$staging_dir/bin"
  cp -r buildkite-agent-hooks $staging_dir

  cd build
  tar -czvf "ci-studio-common-$version-$platform-$suffix.tar.gz" ci-studio-common

  cd -
  make "clean-$platform"
  rm -rf "$staging_dir"
}

function upload_s3_asset() {
  platform=${1:-linux}
  asset="ci-studio-common-$version-$platform-$suffix.tar.gz"

  aws --profile chef-cd s3 cp build/$asset s3://chef-cd-artifacts/ci-studio-common/$asset --acl public-read
}

clean_all

build_tarball linux
build_tarball darwin
build_tarball windows

if [[ $upload == "github" ]]; then
  create_github_release
  upload_release_asset linux
elif [[ $upload == "s3" ]]; then
  upload_s3_asset linux
  upload_s3_asset darwin
  upload_s3_asset windows
fi

clean_all

