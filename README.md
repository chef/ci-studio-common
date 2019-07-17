# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

<!-- You don't need to modify this TOC. It will automatically update when a PR is merged using Expeditor. -->

<!-- toc -->

- [Announcing the 2.0 release](#announcing-the-20-release)
- [CI Services (Travis, Buildkite, etc)](#ci-services-travis-buildkite-etc)
  * [Installation](#installation)
    + [Linux & macOS](#linux--macos)
    + [Windows (Powershell)](#windows-powershell)
    + [Windows (cmd.exe)](#windows-cmdexe)
  * [Commands](#commands)
    + [`ci-studio-common-util`](#ci-studio-common-util)
    + [`did-modify`](#did-modify)
    + [`file-mod`](#file-mod)
    + [`install-buildkite-agent`](#install-buildkite-agent)
    + [`install-habitat`](#install-habitat)
    + [`vault-util`](#vault-util)
- [Habitat Studio](#habitat-studio)
  * [Installation](#installation-1)
  * [What happens when you source studio-common?](#what-happens-when-you-source-studio-common)
  * [Setting a Habitat Version](#setting-a-habitat-version)
  * [`.studiorc` Helpers](#studiorc-helpers)
    + [`document`](#document)
    + [`add_alias`](#add_alias)
    + [`describe`](#describe)
    + [`getting_started`](#getting_started)
    + [`clone_function`](#clone_function)
  * [The `.studio` directory](#the-studio-directory)
- [Running Integration Tests](#running-integration-tests)
  * [Getting Started](#getting-started)
  * [Overwriting Integration Testing Functions](#overwriting-integration-testing-functions)
- [Development](#development)
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

#### Linux & macOS

```bash
curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash
```

If you would like to make an installation from a branch that is under development, you can add `-s -- BRANCH_NAME` at the end of the `curl` command.

```bash
curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash -s -- BRANCH_NAME
```

#### Windows (Powershell)

```powershell
. { iwr -useb https://raw.githubusercontent.com/chef/ci-studio-common/master/install.ps1 } | iex; install
```

If you would like to make an installation from a branch that is under development, you can add it to the command.

```powershell
. { iwr -useb https://raw.githubusercontent.com/chef/ci-studio-common/master/install.ps1 } | iex; install -branch BRANCH_NAME
```

#### Windows (cmd.exe)

```
@"%SystemRoot%\System32\WindowsPowerShell\v1.0\powershell.exe" -NoProfile -InputFormat None -ExecutionPolicy Bypass -Command ". { iwr -useb https://raw.githubusercontent.com/chef/ci-studio-common/master/install.ps1 } | iex; install"
```

### Commands
<!--
  Many of the Helpers are self-documenting. If you see the stdout comment tags, that means that documentation block
  is automatically updated everytime a PR is merged by executing the .expeditor/update_readme.sh script. The implication
  there is that you do not need to manually update those docs.
-->

#### `ci-studio-common-util`

<!-- stdout "./build/linux/ci-studio-common-util --help" -->
```
Utility operations to manage the installation of ci-studio-common

Usage:
  ci-studio-common-util [command]

Available Commands:
  allow       Allow USER to perform certain necessary operations with sudo.
  help        Help about any command
  update      Update the ci-studio-common install

Flags:
  -h, --help   help for ci-studio-common-util

Use "ci-studio-common-util [command] --help" for more information about a command.
```
<!-- stdout -->

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

#### `install-buildkite-agent`

<!-- stdout "./build/linux/install-buildkite-agent --help" -->
```
Manage the Buildkite Agent installation

Usage:
  install-buildkite-agent [command]

Available Commands:
  help        Help about any command
  hook        Install one of the supported HOOKS as an Buildkite Agent Hook.

Flags:
  -h, --help   help for install-buildkite-agent

Use "install-buildkite-agent [command] --help" for more information about a command.
```
<!-- stdout -->

#### `install-habitat`

<!-- stdout "./build/linux/install-habitat --help" -->
```
Install VERSION of Chef Habitat from CHANNEL.

Usage:
  install-habitat [flags]
  install-habitat [command]

Available Commands:
  help        Help about any command
  remove      Completely uninstall Chef Habitat from the system.

Flags:
  -c, --channel string   The channel from which you wish to install Habitat. (default "stable")
  -h, --help             help for install-habitat
  -t, --target string    The kernel target for this installation. (default "x86_64-linux")
  -v, --version string   Which version of Habitat you wish to install. (default "0.82.0")

Use "install-habitat [command] --help" for more information about a command.
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
  print-git-credentials Utility that will print credentials for a user from Vault in git-credential-helper format.

Flags:
      --config string   configuration file (default is /var/opt/ci-studio-common/vault-util.toml)
  -h, --help            help for vault-util

Use "vault-util [command] --help" for more information about a command.
```
<!-- stdout -->

## Habitat Studio

### Installation

Please put the following at the **top** of your `.studiorc` file.

```bash
# shellcheck disable=1090
if [ -d "${CI_STUDIO_COMMON:-}" ]; then
  echo "CI_STUDIO_COMMON override in effect; using $CI_STUDIO_COMMON, not chef/ci-studio-common habitat package"
  source "$CI_STUDIO_COMMON/bin/studio-common"
else
  hab pkg install chef/ci-studio-common
  source "$(hab pkg path chef/ci-studio-common)/bin/studio-common"
fi
```

Technically only the `else` body is necessary to _use_ `ci-studio-common`, but the `if` half allows the use of a locally modified copy of `ci-studio-common`, so changes can be delevoped and tested without building, promoting and installing a new `chef/ci-studio-common` habitat package. See [Development](#development) for more.

### What happens when you source studio-common?

1. Source all the Helper functions that come with `ci-studio-common`.
2. Source any helpers defined in files in your `.studio` directory.
2. Source your `/src/.secrets` file (if it exists).

### Setting a Habitat Version

`ci-studio-common` allows you to let your users know when they are running a version of Habitat that is either newer or older than the recommended version. By default, we recommend the version stored in the `.hab-version` file. However, you can set your own recommended version by setting the `RECOMMENDED_HAB_VERSION` variable before sourcing the studio code.

```bash
RECOMMENDED_HAB_VERSION="0.69.0"
hab pkg install chef/ci-studio-common
source "$(hab pkg path chef/ci-studio-common)/bin/studio-common"
```

### `.studiorc` Helpers

#### `document`

The `document` function gives developers a quick and easy way to document the various commands exposed in their Studios, and make that documentation available to the users of the Studio. The `document` function takes two arguments: the name of the function, and a HEREDOC containing documentation about how the command should be used. For examples, you can take a look at the various functions exposed in the `lib` directory of this repository.

The documentation provided via this function is collected and exposed via the `describe` command detailed below.

#### `add_alias`

The `add_alias` function performs double duty. It will a) add the alias you specify, and b) automatically document that the alias exists.

By specifying `add_alias` in conjunction with a function and a `document` block, the `describe` block will not only show you the functions but the aliases as well.

For example, let's look at the imaginary function `foobear`. We want to alias `foobear` to `fb`. The resulting code block should look like this.

```bash
document "foobear" <<DOC
  Foobear is not a bear.
DOC
add_alias "foobear" "fb"
function foobear() {
  echo "Foo"
}
```

Now, when we run `describe`, we'll see `foobear` documented like so:

```
foobear [alias: fb]
  Foobear is not a bear.

...

ALIASES:
  fb   Alias for: foobear
```

#### `describe`

The `describe` command allows a Studio users in non-CI environments to quickly view the helper functions that have been made available to them by `chef/ci-studio-common` as well as functions available in their own `.studiorc` and `.studio` environments.

Running the `describe` command with no arguments will bring up a list of all the documented functions available to you (along with the first line of the corresponding documentation as a quick intro).

```
$ describe

The following functions are available for your use in this studio:

  build
    Wrapper around the 'hab pkg build' command.
  clone_function
    Create a copy of the ci-studio-common function in your .studio/ directory.
  configure_host
    Configure a host (in /etc/hosts) inside the studio.
  enforce_integration_testing
    Print an error message if any of the required Integration Testing functions are missing.
  export_docker_image
    Export a Docker Image of this Habitat package.
  generate_netrc_config
    Create a .netrc file that can be used to authenticate against GitHub.
  inspec_exec
    (MISSING) You must provide a custom inspec_exec() function.
  install
    Install the specified Habitat packages (defaults to chef/ci-studio-common).
  install_if_missing
    Install the package and binlink the binary only if it is missing.
  integration_tests
    Provision, Verify, and Destroy your full Integration Testing environment.
  ipaddress
    Returns the IPv4 address of the studio.
  start_dependencies
    (MISSING) You must provide a custom start_dependencies() function.
  start_service
    (MISSING) You must provide a custom start_service() function.
  stop_dependencies
    (MISSING) You must provide a custom stop_dependencies() function.
  stop_service
    (MISSING) You must provide a custom stop_service() function.
  sup-log [alias: sl]
    Tail the Habtiat Supervisor's output.
  sup-run [alias: sr]
    Launch the Habtiat Supervisor in the background.
  sup-term [alias: st]
    Kill the Habitat Supervisor running in the background.
  wait_for_port_to_listen
    Wait for a port to be listening.
  wait_for_success
    Wait for the given command to succeed.
  wait_for_svc_to_load
    Helper function to wait for a Habitat service (hab svc) to be loaded by the Habitat Supervisor.

ALIASES:
  sl     Alias for: sup-log
  sr     Alias for: sup-run
  st     Alias for: sup-term

To learn more about a particular function, run 'describe <function>'.
```

Passing the name of a function into the `describe` command will show you the full documentation for that command.

```
$ describe generate_netrc_config

Create a .netrc file that can be used to authenticate against GitHub.

Some projects require access to private GitHub repositories. The recommended
pattern is for projects to use the git+https protocol in conjuction with a
.netrc file.

To learn more about .netrc files, you can check out the following documentation:

  https://www.gnu.org/software/inetutils/manual/html_node/The-_002enetrc-file.html
```

#### `getting_started`

The `getting_started` helper is intended to be used in your `.studiorc`. In non-CI environments, the contents of the HEREDOC you pass into this function will get printed when the user first launches the studio.

```bash
getting_started <<GETTING_STARTED
  Welcome to the Habitat-based development environment for the Best Service.

  === Getting Started ===

  From the studio run:

  # build
  # start_service

  Then on your host, you can hit converge service:

  $ curl http://localhost:1234/version
GETTING_STARTED
```

#### `clone_function`

Despite our best efforts, sometimes a helper function that we provide to you may not suit your needs exactly. In that scenario, you can use the `clone_function` helper to clone an existing function into your local `.studio/` directory under the file `override-common`. There, you can modify the function to meet your needs. From that point on, anytime you open up a new Habitat Studio _your_ version of the helper will load, rather than the version that comes with `ci-studio-common`. If, at any time you wish to go back to the original function, you can simply remove the function declaration from your source code.

```bash
[01][default:/src:0]# ls -alh .studio
ls: cannot access '.studio': No such file or directory
[02][default:/src:2]# clone_function install
[47][default:/src:0]# ls -alh .studio
total 4.0K
drwxr-xr-x  3 root root  102 Oct  3 18:09 .
drwxr-xr-x 21 root root  714 Oct  3 18:09 ..
-rw-r--r--  1 root root 1.5K Oct  3 18:09 override-common
[03][default:/src:0]# cat .studio/override-common
#!/bin/bash
# In this file, you can override functions included with ci-studio-common to meet your needs.

document "install" <<DOC
  Install the specified Habitat packages (defaults to chef/ci-studio-common).

  @(arg:*) The array of packages you wish to install (none will build chef/ci-studio-common)

  Example 1 :: Install the package described in plan.sh
  -----------------------------------------------------
  install

  Example 2 :: Install (and binlink) the listed packages from the 'stable' channel.
  ---------------------------------------------------------------------------------
  OPTS="--binlink" install core/curl core/git

  Example 3 :: Install the listed packages from the 'unstable' channel.
  ---------------------------------------------------------------------
  dev_dependencies=(core/curl core/git)
  OPTS="--channel unstable" install """
DOC
function install ()
{
    install_cmd="install";
    if [[ "x$OPTS" != "x" ]]; then
        install_cmd="$install_cmd $OPTS";
    fi;
    if [[ "x$1" == "x" ]]; then
        pushd /src > /dev/null;
        if [[ ! -f results/last_build.env ]]; then
            build .;
            source results/last_build.env;
            eval "hab pkg $install_cmd $pkg_ident >/dev/null";
        fi;
        popd;
    else
        for pkg in "$@";
        do
            echo "Installing $pkg";
            eval "hab pkg $install_cmd $pkg >/dev/null";
        done;
    fi
}
```


### The `.studio` directory

If you have a lot of helpers, putting them all in your `.studiorc` file can quickly result in a large, difficult to comprehend file. `ci-studio-common` allows you to split up those helpers into logical files and store them in a `.studio` directory (much like the `dot-studio` folder of this repository).

When you `source "$(hab pkg path chef/ci-studio-common)/bin/studio-common"`, all the files in your `.studio` directory will automatically be sourced. Any `document` tags you specify will also automatically be made available under `describe`.

## Running Integration Tests

One of the functions of the `ci-studio-common` library is to allow you to quickly and easily spin up an Integration environment and run InSpec tests against it. This can be done by using a simple command like this:

```
hab studio run integration_tests
```

It can also be done from inside the studio by executing the `integration_tests` function.

```
[1][default:/src:0]# integration_tests
```

### Getting Started

To get started, the first thing you'll want to do is add the `enforce_integration_testing` command towards the top of your `.studiorc` file, after you source `ci-studio-common`. There are a number of functions that you will need to overwrite in your own `.studiorc` file or `.studio` directory, and `enforce_integration_testing` ensures that those functions exist.

```bash
hab pkg install chef/ci-studio-common
source "$(hab pkg path chef/ci-studio-common)/bin/studio-common"

enforce_integration_testing
```

The list of necessary functions are:

* `start_service`
* `start_dependencies`
* `stop_dependencies`
* `stop_service`

### Overwriting Integration Testing Functions

We recommend that you create an `integration_testing` file in your `.studio` directory, and put all your function definitions in that file. The recommended function definition would look something like this:

```bash
document "start_service" <<DOC
  A quick description of what this process does. Be specific! This is your user-facing documentation.
DOC
function start_service() {
  # All the things you need to do to start your service
}
```

## Development

If the `CI_STUDIO_COMMON` override pattern is implemented in the `.studiorc` which sources `ci-studio-common` ([example in a2](https://github.com/chef/a2/blob/master/.studiorc#L11-L14)), local changes to can be tested by placing the `ci-studio-common` checkout inside the repo which is sourcing it. This is necessary so that it is accessible from the studio `chroot`ed environment. To set the `CI_STUDIO_COMMON` environment variable inside the studio when `.studiorc` is sourced, we can make use of the `HAB_STUDIO_SECRET` facility ([documented here](https://www.habitat.sh/docs/reference/#environment-variables)). For example, if our `ci-studio-common` is at the same level of our `.studiorc`, we can enter the studio like this:

```
$ env HAB_STUDIO_SECRET_CI_STUDIO_COMMON=/src/ci-studio-common hab studio enter
```

(This is a bit of a hack as the value of CI_STUDIO_COMMON doesn't need to be secret, but it's currently the easiest way to get environment variables set inside the studio before `.studiorc` is sourced.)

## FAQ

### How do you determine a CI environment vs a non-CI environment?

`ci-studio-common` determines whether it is operating in a CI environment by the presence of a `CI` environment variable.
