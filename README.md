# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

<!-- You don't need to modify this TOC. It will automatically update when a PR is merged using Expeditor. -->

<!-- toc -->

- [TravisCI](#travisci)
  * [Installation](#installation)
  * [Environment Variables](#environment-variables)
  * [Helpers](#helpers)
    + [`aws-configure`](#aws-configure)
    + [`ceval`](#ceval)
    + [`citadel`](#citadel)
    + [`hab-origin`](#hab-origin)
    + [`hab-studio`](#hab-studio)
    + [`hab-verify`](#hab-verify)
    + [`install-tool`](#install-tool)
    + [`run-if-changed`](#run-if-changed)
- [Habitat Studio](#habitat-studio)
  * [Installation](#installation-1)
  * [What happens when you source studio-common?](#what-happens-when-you-source-studio-common)
  * [`.studiorc` Helpers](#studiorc-helpers)
    + [`document`](#document)
    + [`add_alias`](#add_alias)
    + [`describe`](#describe)
    + [`getting_started`](#getting_started)
    + [`clone_function`](#clone_function)
  * [The `.studio` directory](#the-studio-directory)
  * [The `.secrets` file](#the-secrets-file)
- [Running Integration Tests](#running-integration-tests)
  * [Getting Started](#getting-started)
  * [Overwriting Integration Testing Functions](#overwriting-integration-testing-functions)
- [FAQ](#faq)
  * [How do you determine a CI environment vs a non-CI environment?](#how-do-you-determine-a-ci-environment-vs-a-non-ci-environment)

<!-- tocstop -->

## TravisCI

### Installation

```yaml
before_install:
  - curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash
  - export PATH="$PATH:$HOME/ci-studio-common/bin:$HOME/tools/bin"
```

If you would like to make an installation from a branch that is under
development, you can add `-s -- BRANCH_NAME` at the end of the `curl` command.
```
curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash -s -- BRANCH_NAME
```

### Environment Variables

Various helpers and other operations assume the presence of certain environment variables. Here are the list of environment variables that should be present:

* `<AWS_PROFILE>_AWS_ACCESS_KEY_ID`
* `<AWS_PROFILE>_AWS_SECRET_ACCESS_KEY`
* `<AWS_PROFILE>_AWS_DEFAULT_REGION`
* `<GITHUB_USER>_GITHUB_TOKEN`

where:

* `<AWS_PROFILE>` is the profile you specify in `aws-configure` (e.g. `CHEF_CD`)
* `<GITHUB_USER>` is the GitHub user associated with the token (e.g. `CHEF_CI`)

> For Chef Software, these values are automatically applied to Travis Projects by Engineering Services.

### Helpers
<!--
  Many of the Helpers are self-documenting. If you see the stdout comment tags, that means that documentation block
  is automatically updated everytime a PR is merged by executing the .expeditor/update_readme.sh script. The implication
  there is that you do not need to manually update those docs.
-->

#### `aws-configure`

<!-- stdout "./bin/aws-configure --help" -->
```
Usage: aws-configure [PROFILE]

A non-interactive version of 'aws configure' that allows you to configure AWS CLI profiles.

To add a new profile, you MUST specify following environment variables:

    * <PROFILE>_AWS_ACCESS_KEY_ID
    * <PROFILE>_AWS_SECRET_ACCESS_KEY

You can optionally specify '<PROFILE>_AWS_DEFAULT_REGION' to determine the default region for this profile.

If no PROFILE is specified, aws-configure will defer to the value specified in AWS_PROFILE.

SUBCOMMANDS:
    is-configured PROFILE     Returns '0' if the profile is configured. Otherwise, it returns '1'.
```
<!-- stdout -->

#### `ceval`

<!-- stdout "./bin/ceval --help" -->
```
Usage: ceval COMMAND

Conditionally evaluate the given COMMAND.

    If the DEBUG environment variable is unset, ceval will evaluate the COMMAND using 'eval'.
    If the DEBUG environment variable is set (to any value), ceval will simply echo the given COMMAND.

GUIDANCE:

  1. This command is intended to wrap desctructive or permanent commands that you do not want executed
     when debugging the parent script. 

        ceval "s3 cp myfile s3://my-bucket/my-file"

  2. If you're command requires double quotes, make sure to escape them.

        ceval "echo \"I'm a little tea pot\""
```
<!-- stdout -->

#### `citadel`

<!-- stdout "./bin/citadel --help" -->
```
Usage: citadel FILE

A Bash utility that prints the contents of the given FILE from the CITADEL_BUCKET in S3 to STDOUT.

Requires that you have an AWS profile configured. To configure an AWS profile, you can use 'aws-configure [PROFILE]'.

ENVIRONMENT VARIABLES:
    CITADEL_BUCKET        The name of the S3 bucket where citadel files are kept. (default: $AWS_PROFILE-citadel)
    CITADEL_PROFILE       The name of the AWS CLI profile with access to citadel. (default: $AWS_PROFILE)
```
<!-- stdout -->

#### `hab-origin`

<!-- stdout "./bin/hab-origin --help" -->
```
Usage: hab-origin [SUBCOMMAND]

Helpers that extend functionality of the hab origin namespace.

SUBCOMMANDS:
    download-sig-key ORIGIN     Download the private signing key for ORIGIN stored in the citadel S3 bucket.
```
<!-- stdout -->

#### `hab-studio`

<!-- stdout "./bin/hab-studio --help" -->
```
Usage: hab-studio [SUBCOMMAND]

Utility to configure aspects of your Habitat Studio prior to launch.

SUBCOMMANDS:
    cleanup                               Remove any state created by any of the hab-studio commands.
    configure-github-account ACCOUNT      Configure studios with GitHub credentials for ACCOUNT.
```
<!-- stdout -->

#### `hab-verify`

<!-- stdout "./bin/hab-verify --help" -->
```
Usage: hab-verify [SUBCOMMAND]

Helpers that perform common verification steps inside Habitat.

SUBCOMMANDS:
    lint              Perform some simple linting against plan.sh files.
    syntax            Perform some simple syntax checks against plan.sh files.
    in-studio CMD     Run CMD in the contexts of a hab-studio started in the root of the project repository.

ENVIRONMENT VARIABLES:
    STUDIO_OPTS       Optional options to pass to the hab-studio.
    HAB_ORIGIN        The Habitat origin associated with the hab-studio.
```
<!-- stdout -->

#### `install-tool`

<!-- stdout "./bin/install-tool --help" -->
```
Usage: install-tool TOOL [VERSION]

Install one of the pre-configured TOOLs (listed below) at the specified VERSION (if specified).

AVAILABLE TOOLS:
    aws
    chefdk [CHANNEL]
    docker-compose [VERSION]
    hab
    terraform [VERSION]
```
<!-- stdout -->

#### `run-if-changed`

<!-- stdout "./bin/run-if-changed --help" -->
```
Usage: GLOBS=<GLOBS> WORKDIR=<WORKDIR> run-if-changed CMD

Only execute CMD from WORKDIR if there are changed files that match GLOBS.

ENVIRONMENT VARIABLES:
    TRAVIS_SUGAR_FORCE_IF_TRAVIS_YAML_CHANGED   Force all tests to run if the '.travis.yml' file was modified. 
    TRAVIS_SUGAR_FILTER_TESTS                   Set to 'false' to force all tests to run. To filter tests only on a PR, use 'TRAVIS_SUGAR_FILTER_TESTS='
    WORKDIR                                     The directory where, if there are changes, CMD should be executed.
    GLOBS                                       A list of glob patterns to inspect to determine if there are changes. Defaults to "<WORKDIR>/*"

Example 1 :: Running Travis tests only if they have changed
-----------------------------------------------------------
  ---
  env:
    global:
      - TRAVIS_SUGAR_FILTER_TESTS=
      - TRAVIS_SUGAR_FORCE_IF_TRAVIS_YAML_CHANGED=true

  matrix:
    include:
      - env: NAME=component_test
        script: WORKDIR=component_dir run-if-changed make test
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

### The `.secrets` file

It is currently difficult/impossible to inject secrets into your Hab Studio. To get around this, `ci-studio-common` will automatically source the `/src/.secrets` file (if it exists). The current recommended practice is to export environment variables containing your secrets.

```bash
export GITHUB_TOKEN="<your personal access token>"
```

> Make sure the `.secret` file is added to your `.gitignore`!

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
