# Check pod_exec

This is used to check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL.

#### Supported Kubernetes Objects

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| deployments            | pod               |
| daemonsets             | pod               |
| replicasets            | pod               |
| statefulsets           | pod               |
| replicationcontrollers | pod               |
| services               | pod               |
| pods                   | pod               |

#### Vars

* `container` - Container name in a Kubernetes Pod
* `cmd` - Exec command. [Default: '/bin/sh']
* `argv` - Exec command arguments. [Format: 'arg; arg; arg']

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_pod_exec --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --argv="ls /var/influxdb/token.ini"
# --host are provided by Icinga2
```
###### Output
```
CRITICAL: Exit Code: 2
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-kube-exec
  namespace: kube-system
  labels:
    alert.appscode.com/objectType: pods
    alert.appscode.com/objectName: monitoring-influxdb-0.12.2-n3lo2
spec:
  check: pod_exec
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    argv: ls /var/influxdb/token.ini
```
