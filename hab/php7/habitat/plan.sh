pkg_name=php7
pkg_origin=bluespark
pkg_distname=php
pkg_version=7.1.2
pkg_maintainer="Basilio Briceno <basilio@bluespark.com>"
pkg_license=('PHP-3.01')
pkg_upstream_url=http://php.net/
pkg_description="PHP is a popular general-purpose scripting language that is especially suited to web development."
pkg_source=https://php.net/get/${pkg_distname}-${pkg_version}.tar.bz2/from/this/mirror
pkg_filename=${pkg_distname}-${pkg_version}.tar.bz2
pkg_dirname=${pkg_distname}-${pkg_version}
pkg_shasum=e0f2214e2366434ee231156ba70cfefd0c59790f050d8727a3f1dc2affa67004
pkg_deps=(core/libxml2 core/curl core/libpng core/libjpeg-turbo core/zlib core/openssl)
pkg_build_deps=(core/bison2 core/gcc core/make core/re2c core/m4 core/pkg-config bluespark/httpd)
pkg_sbin_dirs=(sbin)
pkg_bin_dirs=(bin)
pkg_lib_dirs=(lib)
pkg_include_dirs=(include)
pkg_interpreters=(bin/php)
pkg_expose=(8080)
pkg_svc_user=root
pkg_svc_group=root

do_prepare ()
{
  if [[ ! -r /usr/bin/xml2-config ]]; then
    ln -sv "$(pkg_path_for core/libxml2)/bin/xml2-config" /usr/bin/xml2-config
    _clean_xml2=true
  fi
  if [[ ! -r /usr/include/openssl/evp.h ]]; then
    mkdir -p /usr/include/openssl
    ln -sv "$(pkg_path_for core/openssl)/include/openssl/evp.h" /usr/include/openssl/evp.h
    _clean_openssl=true
  fi
  if [[ ! -r /usr/include/curl/easy.h ]]; then
    ln -sv "$(pkg_path_for core/curl)/include/curl" /usr/include/curl
    _clean_curl=true
  fi
  if [[ ! -r /usr/include/jpeglib.h ]]; then
    ln -sv "$(pkg_path_for core/libjpeg-turbo)/include/jpeglib.h" /usr/include/jpeglib.h
    _clean_libjpeg=true
  fi
  if [[ ! -r /usr/include/png.h ]]; then
    ln -sv "$(pkg_path_for core/libpng)/include/png.h" /usr/include/png.h
    _clean_libpng=true
  fi
  if [[ ! -r /usr/include/zlib.h ]]; then
    ln -sv "$(pkg_path_for core/zlib)/include/zlib.h" /usr/include/zlib.h
    _clean_zlib=true
  fi
}

do_build ()
{
  ./configure --prefix=${pkg_prefix} \
    CFLAGS="$CFLAGS" \
    LDFLAGS="$LDFLAGS" \
    --bindir="${pkg_prefix}/bin" \
    --sbindir="${pkg_prefix}/sbin" \
    --sysconfdir="${pkg_svc_config_path}/etc" \
    --with-config-file-path="${pkg_svc_config_path}" \
    --enable-fpm \
    --enable-embedded-mysqli \
    --enable-mysqlnd \
    --with-pdo-mysql \
    --enable-opcache \
    --without-pear \
    --with-gd \
    --with-curl \
    --with-jpeg-dir \
    --with-zlib-dir \
    --with-openssl \
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

do_end ()
{
  if [[ -n "$_clean_xml2" ]]; then
    rm -fv /usr/bin/xml2-config
  fi
  if [[ -n "$_clean_openssl" ]]; then
    rm -frv /usr/include/openssl
  fi
  if [[ -n "$_clean_curl" ]]; then
    rm -frv /usr/include/curl
  fi
  if [[ -n "$_clean_libjpeg" ]]; then
    rm -fv /usr/include/jpeglib.h
  fi
  if [[ -n "$_clean_libpng" ]]; then
    rm -fv /usr/include/png.h
  fi
  if [[ -n "$_clean_zlib" ]]; then
    rm -fv /usr/include/zlib.h
  fi
}
