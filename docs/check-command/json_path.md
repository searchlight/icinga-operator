### CheckCommand `json_path`

This is used to check any API response parsing JSON using JQ queries.

#### Icinga2 Host Mapping

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | localhost         |

#### Vars

* `url` - URL to get data
* `query` - JQ query
* `secret` - Kubernetes secret name (secret-name.namespace)
* `in_cluster_config` - Use InClusterConfig if hosted in Kubernetes
* `warning` - Warning JQ query which returns [true/false]
* `critical` - Critical JQ query which returns [true/false]

#### Supported Icinga2 State

* OK
* WARNING
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
hyperalert check_json_path --url='https://api.appscode.com/health' --query='.status' --critical='.status!="OK"'
```
###### Output
```
OK: Response looks good
```

##### Configure Alert Object

```
apiVersion: appscode.com/v1beta1
kind: Alert
metadata:
  name: check-api-health
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: json_path
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
  Vars:
    query: ".status"
    url: https://api.appscode.com/health
    critical: .status!="OK"
```
