#!/usr/env/bin bats

load test_helpers

@test "did-modify accepts no parameters" {
  stub git \
    "diff --quiet HEAD~1...HEAD -- * : exit 0"

  run did-modify
  assert_success
  assert_output "false"

  unstub git
}

@test "did-modify accepts git_ref" {
  stub git \
    "diff --quiet origin/master...HEAD -- * : exit 0"

  run did-modify --git_ref="origin/master"
  assert_success
  assert_output "false"

  unstub git
}

@test "did-modify accepts globs" {
  stub git \
    "diff --quiet HEAD~1...HEAD -- tests/fixtures/foo* : exit 1"

  run did-modify --globs="tests/fixtures/foo*"
  assert_success
  assert_output "true"

  unstub git
}

@test "did-modify accepts both git_ref and globs" {
  stub git \
    "diff --quiet origin/master...HEAD -- tests/fixtures/bar* tests/fixtures/foo* : exit 1"

  run did-modify --git_ref="origin/master" --globs="tests/fixtures/bar*,tests/fixtures/foo*"
  assert_success
  assert_output "true"

  unstub git
}
