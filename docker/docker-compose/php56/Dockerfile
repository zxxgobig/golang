FROM php:5.6.36-fpm-stretch
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

ENV TZ=Asia/Shanghai
RUN echo "deb http://mirrors.aliyun.com/debian stretch main contrib non-free" > /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian stretch main contrib non-free" >> /etc/apt/sources.list  && \
    echo "deb http://mirrors.aliyun.com/debian stretch-updates main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian stretch-updates main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free" >> /etc/apt/sources.list && \
    echo "deb-src http://mirrors.aliyun.com/debian-security stretch/updates main contrib non-free" >> /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y iptables \
                          git \
                          unzip \
                          zlib1g-dev \
                          libpng-dev \
                          libfreetype6-dev \
                          libjpeg62-turbo-dev \
                          libpng-dev \
                          libmemcached-dev \
                          libgearman-dev \
                          gperf \
                          libevent-dev \
                          uuid-dev \
                          libxml2-dev \
                          libxslt-dev \
    && pecl install redis-4.0.0 \
    && docker-php-ext-enable redis \
    && curl -fsSL 'https://github.com/msgpack/msgpack-php/archive/msgpack-0.5.4.tar.gz' -o msgpack-0.5.4.tgz\
    && tar zxf msgpack-0.5.4.tgz \
    && rm -rf msgpack-0.5.4.tgz \
    && (\
        cd msgpack-php-msgpack-0.5.4 \
        && phpize \
        && ./configure \
        && make -j "$(nproc)" \
        && make install \
    )\
    && docker-php-ext-enable msgpack \
    && pecl install igbinary-1.2.0 \
    && docker-php-ext-enable igbinary \
    && curl -fsSL 'https://pecl.php.net/get/memcached-2.2.0.tgz' -o memcached-2.2.0.tgz \
    && tar zxf memcached-2.2.0.tgz \
    && rm -rf memcached-2.2.0.tgz \
        && ( \
            cd memcached-2.2.0 \
            && phpize \
            && ./configure  --enable-memcached-igbinary  --enable-memcached-json --enable-memcached-msgpack \
            && make -j "$(nproc)" \
            && make install \
        ) \
    && docker-php-ext-enable memcached \
    && pecl install gearman-1.1.2 \
    && docker-php-ext-enable gearman \
    && curl -fsSL 'http://192.168.20.141/download/zips/php_token_crypt_nocert.zip' -o php_token_crypt_nocert.zip \
    && unzip php_token_crypt_nocert.zip \
    && rm -rf php_token_crypt_nocert.zip \
    && ( \
        cd php_token_crypt_nocert \
        && phpize \
        && ./configure \
        && make -j "$(nproc)" \
        && make install \
    ) \
    && docker-php-ext-enable token_crypt \
    && docker-php-ext-install -j$(nproc) pdo_mysql  \
    && docker-php-ext-install -j$(nproc) bcmath \
    && docker-php-ext-install -j$(nproc) zip \
    && docker-php-ext-install -j$(nproc) sockets \
    && docker-php-ext-install -j$(nproc) wddx \
    && docker-php-ext-install -j$(nproc) xsl \
    && docker-php-ext-install -j$(nproc) opcache \
    && docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
    && docker-php-ext-install -j$(nproc) gd 

WORKDIR /opt

CMD php-fpm