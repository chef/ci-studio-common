#!/usr/bin/env bats

load test_helpers

@test "install-tool fails when run inside a studio" {
  export STUDIO_TYPE=foo

  run install-tool aws
  assert_failure
  assert_match "used in Travis CI only"
}

@test "install-tool fails when no tool is specified" {
  run instal-tool
  assert_failure
}

@test "install-tool habitat" {
  run install-tool hab
  assert_success
}

@test "install-tool aws" {
  run install-tool aws
  assert_success
}

@test "install-tool docker-compose" {
  run install-tool docker-compose
  assert_success
}

@test "install-tool terraform" {
  run install-tool terraform
  assert_success
}

@test "install-tool chefdk" {
  run install-tool chefdk
  assert_success
}