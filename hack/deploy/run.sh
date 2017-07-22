#!/usr/bin/env bash

GOPATH=$(go env GOPATH)
MASTER_IP=$1

ICINGA_WEB_ADMIN_PASSWORD=$(cat /dev/urandom | base64 | tr -d "=+/" | dd bs=16 count=1 2> /dev/null)

read -r -d '' env_data <<EOF
ICINGA_WEB_HOST=127.0.0.1
ICINGA_WEB_PORT=5432
ICINGA_WEB_DB=icingawebdb
ICINGA_WEB_USER=icingaweb
ICINGA_WEB_PASSWORD=$(cat /dev/urandom | base64 | tr -d "=+/" | dd bs=16 count=1 2> /dev/null)
ICINGA_WEB_ADMIN_PASSWORD=$ICINGA_WEB_ADMIN_PASSWORD
ICINGA_IDO_HOST=127.0.0.1
ICINGA_IDO_PORT=5432
ICINGA_IDO_DB=icingaidodb
ICINGA_IDO_USER=icingaido
ICINGA_IDO_PASSWORD=$(cat /dev/urandom | base64 | tr -d "=+/" | dd bs=16 count=1 2> /dev/null)
ICINGA_API_USER=icingaapi
ICINGA_API_PASSWORD=$(cat /dev/urandom | base64 | tr -d "=+/" | dd bs=16 count=1 2> /dev/null)
ICINGA_ADDRESS=searchlight-operator.kube-system
EOF


export ICINGA_SECRET_ENV=$(echo $env_data | base64 -w 0)

# Directory for certificates
certificate_dir=$GOPATH/src/github.com/appscode/searchlight/hack/deploy/certificate
mkdir -p $certificate_dir

pushd $certificate_dir
  # Generate a ca.key with 2048bit:
  openssl genrsa -out ca.key 2048

  # According to the ca.key generate a ca.crt
  openssl req -x509 -new -nodes -key ca.key -subj "/CN=${MASTER_IP}" -days 10000 -out ca.crt

  # Set ICINGA_CA_CERT env to it
  export ICINGA_CA_CERT=$(base64 ca.crt -w 0)

  #Generate a server.key with 2048bit
  openssl genrsa -out server.key 2048

  # Set ICINGA_SERVER_KEY env to it
  export ICINGA_SERVER_KEY=$(base64 server.key -w 0)

  # According to the server.key generate a server.csr
  openssl req -new -key server.key -subj "/CN=${MASTER_IP}" -out server.csr

  # According to the ca.key, ca.crt and server.csr generate the server.crt
  openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 10000

  # Set ICINGA_SERVER_CERT env to it
  export ICINGA_SERVER_CERT=$(base64 server.crt -w 0)

  rm -rf $certificate_dir
popd

# Deploy Secret
curl https://raw.githubusercontent.com/appscode/searchlight/3.0.0/hack/deploy/deployment.yaml \
  | envsubst \
  | kubectl apply -f -

#To login into Icingaweb2, use following authentication information:
echo
echo "To login into Icingaweb2, use following authentication information:"
echo "Username: admin"
echo "Password: $ICINGA_WEB_ADMIN_PASSWORD"
