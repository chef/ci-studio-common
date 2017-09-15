pkg_name=ci-studio-common
pkg_origin=chef
pkg_version="0.1.0"
pkg_maintainer="The Habitat Maintainers <humans@habitat.sh>"
pkg_description="Shared helpers for use inside CIs (like Travis) and a Habitat Studio"
pkg_license=('Apache-2.0')
pkg_bin_dirs=(bin)
pkg_deps=(
  core/busybox
  core/curl
  core/bash
  core/git
)

do_build() {
  return 0
}

do_install() {
  cp -rf $SRC_PATH/bin/* "$pkg_prefix/bin"
}
