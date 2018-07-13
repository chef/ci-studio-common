#!/bin/bash
echo "List of containers"
docker ps -a

for i in $(docker ps -aq)
do
  echo "Stopping container $i"
  docker stop "$i" || true

  # pause to allow docker to auto-remove container if it was run with --rm
  sleep 5

  echo "Removing container $i"
  docker rm --force "$i" || true
done
