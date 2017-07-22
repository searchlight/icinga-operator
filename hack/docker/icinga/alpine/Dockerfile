FROM appscode/alpine:3.6

ENV DATADIR /srv

# Install Icinga2 ###################################
RUN set -x \
  && apk add --update --no-cache ca-certificates openssl curl nginx gettext bash \
    icinga2 nagios-plugins \
    postgresql \
    jq \
    php7-fpm

# copy config templates
RUN set -x && mkdir -p /scripts/icinga2
COPY config/icinga2/ido-pgsql.conf        /scripts/icinga2/ido-pgsql.conf

# Icinga 2 IDO
# http://docs.icinga.org/icinga2/latest/doc/module/icinga2/chapter/icinga2-features
RUN set -x && icinga2 feature enable ido-pgsql syslog

# Command Feature
COPY config/icinga2/appscode-commands.conf /scripts/appscode-commands.conf
COPY config/icinga2/modified-commands.conf /scripts/modified-commands.conf
RUN set -x \
  && icinga2 feature enable command \
  && /usr/bin/install -d -o icinga /run/icinga2/cmd \
  && cat /scripts/*-commands.conf >> /etc/icinga2/conf.d/commands.conf \
  && rm /scripts/*-commands.conf

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
COPY plugins/* /usr/lib/nagios/plugins/

# runit ##############################################
ADD sv /etc/sv/
RUN ln -s /etc/sv /etc/service

COPY runit.sh /runit.sh
ENTRYPOINT ["/runit.sh"]
EXPOSE  60006 5665
