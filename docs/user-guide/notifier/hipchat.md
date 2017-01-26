### Notifier `hipchat`

This will send a notification to hipchat room.

#### Configure

To set `hipchat` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: appscode-notifier
        key: notify_via
  - name: HIPCHAT_AUTH_TOKEN
    valueFrom:
      secretKeyRef:
        name: appscode-notifier
        key: auth_token
  - name: HIPCHAT_TO
    valueFrom:
      secretKeyRef:
        name: appscode-notifier
        key: to
```

##### envconfig for `hipchat`

| Name                | Description                                                       |
| :---                | :---                                                              |
| HIPCHAT_AUTH_TOKEN  | Set hipchat authentication token                                  |
| HIPCHAT_TO          | Set hipchat room ID. For multiple rooms, set comma separated IDs. |


These environment variables will be set using `appscode-notifier` Secret.

> Set `NOTIFY_VIA` to `hipchat`

#### Create Secret

Create secret with following keys to use hipchat

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "hipchat" | base64  -w 0)
```

##### Key `auth_token`
Encode and set `HIPCHAT_AUTH_TOKEN` to it
```sh
export HIPCHAT_AUTH_TOKEN=$(echo <toke> | base64  -w 0)
```

##### Key `to`
Encode and set `HIPCHAT_TO` to it
```sh
export HIPCHAT_TO=$(echo <hipchat room id> | base64  -w 0)
```

Substitute ENV and deploy secret
```sh
# Deploy Secret
curl https://raw.githubusercontent.com/appscode/searchlight/master/hack/deploy/notifier/hipchat.yaml |
envsubst | kubectl apply -f -
```
