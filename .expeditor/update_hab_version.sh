#!/bin/sh
#
# Update the version of Habitat in .hab-version
#
set -evx

echo "$VERSION" > .hab-version
