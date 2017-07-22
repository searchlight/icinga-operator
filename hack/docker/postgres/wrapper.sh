#!/bin/bash

echo "Waiting for icinga configuration ..."
until [ -f /srv/icinga2/config.ini ] > /dev/null; do echo '.'; sleep 5; cat /srv/icinga2/config.ini; done
export $(cat /srv/icinga2/config.ini | xargs)

cp -r $PGDATA/../docker-entrypoint-initdb.d /docker-entrypoint-initdb.d

# exec docker-entrypoint.sh "$@"
# https://superuser.com/a/176788/441206
source docker-entrypoint.sh
