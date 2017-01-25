### CheckCommand `component_status`

This is used to check Kubernetes components.

#### Supported Kubernetes Objects

| Kubernetes Object   | Icinga2 Host Type  |
| :---:               | :---:              |
| cluster             | localhost          |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```sh
hyperalert check_component_status
```
###### Output
```
OK: All components are healthy
```

##### Configure Alert Object

```yaml
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-component-status
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: component_status
  IcingaParam:
    AlertIntervalSec: 300
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
```
