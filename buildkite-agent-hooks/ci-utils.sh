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

echo "Updating 'ci-utils'"
run_cmd "ci-utils update"

update_status_file="/var/opt/ci-utils/upgrade-in-progress"
if [[ $OSTYPE == "msys" ]]; then
  update_status_file="C:\\ci-utils\\upgrade-in-progress"
fi

while [[ -f $update_status_file ]]; do
  sleep 1
done

echo "Updating 'hab'"
run_cmd "install-habitat"

if [[ -n "${VAULT_UTIL_ACCOUNTS:-}" ]]; then
  vault-util configure-accounts
fi

if [[ -n "${VAULT_UTIL_SECRETS:-}" ]]; then
  eval "vault-util fetch-secret-env"
fi
