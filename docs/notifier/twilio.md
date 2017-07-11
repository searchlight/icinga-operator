### Notifier `twilio`

This will send a notification sms using Twilio.

#### Configure

To set `twilio` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: NOTIFY_VIA
  - name: TWILIO_ACCOUNT_SID
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: TWILIO_ACCOUNT_SID
  - name: TWILIO_AUTH_TOKEN
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: TWILIO_AUTH_TOKEN
  - name: TWILIO_FROM
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: TWILIO_FROM
  - name: TWILIO_TO
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: TWILIO_TO
```

##### envconfig for `twilio`

| Name                | Description                                                                        |
| :---                | :---                                                                               |
| TWILIO_ACCOUNT_SID  | Set twilio account SID                                                             |
| TWILIO_AUTH_TOKEN   | Set twilio authentication token                                                    |
| TWILIO_FROM         | Set sender mobile number for notification                                          |
| TWILIO_TO           | Set receipent mobile number. For multiple receipents, set comma separated numbers. |



These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `twilio`

#### Set Environment Variables

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "twilio" | base64  -w 0)
```

##### Key `twilio_account_sid`
Encode and set `TWILIO_ACCOUNT_SID` to it
```sh
export TWILIO_ACCOUNT_SID=$(echo <account sid> | base64  -w 0)
```

##### Key `twilio_auth_token`
Encode and set `TWILIO_AUTH_TOKEN` to it
```sh
export TWILIO_AUTH_TOKEN=$(echo <authentication token> | base64  -w 0)
```

##### Key `twilio_from`
Encode and set `TWILIO_FROM` to it
```sh
export TWILIO_FROM=$(echo <sender mobile number> | base64  -w 0)
```

##### Key `twilio_to`
Encode and set `TWILIO_TO` to it
```sh
export TWILIO_TO=$(echo <receipent mobile numbers> | base64  -w 0)
```
