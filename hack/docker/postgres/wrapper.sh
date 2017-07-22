#!/bin/bash

echo "Waiting for initdb scripts ..."
until [ -f $PGDATA/../dbscripts/initdb.sh ] > /dev/null; do echo '.'; sleep 5; done

cp -r $PGDATA/../dbscripts/initdb.sh /docker-entrypoint-initdb.d/initdb.sh

# exec docker-entrypoint.sh "$@"
# https://superuser.com/a/176788/441206
source docker-entrypoint.sh
