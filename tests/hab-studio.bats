#!/usr/env/bin bats

load test_helpers

@test "hab-studio errors out if environment variable is unset" {
  run hab-studio configure-github-account foo
  assert_failure
  assert_match "hab-studio: Unable to find the required environment variable: FOO_GITHUB_TOKEN"
}

@test "hab-studio errors out if no sub-command is given" {
  run hab-studio
  assert_failure
}

@test "hab-studio configure-github-account configures the .netrc and GITHUB_TOKEN" {
  export FOO_GITHUB_TOKEN=some_token

  run hab-studio configure-github-account foo
  assert_success
  assert_match "The GITHUB_TOKEN environment variable (and ~/.netrc file) will now be available when you start your studio."
  
  source .secrets
  assert_match "true" "$CI"
  assert_match "some_token" "$GITHUB_TOKEN"

  run cat $HOME/.netrc
  assert_match "login $GITHUB_TOKEN"
}

@test "hab-studio configure-github-account errors out if no profile is given" {
  run hab-studio configure-github-account
  assert_failure
}