#!/usr/env/bin bats

load test_helpers

@test "hab-origin errors out if aws cli is missing" {
  run hab-origin download-sig-key foo
  assert_failure
  assert_match "hab-origin: Unable to find 'aws' executable"
}

@test "hab-origin errors out when no sub-command is specified" {
  run hab-origin
  assert_failure
}