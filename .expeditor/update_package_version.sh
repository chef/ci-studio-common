#!/bin/sh
#
# Distribute the version of terraform, stored in tools/packer/VERSION, across the components
# in our projects.
#
set -evx

packer_version=$(cat tools/packer/VERSION)

sed -i -r "s/^version=.*/version=\"\${1:-$packer_version}\"/" tools/packer/install.sh
