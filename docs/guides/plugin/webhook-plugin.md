---
title: Webhook SearchlightPlugin
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-webhook-searchlight-plugin
    name: Webhook SearchlightPlugin
    parent: searchlight-plugin
    weight: 20
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to SearchlightPlugin? Please start [here](/docs/setup/developer-guide/webhook-plugin.md).

# Check Custom Webhook

Searchlight supports adding custom check using SearchlightPlugin CRD. No longer you have to build binary and attach it inside Icinga container.
Simply you can write a HTTP server and tell your check command to call that HTTP server.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: SearchlightPlugin
metadata:
  name: check-pod-count
spec:
  webhook:
    namespace: default
    name: searchlight-plugin
  alertKinds:
  - ClusterAlert
  arguments:
    vars:
      Item:
        warning:
          type: interger
        critical:
          type: interger
  state:
  - OK
  - Critical
  - Unknown
```

Here,

- `metadata.name` will be the name of CheckCommand
- `spec.webhook` states Kubernetes Service information of Webhook
- `spec.alertKinds` determines which kinds of alerts will support this CheckCommand
- `spec.arguments` provides variables information those user will provide to create alert.
- `spec.state` is support Service State to set notifier receivers

```console
$ kubectl apply -f ./docs/examples/plugins/webhook/demo-0.yaml 
searchlightplugin "check-pod-count" created
```

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/add-plugin.svg" width="581" height="362">
</p>

CheckCommand `check-pod-count` is added in Icinga2 configuration. Here, `vars.Item` from `spec.arguments` are added as arguments in CheckCommand.

Few things to be noted here:

- Webhook will be called with URL formatted as bellow:
  
  `http://<spec.webhook.name>.<spec.webhook.namespace>.svc/<metadata.name>`
- Items in `spec.arguments.vars` for example `warning` and `critical` are registered as custom variables. User can provide values for these variables while creating alerts.
- Items in `spec.arguments.host` are added in Icinga CheckCommand arguments.

### Use Icinga Host Variables

You can pass Icinga host variables to your webhook. [Here is the list](https://www.icinga.com/docs/icinga2/latest/doc/03-monitoring-basics/#host-runtime-macros) of available host variables.
Suppose, you need `host.check_attempt` to be forwarded to your webhook, you can add like this

```yaml
spec:
  arguments:
    host:
      attempt: check_attempt
```

Here,

- Icinga host variable `check_attempt` will be forwarded to webhook as variable `attempt`.

> Note: User can't provided value for these variables.


### Create ClusterAlert

Lets create a ClusterAlert for this CheckCommand.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: count-all-pods-demo-0
  namespace: demo
spec:
  check: check-pod-count
  vars:
    warning: 10
    critical: 15
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```

Here,

- `spec.check` is the name of your custom check you added as SearchlightPlugin
- `spec.vars` are variables those are registered when SearchlightPlugin is created with `spec.arguments.vars`

```console
$ kubectl apply -f ./docs/examples/cluster-alerts/count-all-pods/demo-0.yaml
clusteralert "count-all-pods-demo-0" created
```

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/add-alert.svg" width="581" height="362">
</p>

Now periodically, Icinga will call `check_webhook` plugin under `hyperalert`.
And this plugin will call your webhook you have registered in your SearchlightPlugin. According to the response from webhook, Service State will be determined.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/call-webhook.svg" width="581" height="362">
</p>

In the example above, Service State will be **Warning**.

