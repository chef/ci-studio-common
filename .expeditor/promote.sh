#!/bin/bash

set -eou pipefail

if [ "${EXPEDITOR_TARGET_CHANNEL}" = "unstable" ];
then
  echo "This file does not support actions for artifacts promoted to unstable"
  exit 1
fi

version="${EXPEDITOR_PROMOTABLE}"
bucket=chef-automate-artifacts
asset="ci-utils-${paltform}.tar.gz"

function promote_s3_asset() {
  platform=${1:-linux}
  asset="ci-utils-${platform}.tar.gz"

  aws --profile chef-cd s3 cp "s3://${bucket}/files/ci-utils/${version}/${asset}" "s3://${bucket}/${EXPEDITOR_TARGET_CHANNEL}/latest/ci-utils/${asset}" --acl public-read
}

promote_s3_asset linux
promote_s3_asset darwin
promote_s3_asset windows
