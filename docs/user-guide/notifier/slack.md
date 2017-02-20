### Notifier `slack`

This will send a notification to slack channel.

#### Configure

To set `slack` as notifier, we need to set following environment variables in Icinga2 deployment.

```yaml
env:
  - name: NOTIFY_VIA
    valueFrom:
      secretKeyRef:
        name: appscode-icinga
        key: notify_via
  - name: SLACK_AUTH_TOKEN
    valueFrom:
      secretKeyRef:
        name: appscode-icinga
        key: slack_auth_token
  - name: SLACK_CHANNEL
    valueFrom:
      secretKeyRef:
        name: appscode-icinga
        key: slack_channel
```

##### envconfig for `slack`

| Name             | Description                                                               |
| :---             | :---                                                                      |
| SLACK_AUTH_TOKEN | Set slack access authentication token                                     |
| SLACK_CHANNEL    | Set slack channel name. For multiple channels, set comma separated names. |


#### Add Searchlight app
Add Searchlight app in your slack channel and use provided `bot_access_token`.

<a href="https://slack.com/oauth/authorize?scope=bot&client_id=31843174386.143405120770"><img alt="Add to Slack" height="40" width="139" src="https://platform.slack-edge.com/img/add_to_slack.png" srcset="https://platform.slack-edge.com/img/add_to_slack.png 1x, https://platform.slack-edge.com/img/add_to_slack@2x.png 2x" /></a>

#### Set Environment Variables

These environment variables will be set using `appscode-icinga` Secret.

> Set `NOTIFY_VIA` to `slack`

##### Key `notify_via`
Encode and set `NOTIFY_VIA` to it
```sh
export NOTIFY_VIA=$(echo "slack" | base64  -w 0)
```

##### Key `slack_auth_token`
Encode and set `SLACK_AUTH_TOKEN` to it
```sh
export SLACK_AUTH_TOKEN=$(echo <toke> | base64  -w 0)
```

##### Key `slack_channel`
Encode and set `SLACK_CHANNEL` to it
```sh
export SLACK_CHANNEL=$(echo <slack channel name> | base64  -w 0)
```

