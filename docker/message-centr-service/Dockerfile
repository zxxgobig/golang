FROM hub.ifchange.com/bpc/php-from-offical:7.4-fpm
ENV PHP_FPM="/usr/local/sbin/php-fpm" PHP_CLI="/usr/local/bin/php" CRONTAB_START="/etc/init.d/cron start"
ENV PROCESS_FPM_WEB=0 PROCESS_FPM_RPC=0 PROCESS_QUEUE=1 PROCESS_CRONTAB=0
ENV APP_NAME="jsb-message-center-service" APP_PATH="/opt/wwwroot/jsb-message-center-service/" FRAMEWORK="laravel" SPI_NAME="bpc-message-center-service"
WORKDIR /opt/wwwroot/jsb-message-center-service/
COPY docker-config/default /opt/docker-config
COPY ./ /opt/wwwroot/jsb-message-center-service/
RUN mkdir -p /opt/php_cache && chown -R nobody:nobody /opt/php_cache && mkdir -p /opt/log
CMD ["/opt/docker-config/run.sh"]