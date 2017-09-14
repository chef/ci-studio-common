# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

## TravisCI Scripts

These are scripts that can be used inside your `.travis.yml` files.

* `install_terraform`
* `run_if_changed`

### Installation

```yaml
before_install:
  - wget https://github.com/chef/ci-studio-common/archive/master.tar.gz -O /tmp/ci-studio-common.tar.gz
  - tar -xvf /tmp/ci-studio-common.tar.gz
  - export PATH=$PATH:$PWD/ci-studio-common/bin/
```

### Usage

#### Installing Terraform

By default, `install_terraform` installs Terraform `0.8.1`.

```yaml
install:
  - install_terraform 0.8.1
  - export PATH=$PATH:$HOME/tools/terraform-0.8.1
```

## Common `.studiorc`

For Habitat projects, the `.studiorc_common` has some shared/common helpers that you can leverage.

### Installation

Please put the following at the **top** of your `.studiorc` file.

```bash
wget https://raw.githubusercontent.com/chef/ci-studio-common/master/.studiorc_common -O /tmp/.studiorc_common
. /tmp/.studiorc_common
```
