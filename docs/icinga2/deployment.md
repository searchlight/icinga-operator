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
ICINGA_ADDRESS=searchlight-icinga.kube-system
```

We can use following as `ICINGA_ADDRESS`:

* `<HostIP>:<HostPort>`
* `<KubernetesService.KubernetesNamespace>:<ServicePort>`

> Port is optional. Default: 5665
> KubernetesNamespace is optional. Default: default

Encode Secret data and set `ICINGA_SECRET_ENV` to it
```sh
export ICINGA_SECRET_ENV=$(base64 secret.ini -w 0)
```

We need to generate Icinga2 API certificates. See [here](certificate.md)

And also we need to add some keys for notifier in Icinga2 Secret. We are currently supporting following notifiers:

1. [Hipchat](../notifier/hipchat.md#set-environment-variables)
2. [Mailgun](../notifier/mailgun.md#set-environment-variables)
3. [SMTP](../notifier/smtp.md#set-environment-variables)
4. [Twilio](../notifier/twilio.md#set-environment-variables)
5. [Slack](../notifier/slack.md#set-environment-variables)
6. [Plivo](../notifier/plivo.md#set-environment-variables)

If we don't set keys for notifier, notifications will be ignored.

Substitute ENV and deploy secret
```sh
# Deploy Secret
curl https://raw.githubusercontent.com/appscode/searchlight/1.5.9/hack/deploy/icinga2/secret.yaml |
envsubst | kubectl apply -f -
```

###### Create Service
```sh
# Create Service
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/1.5.9/hack/deploy/icinga2/service.yaml
```

###### Create Deployment

To use notifier we need to set some environment variables. See following links to find out how to use different notifiers:

1. [Hipchat](../notifier/hipchat.md#configure)
2. [Mailgun](../notifier/mailgun.md#configure)
3. [SMTP](../notifier/smtp.md#configure)
4. [Twilio](../notifier/twilio.md#configure)
5. [Slack](../notifier/slack.md#configure)
6. [Plivo](../notifier/plivo.md#configure)

```sh
# Create Deployment
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/1.5.9/hack/deploy/icinga2/deployment.yaml
```

### Login

To login into `Icingaweb2`, use following authentication information:
```
Username: admin
Password: <ICINGA_WEB_ADMIN_PASSWORD>
```
Password will be set from Icinga secret.
