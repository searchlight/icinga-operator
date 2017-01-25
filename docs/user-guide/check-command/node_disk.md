### CheckCommand `node_disk`

This is used to check Node Disk stat.

#### Supported Kubernetes Objects

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | node              |
| nodes             | node              |

#### Vars

* `warning` - Warning level value (usage percentage defaults to 75.0)
* `critical` - Critical level value (usage percentage defaults to 90.0)

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```sh
hyperalert check_volume --node_stat --host ip-172-20-0-9.ec2.internal@default
# --node_stat and --host are provided by Icinga2
```
###### Output
```
OK: (Disk & Inodes)
```

##### Configure Alert Object
```yaml
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-node-disk
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: node_disk
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
