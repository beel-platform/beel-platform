FROM 147905667315.dkr.ecr.us-east-1.amazonaws.com/bluespark:base-hab0.20.0-centos7

LABEL name="Bluespark's LAMP (Apache / PHP7 / MySQL) image"

RUN hab pkg install bluespark/php7; \
    hab pkg binlink bluespark/php7 php; \
    hab pkg binlink bluespark/php7 phar; \
    hab pkg binlink bluespark/php7 php-config; \
    hab pkg binlink bluespark/php7 phpdbg; \
    hab pkg binlink bluespark/php7 phpize; \
    hab sup load bluespark/php7; \
    hab pkg install bluespark/httpd; \
    hab sup load bluespark/httpd; \
    hab pkg install bluespark/mysql; \
    hab sup load bluespark/mysql

EXPOSE 80 3306

ENTRYPOINT hab sup run
