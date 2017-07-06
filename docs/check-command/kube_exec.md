### CheckCommand `kube_exec`

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
```sh
hyperalert check_kube_exec --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --argv="ls /var/influxdb/token.ini"
# --host are provided by Icinga2
```
###### Output
```
CRITICAL: Exit Code: 2
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-kube-exec
  namespace: kube-system
  labels:
    alert.appscode.com/objectType: pods
    alert.appscode.com/objectName: monitoring-influxdb-0.12.2-n3lo2
spec:
  CheckCommand: kube_exec
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
  vars:
    argv: ls /var/influxdb/token.ini
```
