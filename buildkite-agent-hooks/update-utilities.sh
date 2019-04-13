#!/bin/bash
# This script has been deprecated in favor of the expeditor environment hook.

export HAB_NONINTERACTIVE=true
export HAB_NOCOLORING=true

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

  echo "Updating 'expeditor-cli'"
  (
    echo "Installing Expeditor CLI with exclusive lock (timeout 120s)..."
    flock --exclusive --wait 120 201
    sudo -E hab pkg install --channel "${EXPEDITOR_CHANNEL:-stable}" chef-es/expeditor-cli
  ) 201>/tmp/hab-pkg-install-expeditor-cli.lock
fi
