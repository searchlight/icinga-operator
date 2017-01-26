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
ICINGA_K8S_SERVICE=k8s-icinga
```

Encode Secret data and set `ICINGA_SECRET_ENV` to it
```sh
set ICINGA_SECRET_ENV (base64 secret.ini -w 0)
```


We need to generate Icinga2 API certificates. See [here](certificate.md)

Substitute ENV and deploy secret
```sh
# Deploy Secret
curl https://raw.githubusercontent.com/appscode/searchlight/master/hack/deploy/icinga2/secret.yaml |
envsubst | kubectl apply -f -
```

###### Create Service
```sh
# Create Service
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/master/hack/deploy/icinga2/service.yaml
```

###### Create Deployment

We need to configure notifier, if we want, by setting some `ENV` in deployment. We are currently supporting following notifiers. Set `ENV` for selected notifier in deployment.

1. [Hipchat](../notifier/hipchat.md)
2. [Mailgun](../notifier/mailgun.md)
3. [SMTP](../notifier/smtp.md)
4. [Twilio](../notifier/twilio.md)

Now we can create deployment. If we don't set notifier `ENV`, notifications will be ignored.

```sh
# Create Deployment
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/master/hack/deploy/icinga2/deployment.yaml
```

### Login

To login into `Icingaweb2`, use following authentication information:
```
Username: admin
Password: <ICINGA_WEB_ADMIN_PASSWORD>
```
Password will be set from Icinga secret.
