#!/usr/env/bin bats

load test_helpers

@test "configure-github-account errors out if environment variable is unset" {
  run configure-github-account foo

  assert_failure
  assert_output --partial "configure-github-account: Unable to find the required environment variable: FOO_GITHUB_TOKEN"
}

@test "configure-github-account configures the .netrc" {
  export GITHUB_HOME="$HOME"
  export FOO_GITHUB_TOKEN=some_token

  run configure-github-account foo
  assert_success

  run cat $GITHUB_HOME/.netrc
  assert_output --partial "login some_token"

  rm $GITHUB_HOME/.netrc

  unset FOO_GITHUB_TOKEN
  unset GITHUB_HOME
}

@test "configure-github-account errors out if no profile is given" {
  run configure-github-account
  assert_failure
}
