---
title: Node Alert Overview
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: node-alert-overview
    name: Node Alert
    parent: alert-types
    weight: 10
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: concepts
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# NodeAlerts

## What is NodeAlert
A `NodeAlert` is a Kubernetes `Custom Resource Definition` (CRD). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for Nodes in a Kubernetes native way. You only need to describe the desired check command and notifier in a NodeAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.

## NodeAlert Spec
As with all other Kubernetes objects, a NodeAlert needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example NodeAlert object.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: webstore
  namespace: demo
spec:
  selector:
    beta.kubernetes.io/os: linux
  check: node-volume
  vars:
    warning: '70'
    critical: '95'
  checkInterval: 5m
  alertInterval: 3m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Warning
    to: ["ops@example.com"]
  - notifier: Twilio
    state: Critical
    to: ["+1-234-567-8901"]
```

This object will do the followings:

- This Alert is set on nodes with matching label `beta.kubernetes.io/os=linux`.
- Check command `node-volume` will be used.
- Icinga will check for volume size every 5m.
- Notifications will be sent every 3m if any problem is detected, until acknowledged.
- When the disk is 70% full, it will reach `Warning` state and emails will be sent to _ops@example.com_ via Mailgun as notification.
- When the disk is 95% full, it will reach `Critical` state and SMSes will be sent to _+1-234-567-8901_ via Twilio as notification.

Any NodeAlert object has 3 main sections:

### Node Selection
Any NodeAlert can specify nodes in 2 ways:

- `spec.nodeName` can be used to specify a node by name.

- `spec.selector` is a node selector. Searchlight operator will update Icinga as nodes with matching labels are added/removed.

### Check Command
Check commands are used by Icinga to periodically test some condition. If the test return positive appropriate notifications are sent. The following check commands are supported for nodes:
- [node-status](/docs/guides/node-alerts/node-status.md) - To check Kubernetes Node status.
- [node-volume](/docs/guides/node-alerts/node-volume.md) - To check Node Disk stat.

Each check command has a name specified in `spec.check` field. Optionally each check command can take one or more parameters. These are specified in `spec.vars` field. To learn about the available parameters for each check command, please visit their documentation. `spec.checkInterval` specifies how frequently Icinga will perform this check. Some examples are: 30s, 5m, 6h, etc.

### Notifiers
When a check fails, Icinga will keep sending notifications until acknowledged via IcingaWeb dashboard. `spec.alertInterval` specifies how frequently notifications are sent. Icinga can send notifications to different targets based on alert state. `spec.receivers` contains that list of targets:

| Name                       | Description                                                  |
|----------------------------|--------------------------------------------------------------|
| `spec.receivers[*].state`  | `Required` Name of state for which notification will be sent |
| `spec.receivers[*].to`     | `Required` To whom notifications will be sent                |
| `spec.receivers[*].method` | `Required` How this notification will be sent                |


## Icinga Objects
You can skip this section if you are unfamiliar with how Icinga works. Searchlight operator watches for NodeAlert objects and turns them into [Icinga objects](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/) accordingly. For each Kubernetes Node which has an NodeAlert configured, an [Icinga Host](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#host) is created with the name `{namespace}@node@{node-name}` and address matching the internal IP of the Node. Now for each NodeAlert, an [Icinga service](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) is created with name matching the NodeAlert name.


## Next Steps
 - Visit the links below to learn about the available check commands for nodes:
    - [node-status](/docs/guides/node-alerts/node-status.md) - To check Kubernetes Node status.
    - [node-volume](/docs/guides/node-alerts/node-volume.md) - To check Node Disk stat.
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md).
 - To periodically run various checks on pods in a Kubernetes cluster, use [PodAlerts](/docs/concepts/alert-types/pod-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
