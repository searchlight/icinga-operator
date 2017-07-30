> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# NodeAlerts

## What is NodeAlert
A `NodeAlert` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for Nodes in a Kubernetes native way. You only need to describe the desired check command and notifier in a NodeAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.

## NodeAlert Spec
As with all other Kubernetes objects, a NodeAlert needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example NodeAlert object.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: nginx-webstore
  namespace: demo
spec:
  selector:
    matchLabels:
      app: nginx
  check: pod_volume
  vars:
    volumeName: webstore
    warning: 70
    critical: 95
  checkInterval: 5m
  alertInterval: 3m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: WARNING
    to: ["ops@example.com"]
  - notifier: twilio
    state: CRITICAL
    to: ["+1-234-567-8901"]
```

This object will do the followings:

- This Alert is set on pods with matching label `app=nginx` in `demo` namespace.
- Check command `pod_volume` will be applied on volume named `webstore`.
- Icinga will check for volume size every 60s.
- Notifications will be sent every 5m if any problem is detected, until acknowledged.
- When the disk is 70% full, it will reach `WARNING` state and emails will be sent to _ops@example.com_ via Mailgun as notification.
- When the disk is 95% full, it will reach `CRITICAL` state and SMSes will be sent to _+1-234-567-8901_ via Twilio as notification.

Any NodeAlert object has 3 main sections:

### Pod Selection
Any NodeAlert can specify pods in 2 ways:

- `spec.podName` can be used to specify a pod by name. Use this if you are creating pods directly.

- `spec.selector` is a label selector for pods. This should be used if pods are created by workload controllers like Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, etc. Searchlight operator will update Icinga as pods with matching labels are created/deleted by workload controllers.

### Check Command
Check commands are used by Icinga to periodically test some condition. If the test return positive appropriate notifications are sent. The following check commands are supported for pods:
- [influx_query](influx_query.md) - To check InfluxDB query result.
- [pod_exec](pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
- [pod_status](pod_status.md) - To check Kubernetes pod status.
- [pod_volume](pod_volume.md) - To check Pod volume stat.

Each check command has a name specified in `spec.check` field. Optionally each check command can take one or more parameters. These are specified in `spec.vars` field. To learn about the available parameters for each check command, please visit their documentation. `spec.checkInterval` specifies how frequently Icinga will perform this check. Some examples are: 30s, 5m, 6h, etc.

### Notifiers
When a check fails, Icinga will keep sending notifications until acknowledged via IcingaWeb dashboard. `spec.alertInterval` specifies how frequently notifications are sent. Icinga can send notifications to different targets based on alert state. `spec.receivers` contains that list of targets:

| Name                       | Description                                                  |
|----------------------------|--------------------------------------------------------------|
| `spec.receivers[*].state`  | `Required` Name of state for which notification will be sent |
| `spec.receivers[*].to`     | `Required` To whom notifications will be sent                |
| `spec.receivers[*].method` | `Required` How this notification will be sent                |












* [component_status](check_component_status.md) - To check Kubernetes components.
* [influx_query](check_influx_query.md) - To check InfluxDB query result.
* [json_path](check_json_path.md) - To check any API response by parsing JSON using JQ queries.
* [node_count](check_node_count.md) - To check total number of Kubernetes node.
* [node_status](check_node_status.md) - To check Kubernetes Node status.
* [pod_exists](check_pod_exists.md) - To check Kubernetes pod existence.
* [pod_status](check_pod_status.md) - To check Kubernetes pod status.
* [prometheus_metric](check_prometheus_metric.md) - To check Prometheus query result.
* [node_volume](check_node_volume.md) - To check Node Disk stat.
* [volume](check_pod_volume.md) - To check Pod volume stat.
* [event](check_event.md) - To check Kubernetes events for all Warning TYPE happened in last 'c' seconds.
* [pod_exec](check_pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
