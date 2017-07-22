### CheckCommand `node_status`

This is used to check Kubernetes Node status.

#### Supported Kubernetes Objects

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
```console
hyperalert check_node_status --host=ip-172-20-0-9.ec2.internal@default
# --host is provided by Icinga2
```
###### Output
```
OK: Node is Ready
```

##### Configure Alert Object
```yaml
# This alert will be set to all nodes individually
apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: check-node-status
  namespace: default
spec:
  check: node_status
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin

# To set alert on specific node, set following labels
# labels:
#   alert.appscode.com/objectType: nodes
#   alert.appscode.com/objectName: ip-172-20-0-9.ec2.internal
```
