#!/bin/bash

set -x

mkdir -p $PGDATA /docker-entrypoint-initdb.d

echo "Waiting for initdb scripts ..."
until [ -f $PGDATA/../scripts/initdb.sh ] > /dev/null; do echo '.'; sleep 5; done

cp -r $PGDATA/../scripts/initdb.sh /docker-entrypoint-initdb.d/initdb.sh

exec docker-entrypoint.sh "$@"
# https://superuser.com/a/176788/441206
# source docker-entrypoint.sh
