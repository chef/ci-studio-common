#!/bin/sh
#
# Distribute the version of Habitat, stored in .hab-version, across the components
# in our projects.
#
set -evx

hab_version=$(cat tools/hab/VERSION)

echo "$hab_version" > .hab-version
sed -i -r "s/bash -s -- -v .+$/bash -s -- -v $hab_version/" tools/hab/install.sh
sed -i -r "s/export SUPPORTED_HAB_VERSION=\".+\"/export SUPPORTED_HAB_VERSION=\"$hab_version\"/" bin/studio-common
