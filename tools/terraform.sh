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

version="${1-'0.8.1'}"

if [ ! -d "$HOME/tools/terraform-$version" ]; then
  mkdir -p "$HOME/tools/terraform-$version"
  pushd $HOME/tools
    curl -sLo terraform.zip "https://releases.hashicorp.com/terraform/$version/terraform_${version}_linux_amd64.zip"
    unzip terraform.zip
    rm -f terraform.zip || true
    ln -sf "$HOME/tools/bin/terraform" "$HOME/tools/terraform-$version/terraform"
  popd
fi

echo "terraform --version"
terraform --version
