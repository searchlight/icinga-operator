### CheckCommand `pod_status`

This is used to check Kubernetes pod status.

#### Supported Kubernetes Objects

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| cluster                | localhost         |
| deployments            | localhost         |
| daemonsets             | localhost         |
| replicasets            | localhost         |
| statefulsets           | localhost         |
| replicationcontrollers | localhost         |
| services               | localhost         |
| pods                   | pod               |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```sh
hyperalert check_pod_status --host='pod_status@default'
# --host is provided by Icinga2
```
###### Output
```json
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
```yaml
# This will check all pod status in default namespace
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-pod-status
  namespace: default
spec:
  check: pod_status
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin


# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
