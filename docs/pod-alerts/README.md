> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# PodAlerts

## What is PodAlert
A `PodAlert` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for pods in a Kubernetes native way. You only need to describe the desired check command and notifier in a PodAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.

## PodAlert Spec
As with all other Kubernetes objects, a PodAlert needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example PodAlert object.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: pod-volume
  namespace: demo
spec:
  selector:
    matchLabels:
      app: nginx
  check: pod_volume
  vars:
    volumeName: website-storage
    warning: 70
    critical: 95
  checkInterval: 30s
  alertInterval: 5m
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

* This Alert is set on ReplicationController named `elasticsearch-logging-v1` in `kube-system` namespace.
* CheckCommand `volume` will be applied.
* Icinga2 Service will check volume every 60s.
* Notifications will be send every 5m if any problem is detected.
* Email will be sent as a notification to admin user for `CRITICAL` state. For other states, no notification will be sent.
* On each Pod under specified RC, volume named `disk` will be checked. If volume is used more than 60%, it is `WARNING`. For 75%, it is `CRITICAL`.

## Explanation

### Alert Object Fields

* apiVersion - The Kubernetes API version.
* kind - The Kubernetes object type.
* metadata.name - The name of the Alert object.
* metadata.namespace - The namespace of the Alert object
* metadata.labels - The Kubernetes object labels. This labels are used to determine for which object this alert will be set.
* spec.check - Icinga CheckCommand name
* spec.checkInterval - How frequently Icinga Service will be checked
* spec.alertInterval - How frequently notifications will be send
* spec.receivers - NotifierParams contains array of information to send notifications for Incident
* spec.vars - Vars contains array of Icinga Service variables to be used in CheckCommand.

#### NotifierParam Fields

* state - For which state notification will be sent
* to - To whom notification will be sent
* method - How this notification will be sent

> `NotifierParams` is only used when notification is sent via `AppsCode`.

## Check Commands
The following check command are supported for pods:
- [influx_query](influx_query.md) - To check InfluxDB query result.
- [pod_exec](pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
- [pod_status](pod_status.md) - To check Kubernetes pod status.
- [pod_volume](pod_volume.md) - To check Pod volume stat.
