#!/usr/env/bin bats

load test_helpers

@test "hab-verify errors out if no sub-command is given" {
  run hab-verify 
  assert_failure
}