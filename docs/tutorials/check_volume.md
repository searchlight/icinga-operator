### CheckCommand `volume`

This is used to check Pod volume stat.

#### Supported Kubernetes Objects

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| deployments            | pod               |
| daemonsets             | pod               |
| replicasets            | pod               |
| statefulsets           | pod               |
| replicationcontrollers | pod               |
| services               | pod               |
| pods                   | pod               |

#### Vars

* `volume_name` - Volume name
* `secret_name` - Kubernetes secret name for hostfacts authentication
* `secret_namespace` - Kubernetes namespace of secret
* `warning` - Warning level value (usage percentage defaults to 75.0)
* `critical` - Critical level value (usage percentage defaults to 90.0)

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Hostfacts Secret keys

* `ca.crt`
* `hostfacts.key`
* `hostfacts.crt`
* `auth_token`
* `username`
* `password`

#### Example
###### Command
```console
hyperalert check_volume --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --volume_name=influxdb-persistent-storage --warning=70 --critical=85
# --host are provided by Icinga2
```
###### Output
```
WARNING: Disk used more than 70%
```

#### Required Hostfacts
Before using this CheckCommand, you must need to run `hostfacts` service in each Kubernetes node.
Volume stat of kubernetes pod is collected from `hostfacts` service.

See Hostfacts [deployment guide](hostfacts.md)


##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-pod-volume-1
  namespace: kube-system
spec:
  check: volume
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin
  vars:
    volume_name: influxdb-persistent-storage
    warning: 70
    critical: 85
```
