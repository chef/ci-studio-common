#!/bin/sh
#
# Distribute the version of Habitat, stored in .hab-version, across the components
# in our projects.
#
set -evx

hab_version=$(cat .hab-version)

sed -i -r "s/bash -- -v .+$/bash -- -v $hab_version/" tools/hab.sh
sed -i -r "s/export SUPPORTED_HAB_VERSION=\".+\"/export SUPPORTED_HAB_VERSION=\"$hab_version\"/" bin/studio-common
