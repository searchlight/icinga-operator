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
  alertInterval: 5m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin
```
