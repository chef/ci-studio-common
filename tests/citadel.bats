#!/usr/env/bin bats

load test_helpers

@test "citadel errors out if aws cli is missing" {
  run citadel foo-file
  assert_failure
  assert_match "citadel: Unable to find 'aws' executable"
}

@test "citadel errors out if the CITADEL_PROFILE is not configured" {
  export CITADEL_PROFILE=foo
  install-tool aws

  run citadel foo-file
  assert_failure 
  assert_match "citadel: The 'foo' aws profile is not configured"
}

@test "citadel errors out when no file is specified" {
  install-tool aws

  run citaldel
  assert_failure
}

# When able, I want to be able to provide a mock wrapper around the aws s3.