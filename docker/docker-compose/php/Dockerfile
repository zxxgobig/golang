FROM php:7.1.8-fpm
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone \
    && echo "deb http://mirrors.aliyun.com/debian stretch main contrib non-free" > /etc/apt/sources.list \
    &&  echo "deb http://mirrors.aliyun.com/debian stretch-proposed-updates main contrib non-free" >> /etc/apt/sources.list \
    &&  echo "deb http://mirrors.aliyun.com/debian stretch-updates main contrib non-free" >> /etc/apt/sources.list \
    &&  echo "deb http://mirrors.aliyun.com/debian-security/ stretch/updates main non-free contrib">> /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y iptables \
    && apt-get install -y git libmemcached-dev zlib1g-dev libxml2-dev libxslt-dev \
    && pecl install redis \
    && docker-php-ext-enable redis \
    && pecl install msgpack-2.0.3 \
    && docker-php-ext-enable msgpack \
    && pecl install memcached-3.1.3 \
    && docker-php-ext-enable memcached \
    && pecl install igbinary-3.0.1 \
    && docker-php-ext-enable igbinary \
    && docker-php-ext-install -j$(nproc) pdo_mysql  \
    && docker-php-ext-install -j$(nproc) bcmath \
    && docker-php-ext-install -j$(nproc) zip \
    && docker-php-ext-install -j$(nproc) calendar \
    && docker-php-ext-install -j$(nproc) exif \
    && docker-php-ext-install -j$(nproc) gettext  \
    && docker-php-ext-install -j$(nproc) mysqli \
    && docker-php-ext-install -j$(nproc) pcntl \
    && docker-php-ext-install -j$(nproc) shmop \
    && docker-php-ext-install -j$(nproc) sockets \
    && docker-php-ext-install -j$(nproc) sysvmsg \
    && docker-php-ext-install -j$(nproc) sysvsem  \
    && docker-php-ext-install -j$(nproc) sysvshm \
    && docker-php-ext-install -j$(nproc) wddx \
    && docker-php-ext-install -j$(nproc) xsl \
    && docker-php-ext-install -j$(nproc) opcache \
    && apt-get install -y \
        libfreetype6-dev \
        libjpeg62-turbo-dev \
        libpng-dev \
    && docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ \
    && docker-php-ext-install -j$(nproc) gd

WORKDIR /opt

CMD php-fpm
