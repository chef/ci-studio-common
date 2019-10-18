#!/bin/bash
#
# Copyright:: Copyright 2017 Chef Software, Inc.
# License:: Apache License, Version 2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# The branch we will install the source code from
#
# Example: Install from the 'foo' branch
# => curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash -s -- foo
install_dir="/opt/ci-studio-common"
settings_dir="/var/opt/ci-studio-common"

UNAME=$(uname -sm | awk '{print tolower($0)}')

if [[ ($UNAME == *"mac os x"*) || ($UNAME == *darwin*) ]]; then
  PLATFORM="darwin"
elif [[ ($UNAME == *"freebsd"*) ]]; then
  PLATFORM="freebsd"
else
  PLATFORM="linux"
fi

REMOTE_ASSET="https://chef-cd-artifacts.s3-us-west-2.amazonaws.com/ci-studio-common/ci-studio-common-2.0.0-$PLATFORM-rc.tar.gz"

NEW_ETAG=$(curl -sI $REMOTE_ASSET | grep -Fi Etag | awk '{ print $2 }')

ETAG_PATH="$settings_dir/etag"
OLD_ETAG=""

if [[ -f $ETAG_PATH ]]; then
  OLD_ETAG=$(cat $ETAG_PATH)
fi

function make_directories() {
  rm -rf "$install_dir"
  mkdir -p /opt
  mkdir -p "$settings_dir"
}

function download_and_install_asset() {
  echo "Downloading ci-studio-common for $PLATFORM"
  curl -sL "$REMOTE_ASSET" -o /tmp/ci-studio-common.tar.gz
  tar -xzvf /tmp/ci-studio-common.tar.gz -C /opt
}

function make_symlinks() {
  if [[ -w /usr/bin ]]; then
    ln -sf /opt/ci-studio-common/bin/* /usr/bin
  else
    echo "\

=== WARNING ===

ci-studio-common does not have permission to install binaries into /usr/bin.
Please make sure to add /opt/ci-studio-common/bin to your PATH.

export PATH=$PATH:/opt/ci-studio-common/bin

"
  fi
}

if [[ $NEW_ETAG != $OLD_ETAG ]]; then
  make_directories
  download_and_install_asset
  make_symlinks

  echo -n $NEW_ETAG > $ETAG_PATH
else
  echo "ci-studio-common is up-to-date"
fi
