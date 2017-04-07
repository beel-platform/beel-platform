FROM bluespark/habitat

LABEL name="Bluespark's Web (Apache/PHP7) image"

RUN hab pkg install bluespark/php7
RUN hab pkg binlink bluespark/php7 php
RUN hab pkg binlink bluespark/php7 phar
RUN hab pkg binlink bluespark/php7 php-config
RUN hab pkg binlink bluespark/php7 phpdbg
RUN hab pkg binlink bluespark/php7 phpize
RUN hab pkg install bluespark/httpd

EXPOSE 80

ENTRYPOINT hab sup start bluespark/httpd && /bin/bash
