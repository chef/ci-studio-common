#!/usr/env/bin bats

load test_helpers

@test "run-if-changed errors out if no command is specified" {
  run run-if-changed
  assert_failure
}