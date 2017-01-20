### CheckCommand `node_status`

This is used to check Kubernetes Node status.

#### Icinga2 Host Mapping

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | node              |
| nodes             | node              |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
hyperalert check_node_status --host=ip-172-20-0-9.ec2.internal@default
# --host is provided by Icinga2
```
###### Output
```
OK: Node is Ready
```

##### Configure Alert Object

```
# This alert will be set to all nodes individually
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-node-status
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: node_status
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin

# To set alert on specific node, set following labels
# labels:
#   alert.appscode.com/objectType: nodes
#   alert.appscode.com/objectName: ip-172-20-0-9.ec2.internal
```
