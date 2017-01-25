# Deployment Guide

This guide will walk you through deploying the icinga2.

### Deploy Icinga

###### Deploy Secret

We need to create secret object for Icinga2. We Need following data for secret object

1. .env: `$ICINGA_SECRET_ENV`
2. ca.crt: `$ICINGA_CA_CERT`
3. icinga.key: `$ICINGA_SERVER_KEY`
4. icinga.crt: `$ICINGA_SERVER_CERT` 


Save the following contents to `secret.ini`:
```ini
ICINGA_WEB_HOST=127.0.0.1
ICINGA_WEB_PORT=5432
ICINGA_WEB_DB=icingawebdb
ICINGA_WEB_USER=icingaweb
ICINGA_WEB_PASSWORD=12345678
ICINGA_WEB_ADMIN_PASSWORD=admin
ICINGA_IDO_HOST=127.0.0.1
ICINGA_IDO_PORT=5432
ICINGA_IDO_DB=icingaidodb
ICINGA_IDO_USER=icingaido
ICINGA_IDO_PASSWORD=12345678
ICINGA_API_USER=icingaapi
ICINGA_API_PASSWORD=12345678
ICINGA_SERVICE=k8s-icinga
```

Encode Secret data and set `ICINGA_SECRET_ENV` to it
```sh
set ICINGA_SECRET_ENV (base64 secret.ini -w 0)
```


We need to generate Icinga2 API certificates. See [here](certificate.md)

Substitute ENV and deploy secret
```sh
# Deploy Secret
curl https://raw.githubusercontent.com/appscode/searchlight/master/hack/kubernetes/icinga2/secret.yaml |
envsubst | kubectl apply -f -
```

###### Create Service
```sh
# Create Service
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/master/hack/kubernetes/icinga2/service.yaml
```

###### Create Deployment
```sh
# Create Deployment
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/master/hack/kubernetes/icinga2/deployment.yaml
```

### Login

To login into `Icingaweb2`, use following authentication information:
```
Username: admin
Password: <ICINGA_WEB_ADMIN_PASSWORD>
```
Password will be set from Icinga secret.
