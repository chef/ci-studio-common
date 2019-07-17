@ECHO OFF

SET HAB_NONINTERACTIVE="true"
SET HAB_NOCOLORING="true"

ci-utils.exe update

:CheckForFile
IF NOT EXIST "C:\ci-utils\upgrade-in-progress" GOTO UpgradeComplete

REM Wait a bit and then check for the file again
ping 127.0.0.1 -n1 -w 10000 >NUL
GOTO :CheckForFile

:UpgradeComplete

install-habitat.exe

IF DEFINED VAULT_UTIL_ACCOUNTS (
  vault-util.exe configure-accounts
)

IF DEFINED VAULT_UTIL_SECRETS (
  FOR /F "USEBACKQ TOKENS=*" %%F IN (`vault-util.exe fetch-secret-env --format batch`) DO SET %%F
)
