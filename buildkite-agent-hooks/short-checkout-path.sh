#!/bin/bash

# Check out code into a shorter code path
if [[ $OSTYPE == "msys" ]]; then
  BUILDKITE_BUILD_CHECKOUT_PATH="C:\\bk${BUILDKITE_JOB_ID:0:8}"
else
  BUILDKITE_BUILD_CHECKOUT_PATH="/tmp/bk${BUILDKITE_JOB_ID:0:8}"
fi

export BUILDKITE_BUILD_CHECKOUT_PATH