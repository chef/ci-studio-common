#!/usr/bin/env bats

load test_helpers

@test "aws-configure errors out when aws-cli is missing" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz

  run aws-configure foo
  assert_failure
  assert_match "aws-configure: Unable to find 'aws' executable"
}

@test "aws-configure errors out when no PROFILE is specified" {
  unset AWS_PROFILE
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz
  install-tool aws

  run aws-configure
  assert_failure
}

@test "aws-configure errors out when no AWS variables are specified" {
  install-tool aws

  run aws-configure foo
  assert_failure
}

@test "aws-configure configures correctly when AWS_PROFILE is specified" {  
  export AWS_PROFILE=foo
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz
  install-tool aws

  run aws-configure
  assert_success
}

@test "aws-configure configures correctly when PROFILE is specified" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz
  install-tool aws

  run aws-configure foo
  assert_success
}

@test "aws-configure is-configure correctly detects when PROFILE is configured" {
  export FOO_AWS_ACCESS_KEY_ID=bar
  export FOO_AWS_SECRET_ACCESS_KEY=baz
  install-tool aws 

  run aws-configure is-configured foo
  assert_failure

  aws-configure foo
  run aws-configure is-configured foo
  assert_success

}