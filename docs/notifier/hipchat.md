### Notifier `hipchat`

This will send a notification to hipchat room.

#### Configure

To set `hipchat` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: NOTIFY_VIA
  - name: HIPCHAT_AUTH_TOKEN
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: HIPCHAT_AUTH_TOKEN
  - name: HIPCHAT_TO
    valueFrom:
      secretKeyRef:
        name: searchlight-operator
        key: HIPCHAT_TO
```

##### envconfig for `hipchat`

| Name                | Description                                                       |
| :---                | :---                                                              |
| HIPCHAT_AUTH_TOKEN  | Set hipchat authentication token                                  |
| HIPCHAT_TO          | Set hipchat room ID. For multiple rooms, set comma separated IDs. |


These environment variables will be set using `searchlight-icinga` Secret.

> Set `NOTIFY_VIA` to `hipchat`

#### Set Environment Variables

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "hipchat" | base64  -w 0)
```

##### Key `hipchat_auth_token`
Encode and set `HIPCHAT_AUTH_TOKEN` to it
```sh
export HIPCHAT_AUTH_TOKEN=$(echo <toke> | base64  -w 0)
```

##### Key `hipchat_to`
Encode and set `HIPCHAT_TO` to it
```sh
export HIPCHAT_TO=$(echo <hipchat room id> | base64  -w 0)
```
