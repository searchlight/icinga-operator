### Notifier `hipchat`

This will send a notification to hipchat room.

#### Configure

To set `hipchat` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    value: hipchat
  - name: HIPCHAT_AUTH_TOKEN
    value: <token>
  - name: HIPCHAT_TO
    value: <id>
```

> Set `NOTIFY_VIA` to `hipchat`

##### envconfig for `hipchat`

| Name                | Description                                                       |
| :---                | :---                                                              |
| HIPCHAT_AUTH_TOKEN  | Set hipchat authentication token                                  |
| HIPCHAT_TO          | Set hipchat room ID. For multiple rooms, set comma separated IDs. |
