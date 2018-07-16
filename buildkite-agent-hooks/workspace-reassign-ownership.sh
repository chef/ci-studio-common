#!/bin/bash
echo "Re-assigning ownership of all files in the working directory to buildkite-agent so it can delete them"
sudo chown -R buildkite-agent "${HOME}/builds/${BUILDKITE_AGENT_NAME}/${BUILDKITE_ORGANIZATION_SLUG}" || true
