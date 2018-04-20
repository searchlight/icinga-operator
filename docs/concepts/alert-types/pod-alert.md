---
title: Pod Alert Overview
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: pod-alert-overview
    name: Pod Alert
    parent: alert-types
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: concepts
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# PodAlerts

## What is PodAlert
A `PodAlert` is a Kubernetes `Custom Resource Definition` (CRD). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for Pods in a Kubernetes native way. You only need to describe the desired check command and notifier in a PodAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.

## PodAlert Spec
As with all other Kubernetes objects, a PodAlert needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example PodAlert object.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: nginx-webstore
  namespace: demo
spec:
  selector:
    matchLabels:
      app: nginx
  check: pod-volume
  vars:
    volumeName: webstore
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

- This Alert is set on pods with matching label `app=nginx` in `demo` namespace.
- Check command `pod-volume` will be applied on volume named `webstore`.
- Icinga will check for volume size every 5m.
- Notifications will be sent every 3m if any problem is detected, until acknowledged.
- When the disk is 70% full, it will reach `Warning` state and emails will be sent to _ops@example.com_ via Mailgun as notification.
- When the disk is 95% full, it will reach `Critical` state and SMSes will be sent to _+1-234-567-8901_ via Twilio as notification.

Any PodAlert object has 3 main sections:

### Pod Selection
Any PodAlert can specify pods in 2 ways:

- `spec.podName` can be used to specify a pod by name. Use this if you are creating pods directly.

- `spec.selector` is a label selector for pods. This should be used if pods are created by workload controllers like Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, etc. Searchlight operator will update Icinga as pods with matching labels are created/deleted by workload controllers.

### Check Command
Check commands are used by Icinga to periodically test some condition. If the test return positive appropriate notifications are sent. The following check commands are supported for pods:
- [pod-exec](/docs/guides/pod-alerts/pod-exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns Critical
- [pod-status](/docs/guides/pod-alerts/pod-status.md) - To check Kubernetes pod status.
- [pod-volume](/docs/guides/pod-alerts/pod-volume.md) - To check Pod volume usage stat.

Each check command has a name specified in `spec.check` field. Optionally each check command can take one or more parameters. These are specified in `spec.vars` field. To learn about the available parameters for each check command, please visit their documentation. `spec.checkInterval` specifies how frequently Icinga will perform this check. Some examples are: 30s, 5m, 6h, etc.

### Notifiers
When a check fails, Icinga will keep sending notifications until acknowledged via IcingaWeb dashboard. `spec.alertInterval` specifies how frequently notifications are sent. Icinga can send notifications to different targets based on alert state. `spec.receivers` contains that list of targets:

| Name                       | Description                                                  |
|----------------------------|--------------------------------------------------------------|
| `spec.receivers[*].state`  | `Required` Name of state for which notification will be sent |
| `spec.receivers[*].to`     | `Required` To whom notifications will be sent                |
| `spec.receivers[*].method` | `Required` How this notification will be sent                |


## Icinga Objects
You can skip this section if you are unfamiliar with how Icinga works. Searchlight operator watches for PodAlert objects and turns them into [Icinga objects](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/) accordingly. For each Kubernetes Pod which has an PodAlert configured, an [Icinga Host](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#host) is created with the name `{namespace}@pod@{pod-name}` and address matching the IP of the Pod. Now for each PodAlert, an [Icinga service](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) is created with name matching the PodAlert name.


## Next Steps
 - Visit the links below to learn about the available check commands for pods:
    - [pod-exec](/docs/guides/pod-alerts/pod-exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns Critical
    - [pod-status](/docs/guides/pod-alerts/pod-status.md) - To check Kubernetes pod status.
    - [pod-volume](/docs/guides/pod-alerts/pod-volume.md) - To check Pod volume stat.
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md).
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
