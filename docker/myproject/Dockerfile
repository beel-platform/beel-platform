FROM 147905667315.dkr.ecr.us-east-1.amazonaws.com/bluespark:base-hab0.20.0-centos7

LABEL name="Bluespark's LAMP (Apache / PHP7 / MySQL) image"

RUN hab pkg install bluespark/php7
RUN hab pkg binlink bluespark/php7 php
RUN hab pkg binlink bluespark/php7 phar
RUN hab pkg binlink bluespark/php7 php-config
RUN hab pkg binlink bluespark/php7 phpdbg
RUN hab pkg binlink bluespark/php7 phpize
RUN hab sup load bluespark/php7
RUN hab pkg install bluespark/httpd
RUN hab sup load bluespark/httpd
RUN hab pkg install bluespark/mysql
RUN hab sup load bluespark/mysql

EXPOSE 80 3306

ENTRYPOINT hab sup run
