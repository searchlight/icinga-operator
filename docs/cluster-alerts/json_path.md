# Check json_path

This is used to check any API response parsing JSON using JQ queries.

#### Supported Kubernetes Objects

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | localhost         |

#### Vars

* `url` - URL to get data
* `query` - JQ query
* `secret` - Kubernetes secret name.
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
```console
hyperalert check_json_path --url='https://api.appscode.com/health' --query='.status' --critical='.status!="OK"'
```
###### Output
```
OK: Response looks good
```

##### Configure Alert Object

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-api-health
  namespace: demo
  labels:
    alert.appscode.com/objectType: cluster
spec:
  check: json_path
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    query: ".status"
    url: https://api.appscode.com/health
    critical: .status!="OK"
```
