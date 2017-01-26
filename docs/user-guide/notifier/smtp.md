### Notifier `smtp`

This will send a notification email using smtp.

#### Configure

To set `smtp` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    value: smtp
  - name: SMTP_HOST
    value: <SMTP server address>
  - name: SMTP_PORT
    value: <SMTP server port>
  - name: SMTP_USERNAME
    value: <username>
  - name: SMTP_PASSWORD
    value: <password>
  - name: SMTP_FROM
    value: <email address>
  - name: SMTP_TO
    value: <email address>
```

> Set `NOTIFY_VIA` to `smtp`

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
