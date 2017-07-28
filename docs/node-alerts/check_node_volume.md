### CheckCommand `node_volume`

This is used to check Node Disk stat.

#### Supported Kubernetes Objects

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | node              |
| nodes             | node              |

#### Vars

* `secret_name` - Kubernetes secret name for hostfacts authentication
* `secret_namespace` - Kubernetes namespace of secret
* `warning` - Warning level value (usage percentage defaults to 75.0)
* `critical` - Critical level value (usage percentage defaults to 90.0)

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_volume --node_stat --host ip-172-20-0-9.ec2.internal@default
# --node_stat and --host are provided by Icinga2
```
###### Output
```
OK: (Disk & Inodes)
```

#### Required Hostfacts
Before using this CheckCommand, you must need to run `hostfacts` service in each Kubernetes node.
Node disk stat is collected from `hostfacts` service deployed in each node.

See Hostfacts [deployment guide](hostfacts.md)


##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: check-node-disk
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  check: node_volume
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - method: EMAIL
    state: CRITICAL
    to: system-admin

# To set alert on specific node, set following labels
#  selector:
#    kubernetes.io/hostname: ip-172-20-0-9.ec2.internal
```