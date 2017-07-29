# Check any_http

# Check pod_exists`

This is used to check Kubernetes pod existence.

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

#### Vars

* `count` - Number of expected Kubernetes Node


#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_pod_exists --host='pod_exists@default' --count=7
# --host is provided by Icinga2
```
###### Output
```
OK: Found all pods
```

##### Configure Alert Object
```yaml
# This will check if any pod exists in default namespace
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: check-pod-exist-1
  namespace: demo
spec:
  check: pod_exists
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]

# To check with expected pod number, suppose 8, add following in spec.vars
# vars:
#   count: 8

# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
