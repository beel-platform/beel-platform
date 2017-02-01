pkg_origin=bbh
pkg_name=php5
pkg_distname=php
pkg_version=5.6.29
pkg_maintainer="Basilio Briceno <bbh@briceno.mx>"
pkg_license=('PHP-3.01')
pkg_upstream_url=http://php.net/
pkg_description="PHP is a popular general-purpose scripting language that is especially suited to web development."
pkg_source=https://php.net/get/${pkg_distname}-${pkg_version}.tar.bz2/from/this/mirror
pkg_filename=${pkg_distname}-${pkg_version}.tar.bz2
pkg_dirname=${pkg_distname}-${pkg_version}
pkg_shasum=499b844c8aa7be064c111692e51a093ba94e54d2d9abb01e70ea76183a1825bb
pkg_deps=(core/libxml2 core/openssl core/curl core/libpng core/libjpeg-turbo core/zlib)
pkg_build_deps=(core/bison2 core/gcc core/make core/re2c)
pkg_sbin_dirs=(sbin)
pkg_bin_dirs=(bin)
pkg_lib_dirs=(lib)
pkg_include_dirs=(include)
pkg_interpreters=(bin/php)
pkg_expose=(8080)
pkg_svc_user=root
pkg_svc_group=root

do_prepare() {
  # The configure script expects binaries to either be in `/usr/bin`, `/usr/local/bin` or be
  # passed in as a configure param. Instead of overriding the entire do_build, symlink the
  # required executable into place.
  if [[ ! -r /usr/bin/xml2-config ]]; then
    ln -sv "$(pkg_path_for libxml2)/bin/xml2-config" /usr/bin/xml2-config
    _clean_xml2=true
  fi
  # if [[ ! -r /usr/include/openssl/evp.h ]]; then
  #   ln -sv "$(pkg_path_for openssl)/include/openssl" /usr/include/openssl
  #   _clean_openssl=true
  # fi
  if [[ ! -r /usr/include/curl/easy.h ]]; then
    ln -sv "$(pkg_path_for curl)/include/curl" /usr/include/curl
    _clean_curl=true
  fi
  if [[ ! -r /usr/include/jpeglib.h ]]; then
    ln -sv "$(pkg_path_for libjpeg-turbo)/include/jpeglib.h" /usr/include/jpeglib.h
    _clean_libjpeg=true
  fi
  if [[ ! -r /usr/include/png.h ]]; then
    ln -sv "$(pkg_path_for libpng)/include/png.h" /usr/include/png.h
    _clean_libpng=true
  fi
  if [[ ! -r /usr/include/zlib.h ]]; then
    ln -sv "$(pkg_path_for zlib)/include/zlib.h" /usr/include/zlib.h
    _clean_zlib=true
  fi
}

do_build ()
{
  ./configure --prefix=$pkg_prefix \
    CFLAGS="$CFLAGS" \
    LDFLAGS="$LDFLAGS" \
    CPPFLAGS="$CPPFLAGS" \
    CXXFLAGS="$CXXFLAGS" \
    --bindir="$pkg_prefix/bin" \
    --sbindir="$pkg_prefix/sbin" \
    --sysconfdir="$pkg_prefix/etc" \
    --with-config-file-path="$pkg_prefix/config" \
    --enable-fpm \
    --with-mysql=mysqlnd \
    --with-pdo-mysql=mysqlnd \
    --enable-opcache \
    --without-pear \
    --with-gd \
    --with-curl \
    --with-jpeg-dir \
    --with-zlib-dir
    # --with-openssl=/usr
  make -j8
}

do_install ()
{
  make -j8 install
  if [[ -e $pkg_prefix/etc/php-fpm.conf ]]; then rm -fv $pkg_prefix/etc/php-fpm.conf; fi
  ln -sv $pkg_prefix/config/php-fpm.conf $pkg_prefix/etc/php-fpm.conf
}

do_check()
{
  make -j8 test
}

do_end()
{
  # Clean up the `xml2-config` link, if we set it up.
  if [[ -n "$_clean_xml2" ]]; then
    rm -fv /usr/bin/xml2-config
  fi
  # if [[ -n "$_clean_openssl" ]]; then
  #   rm -frv /usr/include/openssl
  # fi
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
