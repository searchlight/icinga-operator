# Check node_exists

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
```console
hyperalert check_node_exists --count=3
```
###### Output
```
CRITICAL: Found 2 node(s) instead of 3
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: check-node-count
  namespace: demo
spec:
  check: node_exists
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    count: 3
```
