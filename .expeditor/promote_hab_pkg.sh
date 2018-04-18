#!/bin/bash

set -eou pipefail

# We pipe this to jq here so we can get only the ident we care about, as there may be invalid characters
# in the full results that cause jq to exit with this:
#
# parse error: Invalid string: control characters from U+0000 through U+001F must be escaped at line 15, column 12
#
results=$(curl --silent https://willem.habitat.sh/v1/depot/channels/chef/dev/pkgs/ci-studio-common/latest | jq '.ident')

pkg_origin=$(echo "$results" | jq -r .origin)
pkg_name=$(echo "$results" | jq -r .name)
pkg_version=$(echo "$results" | jq -r .version)
pkg_release=$(echo "$results" | jq -r .release)

hab pkg promote "${pkg_origin}/${pkg_name}/${pkg_version}/${pkg_release}" stable
