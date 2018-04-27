---
title: Cluster Alert Overview
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: cluster-alert-overview
    name: Cluster Alert
    parent: alert-types
    weight: 5
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: concepts
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).
# ClusterAlerts

## What is ClusterAlert
A `ClusterAlert` is a Kubernetes `Custom Resource Definition` (CRD). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for cluster level alerts in a Kubernetes native way. You only need to describe the desired check command and notifier in a ClusterAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.

## ClusterAlert Spec
As with all other Kubernetes objects, a ClusterAlert needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example ClusterAlert object.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: pod-exists-demo-0
  namespace: demo
spec:
  check: pod-exists
  vars:
    selector: app=nginx
    count: '2'
  checkInterval: 60s
  alertInterval: 3m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Twilio
    state: Critical
    to: ["+1-234-567-8901"]
```

This object will do the followings:

- This Alert is set at cluster level in `demo` namespace.
- Check command `pod-exists` will check for 2 pods matching the label `app=nginx` in `demo` namespace.
- Icinga will check for the existence of pods every 60s.
- Notifications will be sent every 3m if any problem is detected, until acknowledged.
- When the number of pods with label app=nginx is not 2, it will reach `Critical` state and SMSes will be sent to _+1-234-567-8901_ via Twilio as notification.


Any ClusterAlert object has 2 main sections:

### Check Command
Check commands are used by Icinga to periodically test some condition. If the test return positive appropriate notifications are sent. The following check commands are supported for pods:

- [ca-cert](/docs/guides/cluster-alerts/ca-cert.md) - To check expiration of CA certificate used by Kubernetes api server.
- [component-status](/docs/guides/cluster-alerts/component-status.md) - To check Kubernetes component status.
- [event](/docs/guides/cluster-alerts/event.md) - To check Kubernetes Warning events.
- [json-path](/docs/guides/cluster-alerts/json-path.md) - To check any JSON HTTP response using [jq](https://stedolan.github.io/jq/).
- [node-exists](/docs/guides/cluster-alerts/node-exists.md) - To check existence of Kubernetes nodes.
- [pod-exists](/docs/guides/cluster-alerts/pod-exists.md) - To check existence of Kubernetes pods.

Each check command has a name specified in `spec.check` field. Optionally each check command can take one or more parameters. These are specified in `spec.vars` field. To learn about the available parameters for each check command, please visit their documentation. `spec.checkInterval` specifies how frequently Icinga will perform this check. Some examples are: 30s, 5m, 6h, etc.

### Notifiers
When a check fails, Icinga will keep sending notifications until acknowledged via IcingaWeb dashboard. `spec.alertInterval` specifies how frequently notifications are sent. Icinga can send notifications to different targets based on alert state. `spec.receivers` contains that list of targets:

| Name                       | Description                                                  |
|----------------------------|--------------------------------------------------------------|
| `spec.receivers[*].state`  | `Required` Name of state for which notification will be sent |
| `spec.receivers[*].to`     | `Required` To whom notifications will be sent                |
| `spec.receivers[*].method` | `Required` How this notification will be sent                |


## Icinga Objects
You can skip this section if you are unfamiliar with how Icinga works. Searchlight operator watches for ClusterAlert objects and turns them into [Icinga objects](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/) accordingly. A single [Icinga Host](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#host) is created with the name `{namespace}@cluster` and address `127.0.0.1` for all ClusterAlerts in a Kubernetes namespace. Now for each ClusterAlert, an [Icinga service](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) is created with name matching the ClusterAlert name.


## Next Steps
 - Visit the links below to learn about the available check commands for a cluster:
    - [ca-cert](/docs/guides/cluster-alerts/ca-cert.md) - To check expiration of CA certificate used by Kubernetes api server.
    - [component-status](/docs/guides/cluster-alerts/component-status.md) - To check Kubernetes component status.
    - [event](/docs/guides/cluster-alerts/event.md) - To check Kubernetes Warning events.
    - [json-path](/docs/guides/cluster-alerts/json-path.md) - To check any JSON HTTP response using [jq](https://stedolan.github.io/jq/).
    - [node-exists](/docs/guides/cluster-alerts/node-exists.md) - To check existence of Kubernetes nodes.
    - [pod-exists](/docs/guides/cluster-alerts/pod-exists.md) - To check existence of Kubernetes pods.
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - To periodically run various checks on pods in a Kubernetes cluster, use [PodAlerts](/docs/concepts/alert-types/pod-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
