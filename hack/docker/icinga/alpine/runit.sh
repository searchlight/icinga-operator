#!/bin/bash

# Copyright AppsCode Inc. and Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
set -o errexit
set -o pipefail

echo "Waiting for icinga configuration ..."
until [ -f /srv/searchlight/config.ini ] >/dev/null; do
    echo '.'
    sleep 5
    cat /srv/searchlight/config.ini
done
export $(cat /srv/searchlight/config.ini | xargs)

if [ ! -f "/scripts/.icingaweb2" ]; then
    envsubst </scripts/icingaweb2/authentication.ini >/etc/icingaweb2/authentication.ini
    envsubst </scripts/icingaweb2/config.ini >/etc/icingaweb2/config.ini
    envsubst </scripts/icingaweb2/groups.ini >/etc/icingaweb2/groups.ini
    envsubst </scripts/icingaweb2/resources.ini >/etc/icingaweb2/resources.ini
    touch /scripts/.icingaweb2
fi

if [ ! -f "$DATADIR/.lib_icinga2" ]; then
    mv /scripts/lib $DATADIR/icinga2
    touch $DATADIR/.lib_icinga2
fi
chown -R icinga:icinga $DATADIR/icinga2
rm -rf /var/lib/icinga2
ln -sv -T $DATADIR/icinga2 /var/lib/icinga2

mkdir -p $DATADIR/custom.d
chown -R icinga:icinga $DATADIR/custom.d
rm -rf /etc/icinga2/custom.d
ln -sv -T $DATADIR/custom.d /etc/icinga2/custom.d

cat >>/etc/icinga2/icinga2.conf <<EOL
/**
 * Used to add configuration for user defined Commands.
 * Each of these files must have the file extension ".conf".
 */
include_recursive "custom.d"
EOL

chown -R icinga:icinga /usr/lib/nagios/plugins
chmod -R 755 /usr/lib/nagios/plugins

# Fix bash interpreter
# sed -i 's/\/sbin\/openrc-run/\/bin\/bash/g' /etc/init.d/icinga2

mkdir -p $DATADIR/scripts
cp /usr/share/icinga2-ido-pgsql/schema/pgsql.sql $DATADIR/scripts/icinga2-ido.schema.sql
cp /usr/share/icingaweb2/etc/schema/pgsql.schema.sql $DATADIR/scripts/icingaweb2.schema.sql

cat >$DATADIR/scripts/.initdb.sh <<EOL
#!/bin/bash
set -x

psql -U postgres -c "CREATE ROLE $ICINGA_IDO_USER WITH LOGIN PASSWORD '$ICINGA_IDO_PASSWORD'";
psql -U postgres -c "CREATE DATABASE $ICINGA_IDO_DB WITH OWNER $ICINGA_IDO_USER";
psql -U $ICINGA_IDO_USER -d $ICINGA_IDO_DB < \$PGDATA/../scripts/icinga2-ido.schema.sql;

psql -U postgres -c "CREATE ROLE $ICINGA_WEB_USER WITH LOGIN PASSWORD '$ICINGA_WEB_PASSWORD'";
psql -U postgres -c "CREATE DATABASE $ICINGA_WEB_DB WITH OWNER $ICINGA_WEB_USER";
psql -U $ICINGA_WEB_USER -d $ICINGA_WEB_DB < \$PGDATA/../scripts/icingaweb2.schema.sql;

# Add "Administrators" icingaweb_group; This group has admin permission
psql -U $ICINGA_WEB_USER -d $ICINGA_WEB_DB <<EOF
INSERT INTO icingaweb_group (id, name) VALUES (1, 'Administrators');
EOF
EOL

# Set icingaweb2 UI admin password, if provided
if [ -n "$ICINGA_WEB_UI_PASSWORD" ]; then
    cat >>$DATADIR/scripts/.initdb.sh <<EOL
passhash=\$(openssl passwd -1 "$ICINGA_WEB_UI_PASSWORD")
psql -U $ICINGA_WEB_USER -d $ICINGA_WEB_DB <<EOF
INSERT INTO icingaweb_user (name, active, password_hash) VALUES ('admin', 1, '\$passhash');
INSERT INTO icingaweb_group_membership (group_id, username) VALUES (1, 'admin');
EOF
EOL
fi

chmod 755 $DATADIR/scripts/*
# This line will trigger postgres. So, always keep this as the last line for db setup operations.
mv $DATADIR/scripts/.initdb.sh $DATADIR/scripts/initdb.sh

# IcingaWeb reads namespace and api_endpoint from configmap
mkdir -p /var/run/config/appscode
chmod -R 0755 /var/run/config/appscode

# Wait for postgres to start
# ref: http://unix.stackexchange.com/a/5279
echo "Waiting for postgres to become ready ..."
until pg_isready -h 127.0.0.1 >/dev/null; do
    echo '.'
    sleep 5
done

# Ensure icinga plugins can read ENV cars
echo "export KUBERNETES_SERVICE_HOST=${KUBERNETES_SERVICE_HOST}" >>/etc/profile.d/icinga2
echo "export KUBERNETES_SERVICE_PORT=${KUBERNETES_SERVICE_PORT}" >>/etc/profile.d/icinga2
echo "export APPSCODE_ANALYTICS_CLIENT_ID=$(/usr/lib/monitoring-plugins/hyperalert analytics_id)" >>/etc/profile.d/icinga2

export >/etc/envvars

echo "Starting runit..."
exec /sbin/runsvdir -P /etc/service
