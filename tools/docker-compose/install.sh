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

echo "WARNING: The 'docker-compose' tool is being deprecated. Please use the 'core/docker-compose' Habitat package intead."

version="${1:-1.14.0}"
system=$(uname -s)
machine=$(uname -m)

if [ ! -d "$HOME/tools/docker-compose-$version" ]; then
  mkdir -p $HOME/tools/docker-compose-$version
  pushd $HOME/tools > /dev/null
    curl -sLo docker-compose "https://github.com/docker/compose/releases/download/$version/docker-compose-$system-$machine"
    chmod +x docker-compose
    mv docker-compose $HOME/tools/docker-compose-$version
    ln -sf $HOME/tools/docker-compose-$version/docker-compose $HOME/tools/bin/docker-compose
  popd > /dev/null
fi

echo ""
echo "docker-compose --version"
$HOME/tools/bin/docker-compose --version
