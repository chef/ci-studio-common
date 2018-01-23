#!/bin/sh
#
# Distribute the version of terraform, stored in tools/terraform/VERSION, across the components
# in our projects.
#
set -evx

terraform_version=$(cat tools/terraform/VERSION)

sed -i -r "s/^version=.*/version=\"\${1:-$terraform_version}\"/" tools/terraform/install.sh
