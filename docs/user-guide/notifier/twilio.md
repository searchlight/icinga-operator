### Notifier `twilio`

This will send a notification sms using twilio`.

#### Configure

To set `twilio` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    value: twilio
  - name: TWILIO_ACCOUNT_SID
    value: <account SID>
  - name: TWILIO_AUTH_TOKEN
    value: <authentication token>
  - name: TWILIO_FROM
    value: <mobile number>
  - name: TWILIO_TO
    value: <mobile number>
```

> Set `NOTIFY_VIA` to `twilio`

##### envconfig for `twilio`

| Name                | Description                                                                        |
| :---                | :---                                                                               |
| TWILIO_ACCOUNT_SID  | Set twilio account SID                                                             |
| TWILIO_AUTH_TOKEN   | Set twilio authentication token                                                    |
| TWILIO_FROM         | Set sender mobile number for notification                                          |
| TWILIO_TO           | Set receipent mobile number. For multiple receipents, set comma separated numbers. |
