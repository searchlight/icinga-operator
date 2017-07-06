### CheckCommand `node_count`

This is used to check total number of Kubernetes node.

#### Supported Kubernetes Objects

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | localhost         |

#### Vars

* `count` - Number of expected Kubernetes Node

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```sh
hyperalert check_node_count --count=3
```
###### Output
```
CRITICAL: Found 2 node(s) instead of 3
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-node-count
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: node_count
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
  Vars:
    count: 3
```
