#!/bin/bash
# This script can be used to run the expeditor CLI commands in the environment hook

export HAB_NONINTERACTIVE=true
export HAB_NOCOLORING=true

# Sudo does not exist in msys
run_cmd() {
  if [[ $OSTYPE == "msys" ]]; then
    $@
  else
    sudo -E $@
  fi
}

if [[ -n "${VAULT_UTIL_ACCOUNTS:-}" ]]; then
  vault-util configure-accounts
fi

if [[ -n "${VAULT_UTIL_SECRETS:-}" ]]; then
  . <(vault-util fetch-secret-env)
fi
