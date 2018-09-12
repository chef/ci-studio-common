pkg_name=ci-studio-common
pkg_origin=chef
pkg_maintainer="Engineering Services <eng-services@chef.io>"
pkg_description="Shared helpers for use inside CIs (like Travis) and a Habitat Studio"
pkg_license=('Apache-2.0')
pkg_upstream_url=https://github.com/chef/ci-studio-common
pkg_deps=core/curl
pkg_bin_dirs=(bin)

pkg_version() {
  cat "$SRC_PATH/VERSION"
}

do_before() {
  do_default_before
  update_pkg_version
}

do_build() {
  return 0
}

do_install() {
  cp -rf "$SRC_PATH"/bin/* "$pkg_prefix/bin"
  cp -rf "$SRC_PATH/dot-studio" "$pkg_prefix"
}
