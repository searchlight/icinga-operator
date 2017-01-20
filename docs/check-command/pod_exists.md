### CheckCommand `pod_exists`

This is used to check Kubernetes pod existence.

#### Icinga2 Host Mapping

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| cluster                | localhost         |
| deployments            | localhost         |
| daemonsets             | localhost         |
| replicasets            | localhost         |
| petsets                | localhost         |
| replicationcontrollers | localhost         |
| services               | localhost         |

#### Vars

* `count` - Number of expected Kubernetes Node


#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
hyperalert check_pod_exists --host='pod_exists@default' --count=7
# --host is provided by Icinga2
```
###### Output
```
OK: Found all pods
```

##### Configure Alert Object

```
# This will check if any pod exists in default namespace
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-pod-exist-1
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: pod_exists
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin

# To check with expected pod number, suppose 8, add following in spec.vars
# vars:
#   count: 8

# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
