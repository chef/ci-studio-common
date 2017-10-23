#!/usr/bin/env bats

@test "habitat" {
  run hab --version
  [ "$status" -eq 0 ]
}

@test "aws" {
  run aws --version
  [ "$status" -eq 0 ]
}

@test "docker-compose" {
  run docker-compose --version
  [ "$status" -eq 0 ]
}

@test "terraform" {
  run terraform --version
  [ "$status" -eq 0 ]
}