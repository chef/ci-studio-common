#!/usr/bin/env bats

load test_helpers

@test "install-tool fails when no tool is specified" {
  run install-tool
  assert_failure
}

@test "install-tool habitat" {
  stub ceval \
    "'bash $INSTALL_TOOL_DIR/hab/install.sh ' : echo installed hab"

  run install-tool hab
  assert_success

  unstub ceval
}

@test "install-tool aws" {
  stub ceval \
    "'bash $INSTALL_TOOL_DIR/aws/install.sh ' : echo installed aws"

  run install-tool aws
  assert_success

  unstub ceval
}

@test "install-tool docker-compose" {
  stub ceval \
    "'bash $INSTALL_TOOL_DIR/docker-compose/install.sh ' : echo installed docker-compose"

  run install-tool docker-compose
  assert_success

  unstub ceval
}

@test "install-tool terraform" {
  stub ceval \
    "'bash $INSTALL_TOOL_DIR/terraform/install.sh ' : echo installed terraform"

  run install-tool terraform
  assert_success

  unstub ceval
}

@test "install-tool chefdk" {
  stub ceval \
    "'bash $INSTALL_TOOL_DIR/chefdk/install.sh ' : echo installed chefdk"

  run install-tool chefdk
  assert_success

  unstub ceval
}
