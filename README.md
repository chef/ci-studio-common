# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

## TravisCI

### Installation

```yaml
before_install:
  - curl https://raw.githubusercontent.com/chef/ci-studio-common/master/install.sh | bash
  - export PATH="$PATH:$HOME/ci-studio-common/bin:$HOME/tools/bin"
```

### Helpers

#### `hab-origin`

Helpers that extend functionality of the `hab origin` namespace.

##### `hab-origin download-sig-key $ORIGIN`

Download the private signing key stored in `AWS_S3_BUCKET` (defaults to `ci-studio-common` for internal Chef Software usage). Requires a valid [aws cli configuration](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html).

> For Chef Software projects, reach out to the #jex-team for getting setup to use the chef-ci credentials.

#### `hab-verify`

Helpers that perform common verification steps inside Habitat.

##### `hab-verify in-studio $COMMAND`

Run `$COMMAND` in the contexts of a `hab studio` started in the root of the project repository.

##### `hab-verify lint`

Perform some simple linting against `plan.sh` files.

##### `hab-verify syntax`

Perform some simple syntax checks against `plan.sh` files.

#### `install-tool`

This utility installs tooling in a very specific pattern specific to Travis CI. You can see all the tools we support in the [tools directory](https://github.com/chef/ci-studio-common/tree/master/tools).

```yaml
install: install-tool aws
```

#### `run-if-changed`

If your project is comprised of many components (or lots of tests), you can use `run_if_changed` to only execute test suites if something changed in given `WORKDIR`.

```yaml
env:
  global:
    - TRAVIS_SUGAR_FILTER_TESTS=true
    - TRAVIS_SUGAR_FORCE_IF_TRAVIS_YAML_CHANGED=true

matrix:
  include:
    - env:
        - NAME=component_test
      script: WORKDIR=component_dir run-if-changed make <test>
```

##### Tuning Options

###### `TRAVIS_SUGAR_FILTER_TESTS`

Tells `run_if_changed` to filter only when we're on a Pull Request. Will run the test, even if no files have changed, if set to false.

To keep things running fast, try to either a) cache as much as possible or b) put the setup behind the `run_if_changed` command as part of the Makefile.

If you only want to filter the tests only on a PR you can use the following setting: `TRAVIS_SUGAR_FILTER_TESTS=$TRAVIS_PULL_REQUEST`

###### `TRAVIS_SUGAR_FORCE_IF_TRAVIS_YAML_CHANGED`

If the .travis.yml file was modified, force all of the tests to run just to be safe. We want to make sure that the changes wouldn't otherwise cause the tests to fail.

## Habitat Studio

### Installation

Please put the following at the **top** of your `.studiorc` file.

```bash
hab pkg install chef/ci-studio-common
source "$(hab pkg path chef/ci-studio-common)/bin/studio-common"
```

### Helpers

#### `document`

The `document` function gives developers a quick and easy way to document the various commands exposed in their Studios, and make that documentation available to the users of the Studio. The `document` function takes two arguments: the name of the function, and a HEREDOC containing documentation about how the command should be used. For examples, you can take a look at the various functions exposed in the `lib` directory of this repository.

The documentation provided via this function is collected and exposed via the `describe` command detailed below.

#### `describe`

The `describe` command allows a Studio users in non-CI environments to quickly view the helper functions that have been made available to them by `chef/ci-studio-common` as well as functions available in their own `.studiorc` and `.studio` environments.

Running the `describe` command with no arguments will bring up a list of all the documented functions available to you (along with the first line of the corresponding documentation as a quick intro).

```
$ describe

The following functions are available for your use in this studio:

  build
    Native studio command to build the Habitat Package (see https://www.habitat.sh/docs/reference/habitat-cli/#hab-pkg-build)
  enforce_hab_version
    Ensure that the installed version of the hab toolchain is >= the given version.
  export_docker_image
    Exports a docker image from the latest habitat build (requires that you have a habitat package already built).
  generate_netrc_config
    Create a .netrc file that can be used to authenticate against GitHub.
  install_hab_packages
    Installs the array of Habitat Packages passed in to the function.
  print_getting_started
    Print out useful instructions on how to get started using this studio.
  wait_for_port_to_listen
    Helper function to wait for a port to be listening
  wait_for_service
    Helper function to wait for services to come online
  wait_for_svc_to_load
    Helper function to wait for a service (svc) to be loaded by the Habitat Supervisor.

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
  # start

  Then on your host, you can hit converge service:

  $ curl http://localhost:1234/version
GETTING_STARTED
```

#### `source_studio_helpers`

Studios may have many, many helpers and the `.studiorc` file can get out of hand pretty quickly. The `source_studio_helpers` command is a quick way to quickly source additional files stored in the `.studio` folder of your `/src` directory.

## FAQ

### How do you determine a CI environment vs a non-CI environment?

`ci-studio-common` determines whether it is operating in a CI environment by the presence of a `CI` environment variable.
