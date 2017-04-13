FROM bluespark/centos7

LABEL name="Bluespark's (Habitat 0.20.0 / CentOS 7) base image"

RUN wget -q https://api.bintray.com/content/habitat/stable/linux/x86_64/hab-%24latest-x86_64-linux.tar.gz?bt_package=hab-x86_64-linux -O hab-x86_64-linux.tar.gz
RUN tar xfz hab-x86_64-linux.tar.gz
RUN mv hab-*/hab /usr/bin/
RUN rm -rf hab-*
RUN hab pkg install core/hab-sup
RUN useradd hab

EXPOSE 80

CMD ["/bin/bash"]
