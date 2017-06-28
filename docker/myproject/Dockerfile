FROM 147905667315.dkr.ecr.us-west-1.amazonaws.com/bluespark:base_centos7-hab0.24.1

LABEL name="Bluespark's (Apache / PHP7) container image"

RUN hab pkg install bluespark/php7 && \
    hab sup load bluespark/php7 && \
    hab pkg binlink bluespark/php7 php; \
    hab pkg binlink bluespark/php7 phar; \
    hab pkg binlink bluespark/php7 php-config; \
    hab pkg binlink bluespark/php7 phpdbg; \
    hab pkg binlink bluespark/php7 phpize; \
    hab pkg install bluespark/httpd && \
    hab sup load bluespark/httpd;

EXPOSE 80

ENTRYPOINT /bin/bash

CMD hab sup run
