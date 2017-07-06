### Notifier `plivo`

This will send a notification sms using plivo.

#### Configure

To set `plivo` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: searchlight-icinga
        key: notify_via
  - name: PLIVO_AUTH_ID
    valueFrom:
      secretKeyRef:
        name: searchlight-icinga
        key: plivo_auth_id
  - name: PLIVO_AUTH_TOKEN
    valueFrom:
      secretKeyRef:
        name: searchlight-icinga
        key: plivo_auth_token
  - name: PLIVO_FROM
    valueFrom:
      secretKeyRef:
        name: searchlight-icinga
        key: plivo_from
  - name: PLIVO_TO
    valueFrom:
      secretKeyRef:
        name: searchlight-icinga
        key: plivo_to
```

##### envconfig for `plivo`

| Name              | Description                                                                        |
| :---              | :---                                                                               |
| PLIVO_AUTH_ID     | Set plivo auth ID                                                                  |
| PLIVO_AUTH_TOKEN  | Set plivo authentication token                                                     |
| PLIVO_FROM        | Set sender mobile number for notification                                          |
| PLIVO_TO          | Set receipent mobile number. For multiple receipents, set comma separated numbers. |



These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `plivo`

#### Set Environment Variables

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "plivo" | base64  -w 0)
```

##### Key `plivo_auth_id`
Encode and set `PLIVO_AUTH_ID` to it
```sh
export PLIVO_AUTH_ID=$(echo <auth id> | base64  -w 0)
```

##### Key `plivo_auth_token`
Encode and set `PLIVO_AUTH_TOKEN` to it
```sh
export PLIVO_AUTH_TOKEN=$(echo <authentication token> | base64  -w 0)
```

##### Key `plivo_from`
Encode and set `PLIVO_FROM` to it
```sh
export PLIVO_FROM=$(echo <sender mobile number> | base64  -w 0)
```

##### Key `plivo_to`
Encode and set `PLIVO_TO` to it
```sh
export PLIVO_TO=$(echo <receipent mobile numbers> | base64  -w 0)
```
