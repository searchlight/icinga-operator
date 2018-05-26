FROM alpine:3.6

COPY runit.sh /runit.sh

ENV DATADIR /srv

# Install Icinga2 ###################################
RUN set -x \
  && apk add --update --no-cache ca-certificates openssl su-exec runit socklog tzdata curl nginx gettext bash openrc shadow iputils \
    icinga2 nagios-plugins jq \
    postgresql \
    php7 \
    php7-common \
    php7-ctype \
    php7-dom \
    php7-fpm \
    php7-gd \
    php7-gettext \
    php7-gmp \
    php7-imagick \
    php7-intl \
    php7-json \
    php7-ldap \
    php7-mcrypt \
    php7-mysqlnd \
    php7-openssl \
    php7-pdo \
    php7-pdo_mysql \
    php7-pdo_pgsql \
    php7-pgsql \
    php7-session \
    php7-sockets

RUN set -x \
  && mkdir -p /run/nginx \
  && chown nginx:nginx -R /run/nginx \
  && sed -i 's/group.*/group\ =\ icingacmd/' /etc/php7/php-fpm.d/www.conf \
  && sed -i 's/;date.timezone =/date.timezone ="UTC"/' /etc/php7/php.ini \
  && sed -i ';s/error_log.*/error_log\ =\ syslog/' /etc/php7/php-fpm.conf \
  && sed -i 's/;syslog.facility/syslog.facility/' /etc/php7/php-fpm.conf \
  && sed -i 's/;syslog.ident/syslog.ident/' /etc/php7/php-fpm.conf \
  && rm -rf /etc/sv /etc/service \
  && echo 'Etc/UTC' > /etc/timezone

# copy config templates
RUN set -x && mkdir -p /scripts/icinga2
COPY config/icinga2/ido-pgsql.conf        /scripts/icinga2/ido-pgsql.conf

# Icinga 2 IDO
# http://docs.icinga.org/icinga2/latest/doc/module/icinga2/chapter/icinga2-features
RUN set -x && icinga2 feature enable ido-pgsql syslog checker

# Command Feature
RUN set -x \
  && icinga2 feature enable command notification \
  && mkdir -p /run/icinga2/cmd \
  && chown icinga:icinga -R /run/icinga2

COPY config/icinga2/notification-command.conf /scripts/notification-command.conf
RUN set -x \
  && cat /scripts/notification-command.conf >> /etc/icinga2/conf.d/notification-command.conf \
  && rm /scripts/notification-command.conf

COPY config/icinga2/modified-commands.conf /scripts/modified-commands.conf
RUN set -x \
  && cat /scripts/modified-commands.conf >> /etc/icinga2/conf.d/modified-commands.conf \
  && rm /scripts/modified-commands.conf

COPY config/icinga2/templates.conf /scripts/templates.conf
RUN set -x \
  && cat /scripts/templates.conf >> /etc/icinga2/conf.d/templates.conf \
  && rm /scripts/templates.conf

COPY config/icinga2/users.conf /scripts/users.conf
RUN set -x \
  && cat /scripts/users.conf >> /etc/icinga2/conf.d/users.conf \
  && rm /scripts/users.conf

RUN set -x && rm -rf /etc/icinga2/conf.d/hosts.conf

# Fix icinga installation location and permission
# This is needed since we are building from source
RUN set -x \
  && mv /var/lib/icinga2 /scripts/lib

# Icingaweb2 ##############################################

# copy config templates
RUN set -x && mkdir -p /scripts/icingaweb2
COPY config/icingaweb2/authentication.ini /scripts/icingaweb2/authentication.ini
COPY config/icingaweb2/config.ini         /scripts/icingaweb2/config.ini
COPY config/icingaweb2/groups.ini         /scripts/icingaweb2/groups.ini
COPY config/icingaweb2/resources.ini      /scripts/icingaweb2/resources.ini

# Add icingaweb2
COPY icingaweb2 /usr/share/icingaweb2
RUN set -x && chown -R nginx:nginx /usr/share/icingaweb2
COPY config/icingaweb2 /etc/icingaweb2/
RUN set -x \
  && mkdir -p /etc/icingaweb2/enabledModules \
  && ln -s /usr/share/icingaweb2/modules/doc        /etc/icingaweb2/enabledModules/doc \
	&& ln -s /usr/share/icingaweb2/modules/monitoring /etc/icingaweb2/enabledModules/monitoring \
	&& ln -s /usr/share/icingaweb2/modules/test       /etc/icingaweb2/enabledModules/test

# Update nginx site configuraiton
COPY config/nginx.conf /etc/nginx/conf.d/default.conf

# Plugins ############################################
COPY plugins/* /usr/lib/monitoring-plugins/

# runit ##############################################
ADD sv /etc/sv/
RUN ln -s /etc/sv /etc/service

COPY sv /etc/sv/

ENV TZ     :/etc/localtime
ENV LANG   en_US.utf8

VOLUME ["$DATADIR"]

ENTRYPOINT ["/runit.sh"]
EXPOSE  60006 5665
