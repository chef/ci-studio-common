#!/usr/env/bin bats

load test_helpers

@test "citadel errors out if the CITADEL_PROFILE is not configured" {
  export CITADEL_PROFILE=foo

  stub aws \
    "--profile foo configure list : exit 255"

  run citadel foo-file
  assert_failure
  assert_output --partial "citadel: The 'foo' aws profile is not configured"

  unstub aws

  unset CITADEL_PROFILE
}

@test "citadel prints out the contents of the file" {
  export CITADEL_PROFILE=foo

  stub aws \
    "--profile foo configure list : exit 0" \
    "s3 cp --profile foo s3://foo-citadel/foo-file - : echo foo-contents"

  run citadel foo-file
  assert_success
  assert_output "foo-contents"

  unstub aws

  unset CITADEL_PROFILE
}

@test "citadel errors out when no file is specified" {
  run citaldel
  assert_failure
}

