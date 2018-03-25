#!/bin/bash
#
# Copyright:: Copyright 2018 Chef Software, Inc.
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

version="${1:-1.2.1}"

if [ ! -d "$HOME/tools/packer-$version" ]; then
  mkdir -p $HOME/tools/packer-$version
  pushd $HOME/tools > /dev/null
    curl -sLo packer.zip "https://releases.hashicorp.com/packer/$version/packer_${version}_linux_amd64.zip"
    unzip packer.zip
    rm -f packer.zip || true
    mv packer $HOME/tools/packer-$version
    ln -sf $HOME/tools/packer-$version/packer $HOME/tools/bin/packer
  popd > /dev/null
fi

echo ""
echo "packer --version"
$HOME/tools/bin/packer --version
