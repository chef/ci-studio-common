#!/usr/env/bin bats

load test_helpers

@test "studio-common errors out if executed (not sourced)" {
  run studio-common
  assert_failure
  assert_output --partial "ERROR: studio-common is designed to be sourced, not executed."
}
