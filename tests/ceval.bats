#!/usr/bin/env bats

load test_helpers

@test "ceval echos COMMAND when DEBUG flag is set" {
  export DEBUG=foo
  run ceval "echo \"I'm a little tea pot\""

  assert_success
  assert_output "echo \"I'm a little tea pot\""
}

@test "ceval evaluates COMMAND when DEBUG flag is unset" {
  run ceval "echo \"I'm a little tea pot\""

  assert_success
  assert_output "I'm a little tea pot"
}