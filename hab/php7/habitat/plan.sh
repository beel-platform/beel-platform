pkg_name=php7
pkg_origin=bluespark
pkg_distname=php
pkg_version=7.1.6
pkg_maintainer="Basilio Briceno <basilio@bluespark.com>"
pkg_license=('PHP-3.01')
pkg_upstream_url=http://php.net/
pkg_description="PHP is a popular general-purpose scripting language that is especially suited to web development."
pkg_source=https://php.net/get/${pkg_distname}-${pkg_version}.tar.bz2/from/this/mirror
pkg_filename=${pkg_distname}-${pkg_version}.tar.bz2
pkg_dirname=${pkg_distname}-${pkg_version}
pkg_shasum=6e3576ca77672a18461a4b089c5790647f1b2c19f82e4f5e94c962609aabffcf
pkg_deps=(
  core/coreutils
  core/curl
  core/glibc
  core/libxml2
  core/libjpeg-turbo
  core/libpng
  core/openssl
  core/zlib
)
pkg_build_deps=(
  core/bison2
  core/gcc
  core/make
  core/re2c
  bluespark/httpd
)
pkg_bin_dirs=(bin sbin)
pkg_lib_dirs=(lib)
pkg_include_dirs=(include)
pkg_interpreters=(bin/php)
pkg_expose=(8080)
pkg_svc_user=root
pkg_svc_group=root

do_build ()
{
  ./configure --prefix=${pkg_prefix} \
    --bindir="${pkg_prefix}/bin" \
    --sbindir="${pkg_prefix}/sbin" \
    --sysconfdir="${pkg_svc_config_path}/etc" \
    --with-config-file-path="${pkg_svc_config_path}" \
    --enable-exif \
    --enable-fpm \
    --with-fpm-user=hab \
    --with-fpm-group=hab \
    --enable-mbstring \
    --enable-opcache \
    --enable-mysqlnd \
    --with-mysql=mysqlnd \
    --with-mysqli=mysqlnd \
    --with-pdo-mysql=mysqlnd \
    --with-curl="$(pkg_path_for curl)" \
    --with-gd \
    --with-jpeg-dir="$(pkg_path_for libjpeg-turbo)" \
    --with-libxml-dir="$(pkg_path_for libxml2)" \
    --with-openssl="$(pkg_path_for openssl)" \
    --with-png-dir="$(pkg_path_for libpng)" \
    --with-xmlrpc \
    --with-zlib="$(pkg_path_for zlib)" \
    --with-apxs2="$(pkg_path_for bluespark/httpd)/bin/apxs"
  make
}

do_install ()
{
  make install
}

do_check ()
{
  make test
}
