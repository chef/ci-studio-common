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

echo "WARNING: The 'aws' tool is being deprecated. Please use the 'core/aws-cli' Habitat package intead."

pip install --user awscli
ln -sf $HOME/.local/bin/aws $HOME/tools/bin/aws

# There is a weird bug in aws where if you specify an AWS_PROFILE that is not
# configured, any 'aws' command (even 'aws --version') will fail. This is an
# issue in CI systems where AWS_PROFILE has already been specified but hasn't
# had time to be configured by the time we install awscli. So rather than
# running 'aws --version', we grab the version from pip.
echo ""
echo "aws --version"
pip show awscli | grep Version
