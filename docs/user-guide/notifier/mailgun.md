### Notifier `mailgun`

This will send a notification email via mailgun.

#### Configure

To set `mailgun` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    value: mailgun
  - name: MAILGUN_DOMAIN
    value: <domain>
  - name: MAILGUN_API_KEY
    value: <api key>
  - name: MAILGUN_FROM
    value: <email address>
  - name: MAILGUN_TO
    value: <email address>
```

> Set `NOTIFY_VIA` to `mailgun`

##### envconfig for `mailgun`

| Name                    | Description                                                                    |
| :---                    | :---                                                                           |
| MAILGUN_DOMAIN          | Set domain name for mailgun configuration                                      |
| MAILGUN_API_KEY         | Set mailgun API Key                                                            |
| MAILGUN_PUBLIC_API_KEY  | Set mailgun public API Key                                                     |
| MAILGUN_FROM            | Set sender address for notification                                            |
| MAILGUN_TO              | Set receipent address. For multiple receipents, set comma separated addresses. |
