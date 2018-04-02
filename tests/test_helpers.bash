#!/bin/bash

 __FILE__="${BASH_SOURCE[0]}"
TESTS_DIR=$( cd "$( dirname "${__FILE__}" )" && pwd )

export CI_STUDIO_COMMON_ROOT_PATH="$BATS_TEST_DIRNAME/.."
export PATH=$CI_STUDIO_COMMON_ROOT_PATH/bin:$HOME/tools/bin:$PATH
export INSTALL_TOOL_DIR="$BATS_TEST_DIRNAME/../tools"

load '/usr/local/lib/bats/load.bash'
