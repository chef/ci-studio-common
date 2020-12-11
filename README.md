# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

<!-- You don't need to modify this TOC. It will automatically update when a PR is merged using Expeditor. -->

<!-- toc -->

- [Announcing the 2.0 release](#announcing-the-20-release)
- [CI Services (Travis, Buildkite, etc)](#ci-services-travis-buildkite-etc)
  * [Installation](#installation)
  * [Commands](#commands)
    + [`did-modify`](#did-modify)
    + [`file-mod`](#file-mod)
    + [`vault-util`](#vault-util)
- [FAQ](#faq)
  * [How do you determine a CI environment vs a non-CI environment?](#how-do-you-determine-a-ci-environment-vs-a-non-ci-environment)

<!-- tocstop -->

## Announcing the 2.0 release

The focus of the 2.0 release is to optimize ci-studio-common for use with multiple Buildkite platform. As such, we have made the breaking following changes:

* Utilities are now being written in Go to ensure consistency and availability across supported platforms.
* Some deprecated utilities have been removed, or are in the process of being removed.
* We no longer ship the CI binaries as part of the Chef Habitat Package


## CI Services (Travis, Buildkite, etc)

### Installation

You can install each component individually or the entirety of ci-studio-common grouped per OS and Arch at https://github.com/chef/ci-studio-common/releases

Currently available for:
OS | ARCH
--- | ---
Darwin | x86-64
Linux | x86-64
Windows | x86-64

### Commands
<!--
  Many of the Helpers are self-documenting. If you see the stdout comment tags, that means that documentation block
  is automatically updated everytime a PR is merged by executing the .expeditor/update_readme.sh script. The implication
  there is that you do not need to manually update those docs.
-->

#### `did-modify`

<!-- stdout "./build/linux/did-modify --help" -->
```
Prints "true" to STDOUT if any files matching GLOBS were modified between HEAD and GITREF. Otherwise, prints "false".

Usage:
  did-modify [flags]

Flags:
      --git-ref string   A valid Git reference (e.g. HEAD, master, origin/master, etc). (default "HEAD~1")
      --globs strings    Comma-separated list of glob patterns to inspect to determine if there are changes. (default [*])
  -h, --help             help for did-modify
```
<!-- stdout -->

#### `file-mod`

<!-- stdout "./build/linux/file-mod --help" -->
```
Command line utility to modify files.

Usage:
  file-mod [command]

Available Commands:
  append-if-missing Append STRING to FILE if not already there.
  find-and-replace  Replace REGEX_STR with STRING in FILE. Supports multiline replace.
  help              Help about any command

Flags:
  -h, --help   help for file-mod

Use "file-mod [command] --help" for more information about a command.
```
<!-- stdout -->

#### `vault-util`
<!-- stdout "./build/linux/vault-util --help" -->
```
Utility to access secrets and account information stored in Hashicorp Vault from CI.

Usage:
  vault-util [command]

Available Commands:
  configure-accounts    Configure the accounts specified in the VAULT_UTIL_ACCOUNTS environment variable.
  fetch-secret-env      Fetch the secrets specified in the VAULT_UTIL_SECRETS environment variable from Vault.
  help                  Help about any command
  print-git-credentials Utility that will print credentials for a GitHub App from Vault in git-credential-helper format.

Flags:
  -h, --help   help for vault-util

Use "vault-util [command] --help" for more information about a command.
```
<!-- stdout -->

## FAQ

### How do you determine a CI environment vs a non-CI environment?

`ci-studio-common` determines whether it is operating in a CI environment by the presence of a `CI` environment variable.
