#!/bin/bash

set -o errexit
set -o pipefail

if [ -f '/srv/icinga2/secrets/.env' ]; then
    export $(cat /srv/icinga2/secrets/.env | xargs)
else
    echo
    echo 'Missing environment file /srv/icinga2/secrets/.env.'
    echo
    exit 1
fi

if [ ! -f "/scripts/.icingaweb2" ]; then
    envsubst < /scripts/icingaweb2/authentication.ini > /etc/icingaweb2/authentication.ini
    envsubst < /scripts/icingaweb2/config.ini > /etc/icingaweb2/config.ini
    envsubst < /scripts/icingaweb2/groups.ini > /etc/icingaweb2/groups.ini
    envsubst < /scripts/icingaweb2/resources.ini > /etc/icingaweb2/resources.ini
    touch /scripts/.icingaweb2
fi

if [ ! -f "$PGROOT/.lib_icinga2" ]; then
    mv /scripts/lib $PGROOT/icinga2
    touch $PGROOT/.lib_icinga2
fi
chown -R icinga:icinga $PGROOT/icinga2
ln -sv -t /var/lib $PGROOT/icinga2

chown -R icinga:icinga /usr/lib/nagios/plugins
chmod -R 755 /usr/lib/nagios/plugins

mkdir -p $PGSCRIPT
cp /usr/share/icinga2-ido-pgsql/schema/pgsql.sql     $PGSCRIPT/icinga2-ido.schema.sql
cp /usr/share/icingaweb2/etc/schema/pgsql.schema.sql $PGSCRIPT/icingaweb2.schema.sql

cat >$PGSCRIPT/.setup-db.sh <<EOL
#!/bin/bash
set -x

psql -c "CREATE ROLE $ICINGA_IDO_USER WITH LOGIN PASSWORD '$ICINGA_IDO_PASSWORD'";
psql -c "CREATE DATABASE $ICINGA_IDO_DB WITH OWNER $ICINGA_IDO_USER";
psql -U $ICINGA_IDO_USER -d $ICINGA_IDO_DB < \$PGSCRIPT/icinga2-ido.schema.sql;

psql -c "CREATE ROLE $ICINGA_WEB_USER WITH LOGIN PASSWORD '$ICINGA_WEB_PASSWORD'";
psql -c "CREATE DATABASE $ICINGA_WEB_DB WITH OWNER $ICINGA_WEB_USER";
psql -U $ICINGA_WEB_USER -d $ICINGA_WEB_DB < \$PGSCRIPT/icingaweb2.schema.sql;

# Add "Administrators" icingaweb_group; This group has admin permission
psql -d $ICINGA_WEB_DB <<EOF
INSERT INTO icingaweb_group (id, name) VALUES (1, 'Administrators');
EOF
EOL

# Set icingaweb2 UI admin password, if provided
if [ -n "$ICINGA_WEB_ADMIN_PASSWORD" ]; then
    cat >>$PGSCRIPT/.setup-db.sh <<EOL
passhash=\$(openssl passwd -1 "$ICINGA_WEB_ADMIN_PASSWORD")
psql -d $ICINGA_WEB_DB <<EOF
INSERT INTO icingaweb_user (name, active, password_hash) VALUES ('admin', 1, '\$passhash');
INSERT INTO icingaweb_group_membership (group_id, username) VALUES (1, 'admin');
EOF
EOL
fi

chmod 755 $PGSCRIPT/.setup-db.sh
# This line will trigger postgres. So, always keep this as the last line for db setup operations.
mv $PGSCRIPT/.setup-db.sh $PGSCRIPT/setup-db.sh

# IcingaWeb reads namespace and api_endpoint from configmap
mkdir -p /var/run/config/appscode; chmod -R 0755 /var/run/config/appscode

# Wait for postgres to start
# ref: http://unix.stackexchange.com/a/5279
echo "Waiting for postgres to become ready ..."
until pg_isready -h 127.0.0.1 > /dev/null; do echo '.'; sleep 5; done

echo "export KUBERNETES_SERVICE_HOST=${KUBERNETES_SERVICE_HOST}" >> /etc/default/icinga2
echo "export KUBERNETES_SERVICE_PORT=${KUBERNETES_SERVICE_PORT}" >> /etc/default/icinga2

export > /etc/envvars

echo "Starting runit..."
exec /usr/sbin/runsvdir-start
