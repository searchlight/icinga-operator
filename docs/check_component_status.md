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
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-component-status
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  check: component_status
  IcingaParam:
    alertInterval: 300
    checkInterval: 60
  receivers:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
```
