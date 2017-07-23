FROM appscode/nginx-php:1.10-5

ENV PGROOT /var/pv
ENV PGSCRIPT /var/db-script

RUN set -x \
  && groupadd -r icinga --gid=888 \
  && groupadd -r icingacmd --gid=777 \
  && useradd -r -G icingacmd -g icinga --uid=888 icinga \
  && usermod -a -G icingacmd www-data \
  && mkdir -p /var/cache/icinga2 \
  && chown -R icinga:icinga /var/cache/icinga2

# Install compilers ###################################
RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates netcat wget curl \
    build-essential \
    libssl-dev \
    libboost-all-dev \
    bison \
    flex \
    libpq-dev \
    libyajl-dev \
    libedit-dev \
    postgresql-client \
    jq

RUN set -x \
  && wget http://www.cmake.org/files/v3.2/cmake-3.2.2.tar.gz \
  && tar xf cmake-3.2.2.tar.gz \
  && cd cmake-3.2.2 \
  && ./configure \
  && make \
  && make install \
  && cd ..

# Icing2 ##############################################
RUN set -x \
  && wget https://github.com/Icinga/icinga2/archive/v2.4.8.tar.gz \
  && tar xf v2.4.8.tar.gz \
  && cd icinga2-2.4.8 \
  && mkdir build && cd build \
  && cmake .. -DICINGA2_WITH_MYSQL=OFF -DCMAKE_INSTALL_PREFIX= \
  && make \
  && make install \
  && cd ../..

# Nagios-plugins ####################################
# https://wiki.icinga.org/display/howtos/Nagios+Plugins
RUN set -x \
  && wget https://www.monitoring-plugins.org/download/monitoring-plugins-2.1.2.tar.gz \
  && tar xf monitoring-plugins-2.1.2.tar.gz \
  && cd monitoring-plugins-2.1.2 \
  && ./configure --prefix=/usr/lib/nagios \
    --with-nagios-user=icinga \
    --with-nagios-group=icinga \
  && make \
  && make install \
  && mv /usr/lib/nagios/libexec /usr/lib/nagios/plugins \
  && cd ..

RUN set -x \
  && rm -rf /monitoring-plugins-2.1.2 /monitoring-plugins-2.1.2.tar.gz \
  && rm -rf /icinga2-2.4.8 /v2.4.8.tar.gz \
  && rm -rf /cmake-3.2.2 /cmake-3.2.2.tar.gz \
  && apt-get purge -y --auto-remove wget build-essential \
  && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

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

RUN set -x && rm /etc/icinga2/conf.d/hosts.conf

# Fix icinga installation location and permission
# This is needed since we are building from source
RUN set -x \
  && mv /share/icinga2-ido-pgsql /usr/share \
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
RUN set -x && chown -R www-data:www-data /usr/share/icingaweb2
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
EXPOSE  80 5665
