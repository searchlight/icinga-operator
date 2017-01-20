### CheckCommand `pod_status`

This is used to check Kubernetes pod status.

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
| pods                   | pod               |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
hyperalert check_pod_status --host='pod_status@default'
# --host is provided by Icinga2
```
###### Output
```
CRITICAL: {
  "objects": [
    {
      "name": "test-pod-0",
      "namespace": "default",
      "status": "Pending"
    },
    {
      "name": "test-pod-1",
      "namespace": "default",
      "status": "Pending"
    }
  ],
  "message": "Found 2 not running pods(s)"
}
```

##### Configure Alert Object

```
# This will check all pod status in default namespace
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-pod-status
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: pod_status
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin


# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
