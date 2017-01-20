### CheckCommand `volume`

This is used to check Pod volume stat.

#### Icinga2 Host Mapping

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| deployments            | pod               |
| daemonsets             | pod               |
| replicasets            | pod               |
| petsets                | pod               |
| replicationcontrollers | pod               |
| services               | pod               |
| pods                   | pod               |

#### Vars

* `name` - Volume name
* `warning` - Warning level value (usage percentage defaults to 75.0)
* `critical` - Critical level value (usage percentage defaults to 90.0)

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
hyperalert check_volume --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --name=influxdb-persistent-storage --warning=70 --critical=85
# --host are provided by Icinga2
```
###### Output
```
WARNING: Disk used more than 70%
```

##### Configure Alert Object
```
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-pod-volume-1
  namespace: kube-system
  labels:
    alert.appscode.com/objectType: services
    alert.appscode.com/objectName: monitoring-influxdb
spec:
  CheckCommand: volume
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
  vars:
    name: influxdb-persistent-storage
    warning: 70
    critical: 85
```
