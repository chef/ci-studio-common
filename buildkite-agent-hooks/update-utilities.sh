#!/bin/bash

habitat_supported_platform() {
  local habitat_supported_archs_regex='.*x86_64|amd64.*'
  local platform="linux"
  local uname="$(uname -sm | awk '{print tolower($0)}')"

  case "$uname" in
    *"mac os x"* | *darwin*) platform="darwin" ;;
    *"freebsd"*) platform="freebsd" ;;
    *) platform="linux" ;;
  esac

  [[ ($platform =~ darwin|linux) && ($uname =~ $habitat_supported_archs_regex) ]]
}

echo "Updating 'ci-studio-common'"
sudo ci-studio-common-util update

if habitat_supported_platform; then
  echo "Updating 'hab'"
  sudo install-habitat

  # Temporarily removing the expeditor-cli
  # echo "Updating 'expeditor' CLI"
  # hab pkg install chef-es/expeditor-ruby
fi
