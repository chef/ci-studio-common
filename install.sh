#!/bin/bash
#
# Copyright:: Copyright 2017 Chef Software, Inc.
# License:: Apache License, Version 2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
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
branch=${1-master}

if [ -d $HOME/ci-studio-common ]
then
  rm -rf $HOME/ci-studio-common
fi

pushd /tmp
  curl -sLo ci-studio-common.zip https://github.com/chef/ci-studio-common/archive/${branch}.zip
  unzip ci-studio-common.zip
  rm -f ci-studio-common.zip || true
  ci_dir=$(echo "ci-studio-common-${branch}" | tr '/' '-')
  mv $ci_dir "$HOME/ci-studio-common"
popd
