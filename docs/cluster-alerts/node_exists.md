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
```console
hyperalert check_node_count --count=3
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
  namespace: default
spec:
  check: node_count
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin
  vars:
    count: 3
```
