FROM php:8.3-apache-bookworm

LABEL maintainer="yayansty"

ENV DEBIAN_FRONTEND=noninteractive

# =====================================================
# Install System Packages
# =====================================================

RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    git \
    curl \
    wget \
    unzip \
    zip \
    nano \
    vim \
    gnupg2 \
    apt-transport-https \
    ca-certificates \
    unixodbc \
    unixodbc-dev \
    libzip-dev \
    libicu-dev \
    libpng-dev \
    libjpeg62-turbo-dev \
    libfreetype6-dev \
    libxml2-dev \
    libxslt1-dev \
    libonig-dev \
    libldap2-dev \
    libssl-dev \
    libcurl4-openssl-dev \
    libgmp-dev \
    && rm -rf /var/lib/apt/lists/*

# =====================================================
# Install PHP Extensions
# =====================================================

RUN docker-php-ext-configure gd \
    --with-freetype \
    --with-jpeg

RUN docker-php-ext-install -j$(nproc) \
    bcmath \
    exif \
    gd \
    gettext \
    intl \
    mbstring \
    opcache \
    pcntl \
    pdo \
    soap \
    sockets \
    xsl \
    gmp \
    zip

# =====================================================
# Apache Modules
# =====================================================

RUN a2enmod rewrite headers expires remoteip

# =====================================================
# Microsoft SQL Server Driver
# =====================================================

RUN curl -fsSL https://packages.microsoft.com/keys/microsoft.asc \
    | gpg --dearmor \
    -o /usr/share/keyrings/microsoft-prod.gpg

RUN echo "deb [arch=amd64 signed-by=/usr/share/keyrings/microsoft-prod.gpg] https://packages.microsoft.com/debian/12/prod bookworm main" \
    > /etc/apt/sources.list.d/mssql-release.list

RUN apt-get update && ACCEPT_EULA=Y apt-get install -y \
    msodbcsql18 \
    mssql-tools18 \
    && rm -rf /var/lib/apt/lists/*

# =====================================================
# SQLSRV Extension
# =====================================================

RUN apt-get update \
    && apt-get install -y --no-install-recommends $PHPIZE_DEPS \
    && pecl install sqlsrv pdo_sqlsrv \
    && docker-php-ext-enable sqlsrv pdo_sqlsrv \
    && apt-get purge -y --auto-remove $PHPIZE_DEPS \
    && rm -rf /var/lib/apt/lists/* /tmp/pear

# =====================================================
# Composer
# =====================================================

COPY --from=composer:2 /usr/bin/composer /usr/bin/composer

# =====================================================
# PHP Configuration
# =====================================================

COPY php.ini /usr/local/etc/php/php.ini

# =====================================================
# Apache Virtual Host
# =====================================================

COPY vhost.conf \
     /etc/apache2/sites-available/000-default.conf

# =====================================================
# Healthcheck
# =====================================================

HEALTHCHECK \
    --interval=30s \
    --timeout=5s \
    --start-period=10s \
    --retries=3 \
CMD curl -f http://localhost/ || exit 1

WORKDIR /var/www/html


COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/entrypoint.sh

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]

CMD ["apache2-foreground"]
