#!/bin/bash

echo "Waiting for icinga configuration ..."
until [ -f /srv/icinga2/config ] > /dev/null; do echo '.'; sleep 5; cat /srv/icinga2/config; done
export $(cat /srv/icinga2/config | xargs)

cp -r $PGDATA/../docker-entrypoint-initdb.d /docker-entrypoint-initdb.d

exec docker-entrypoint.sh "$@"
