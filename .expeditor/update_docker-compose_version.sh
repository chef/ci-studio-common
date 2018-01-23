#!/bin/sh
#
# Distribute the version of docker-compose, stored in tools/docker-compose/VERSION, across the components
# in our projects.
#
set -evx

docker_compose_version=$(cat tools/docker-compose/VERSION)

sed -i -r "s/^version=.*/version=\"\${1:-$docker_compose_version}\"/" tools/docker-compose/install.sh
