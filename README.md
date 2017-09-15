# Common CI / Habitat Studio Functionality

This repository houses some scripts / files that are used across various Chef projects. These are designed and intended to be used by Chef Developers who are working on some specific projects. We are keeping it public as to simplify the download process for our developers.

## TravisCI Scripts

These are scripts that can be used inside your `.travis.yml` files.

### Installation

```yaml
before_install:
  - curl https://raw.githubusercontent.com/habitat-sh/habitat/master/components/hab/install.sh | sudo bash
  - sudo hab pkg install chef/ci-studio-common --binlink
```

## Habitat Studio

### Installation

Please put the following at the **top** of your `.studiorc` file.

```bash
hab pkg install chef/ci-studio-common --binlink
source studio-common
```
