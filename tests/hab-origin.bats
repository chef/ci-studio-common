#!/usr/env/bin bats

load test_helpers

@test "hab-origin errors out when no sub-command is specified" {
  run hab-origin
  assert_failure
}
