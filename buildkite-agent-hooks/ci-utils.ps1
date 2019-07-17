# Stop script execution when a non-terminating error occurs
$ErrorActionPreference = "Stop"

$env:HAB_NONINTERACTIVE = "true"
$env:HAB_NOCOLORING = "true"

Write-Output "Updating 'ci-utils'"
ci-utils.exe update

$update_status_file = "C:\ci-utils\upgrade-in-progress"
while(Test-Path $update_status_file) {
  Start-Sleep 1
}

Write-Output "Updating 'hab'"
install-habitat

if (Test-Path env:VAULT_UTIL_ACCOUNTS) {
  vault-util configure-accounts
}

if (Test-Path env:VAULT_UTIL_SECRETS) {
  Invoke-Expression (vault-util fetch-secret-env --format ps1)
}
