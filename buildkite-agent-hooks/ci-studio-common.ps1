# Stop script execution when a non-terminating error occurs
$ErrorActionPreference = "Stop"

$env:HAB_NONINTERACTIVE = "true"
$env:HAB_NOCOLORING = "true"

if (Test-Path env:VAULT_UTIL_ACCOUNTS) {
  vault-util configure-accounts
}

if (Test-Path env:VAULT_UTIL_SECRETS) {
  Invoke-Expression (vault-util fetch-secret-env --format ps1)
}