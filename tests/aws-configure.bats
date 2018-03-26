#!/usr/bin/env bats

load test_helpers

@test "aws-configure errors out when no PROFILE is specified" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz

  run aws-configure
  assert_failure

  unset FOO_AWS_ACCESS_KEY_ID
  unset FOO_AWS_SECRET_ACCESS_KEY
}

@test "aws-configure errors out when no AWS variables are specified" {
  run aws-configure foo
  assert_failure
}

@test "aws-configure configures correctly when PROFILE is specified" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz

  stub aws \
    "--profile foo configure list : exit 255" \
    "configure set aws_access_key_id bar --profile foo : echo set aws_access_key" \
    "configure set aws_secret_access_key baz --profile foo : echo set secret_access_key" \
    "configure set region us-east-1 --profile foo : echo set region"

  run aws-configure foo
  assert_success

  unstub aws

  unset FOO_AWS_ACCESS_KEY_ID
  unset FOO_AWS_SECRET_ACCESS_KEY
}

@test "aws-configure is-configured correctly detects when PROFILE is configured" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz

  stub aws \
    "--profile foo configure list : exit 0"

  run aws-configure is-configured foo
  assert_success

  unstub aws

  unset FOO_AWS_ACCESS_KEY_ID
  unset FOO_AWS_SECRET_ACCESS_KEY
}

@test "aws-configure is-configured correctly detects when PROFILE is unconfigured" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz

  stub aws \
    "--profile foo configure list : exit 255"

  run aws-configure is-configured foo
  assert_failure

  unstub aws

  unset FOO_AWS_ACCESS_KEY_ID
  unset FOO_AWS_SECRET_ACCESS_KEY
}
