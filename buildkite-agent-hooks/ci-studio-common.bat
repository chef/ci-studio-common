@ECHO OFF

SET HAB_NONINTERACTIVE="true"
SET HAB_NOCOLORING="true"

IF DEFINED VAULT_UTIL_ACCOUNTS (
  vault-util.exe configure-accounts
)

IF DEFINED VAULT_UTIL_SECRETS (
  FOR /F "USEBACKQ TOKENS=*" %%F IN (`vault-util.exe fetch-secret-env --format batch`) DO SET %%F
)