#
# Copyright:: Copyright 2019 Chef Software, Inc.
# License:: Apache License, Version 2.0
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

new-module -name CIStudioCommon -scriptblock {
  [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.SecurityProtocolType]'Tls,Tls11,Tls12'

  function Install-Utility {
    param (
      [string]
      $user,
      [string]
      $suffix="rc"
    )

    # Stop script execution when a non-terminating error occurs
    $ErrorActionPreference = "Stop"

    $InstallDir = "$env:SystemDrive\ci-studio-common"
    $SettingsDir = "$env:SystemDrive\ci-studio-settings"

    $RemoteAsset = "https://chef-cd-artifacts.s3-us-west-2.amazonaws.com/ci-studio-common/ci-studio-common-2.0.0-windows-$suffix.tar.gz"

    # Make the directories
    New-Item -ItemType "directory" -Path $SettingsDir -Force | Out-Null

    # Grant user access to settings directory
    # if(-not [string]::IsNullOrEmpty($user)) {
    #   $SettingsAcl = Get-Acl $SettingsDir
    #   $SettingsAr = New-Object System.Security.AccessControl.FileSystemAccessRule($user, "FullControl", "Allow")
    #   $SettingsAcl.SetAccessRule($SettingsAr)
    #   $SettingsAcl | Set-Acl $SettingsDir
    # }

    #  # Grant user access to install directory
    #  if(-not [string]::IsNullOrEmpty($user)) {
    #   $InstallAcl = Get-Acl $InstallDr
    #   $InstallAr = New-Object System.Security.AccessControl.FileSystemAccessRule($user, "FullControl", "Allow")
    #   $InstallAcl.SetAccessRule($InstallAr)
    #   $InstallAcl | Set-Acl $InstallDir
    # }

    $client = New-Object System.Net.WebClient

    Write-Output "Downloading ci-studio-common tarball"
    $client.DownloadFile($RemoteAsset, "$env:SystemDrive\ci-studio-common.tar.gz")

    Write-Output "Extracting ci-studio-common into $InstallDir"
    & cmd.exe '/C 7z x "C:\ci-studio-common.tar.gz" -so | 7z x -aoa -si -ttar -o"C:\"'

    # Perform post-install operations
    Set-Content -Path "$SettingsDir/etag" -Value $client.ResponseHeaders.Get("ETag")

    [Environment]::SetEnvironmentVariable("PATH", "$InstallDir\bin;$env:PATH", "Machine");
    $env:PATH = "$InstallDir\bin;$env:PATH"

    Remove-Item -Path "C:\ci-studio-common.tar.gz"
  }
  set-alias install -value Install-Utility

  export-modulemember -function 'Install-Utility' -alias 'install'
}
