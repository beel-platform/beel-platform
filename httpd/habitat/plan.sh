pkg_name=httpd
pkg_origin=bbh
pkg_version=2.4.25
pkg_description="The Apache HTTP Server"
pkg_upstream_url="http://httpd.apache.org/"
pkg_maintainer="Basilio Briceno <bbh@briceno.mx>"
pkg_license=('Apache-2.0')
pkg_source="https://archive.apache.org/dist/${pkg_name}/${pkg_name}-${pkg_version}.tar.gz"
pkg_shasum=be6c5eb805216ec205453bb02b1990c82609cb1b145bcb69dc6e99fff45493a9
pkg_deps=(core/glibc core/expat core/libiconv core/apr core/apr-util core/pcre core/zlib core/openssl core/gcc-libs)
pkg_build_deps=(core/patch core/make core/gcc)
pkg_bin_dirs=(bin)
pkg_include_dirs=(include)
pkg_lib_dirs=(lib)
pkg_exports=([port]=serverport)
pkg_exposes=(port)
pkg_svc_run="httpd -DFOREGROUND -f ${pkg_svc_config_path}/httpd.conf"
pkg_svc_user="root"
pkg_svc_group="root"

do_build ()
{
  # --datadir=${pkg_svc_var_path} \
  ./configure --prefix="${pkg_prefix}" \
    --sysconfdir="${pkg_prefix}/etc" \
    --with-expat="$(pkg_path_for core/expat)" \
    --with-iconv="$(pkg_path_for core/libiconv)" \
    --with-pcre="$(pkg_path_for core/pcre)" \
    --with-apr="$(pkg_path_for core/apr)" \
    --with-apr-util="$(pkg_path_for core/apr-util)" \
    --with-z="$(pkg_path_for core/zlib)" \
    --with-ssl="$(pkg_path_for core/openssl)" \
    --enable-modules="none" \
    --enable-mods-static="none" \
    --enable-mods-shared="reallyall" \
    --enable-mpms-shared="prefork event worker"
  make -j4
}

do_install ()
{
  make install
  cp -rv ${PLAN_CONTEXT}/../source/* ${pkg_prefix}/htdocs/
}
