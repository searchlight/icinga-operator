FROM postgres:9.5-alpine

RUN set -x \
  && apk add --update --no-cache openssl

COPY docker-entrypoint.sh /usr/local/bin/
