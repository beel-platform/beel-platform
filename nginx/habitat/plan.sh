pkg_origin=bbh
pkg_name=nginx
pkg_version=1.11.9
pkg_description="NGINX web server."
pkg_maintainer="Basilio Briceno <bbh@briceno.mx>"
pkg_license=('bsd')
pkg_source=https://nginx.org/download/nginx-${pkg_version}.tar.gz
pkg_upstream_url=https://nginx.org/
pkg_shasum=dc22b71f16b551705930544dc042f1ad1af2f9715f565187ec22c7a4b2625748
pkg_deps=(core/glibc core/libedit core/ncurses core/zlib core/bzip2 core/openssl core/pcre)
pkg_build_deps=(core/gcc core/make core/coreutils)
pkg_lib_dirs=(lib)
pkg_bin_dirs=(sbin)
pkg_include_dirs=(include)
pkg_expose=(80 443)
pkg_svc_user=root
pkg_svc_group=root

do_build ()
{
  ./configure --prefix="${pkg_prefix}" \
    --with-cc-opt="$CFLAGS" \
    --with-ld-opt="$LDFLAGS" \
    --sbin-path="${pkg_prefix}/sbin/nginx" \
    --conf-path="${pkg_svc_config_path}/nginx.conf" \
    --lock-path="${pkg_svc_var_path}/nginx.lock" \
    --pid-path="${pkg_svc_var_path}/nginx.pid" \
    --http-client-body-temp-path="${pkg_svc_var_path}/client-body" \
    --http-fastcgi-temp-path="${pkg_svc_var_path}/fastcgi" \
    --http-proxy-temp-path="${pkg_svc_var_path}/proxy" \
    --http-uwsgi-temp-path="${pkg_svc_var_path}/uwsgi" \
    --http-scgi-temp-path="${pkg_svc_var_path}/scgi" \
    --user=hab \
    --group=hab \
    --http-log-path=/dev/stdout \
    --error-log-path=stderr \
    --with-ipv6 \
    --with-pcre \
    --with-pcre-jit \
    --with-file-aio \
    --with-stream=dynamic \
    --with-mail=dynamic \
    --with-http_gunzip_module \
    --with-http_gzip_static_module \
    --with-http_realip_module \
    --with-http_v2_module \
    --with-http_ssl_module \
    --with-http_stub_status_module \
    --with-http_addition_module \
    --with-http_degradation_module \
    --with-http_flv_module \
    --with-http_mp4_module \
    --with-http_secure_link_module \
    --with-http_sub_module \
    --with-http_slice_module
  make -j4
  # cp -v ${HAB_CACHE_SRC_PATH}/${pkg_dirname}/objs/nginx ${pkg_prefix}/sbin/
  mkdir -p ${pkg_prefix}/www
  cp -rv ${PLAN_CONTEXT}/../source/* ${pkg_prefix}/www
}

do_install ()
{
  make install
}
