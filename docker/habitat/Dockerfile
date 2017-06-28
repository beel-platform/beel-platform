FROM bluespark/centos7

LABEL name="Bluespark's (Habitat / CentOS 7) base image"

RUN wget -q https://api.bintray.com/content/habitat/stable/linux/x86_64/hab-%24latest-x86_64-linux.tar.gz?bt_package=hab-x86_64-linux -O hab-x86_64-linux.tar.gz
RUN tar xfz hab-x86_64-linux.tar.gz && \
    mv hab-*/hab /usr/bin/ && \
    rm -rf hab-* && \
    hab pkg install core/hab-sup && \
    useradd hab

CMD ["/bin/bash"]
