### Notifier `smtp`

This will send a notification email using smtp.

#### Configure

To set `smtp` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: NOTIFY_VIA
  - name: SMTP_HOST
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_HOST
  - name: SMTP_PORT
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_PORT
  - name: SMTP_USERNAME
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_USERNAME
  - name: SMTP_PASSWORD
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_PASSWORD
  - name: SMTP_FROM
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_FROM
  - name: SMTP_TO
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: SMTP_TO
```

##### envconfig for `smtp`

| Name                      | Description                                                                    |
| :---                      | :---                                                                           |
| SMTP_HOST                 | Set host address of smtp server                                                |
| SMTP_PORT                 | Set port of smtp server                                                        |
| SMTP_INSECURE_SKIP_VERIFY | Set `true` to skip ssl verification                                            |
| SMTP_USERNAME             | Set username                                                                   |
| SMTP_PASSWORD             | Set password                                                                   |
| SMTP_FROM                 | Set sender address for notification                                            |
| SMTP_TO                   | Set receipent address. For multiple receipents, set comma separated addresses. |


These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `smtp`

#### Set Environment Variables

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "smtp" | base64  -w 0)
```

##### Key `smtp_host`
Encode and set `SMTP_HOST` to it
```sh
export SMTP_HOST=$(echo <host> | base64  -w 0)
```

##### Key `smtp_port`
Encode and set `SMTP_PORT` to it
```sh
export SMTP_PORT=$(echo <post> | base64  -w 0)
```

##### Key `smtp_username`
Encode and set `SMTP_USERNAME` to it
```sh
export SMTP_USERNAME=$(echo <username> | base64  -w 0)
```

##### Key `smtp_password`
Encode and set `SMTP_PASSWORD` to it
```sh
export SMTP_PASSWORD=$(echo <password> | base64  -w 0)
```

##### Key `smtp_from`
Encode and set `SMTP_FROM` to it
```sh
export SMTP_FROM=$(echo <sender email addresses> | base64  -w 0)
```


##### Key `smtp_to`
Encode and set `SMTP_TO` to it
```sh
export SMTP_TO=$(echo <recipient email addresses> | base64  -w 0)
```
