### CheckCommand `prometheus_metric`

This is used to check Prometheus query result.

#### Supported Kubernetes Objects

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| cluster                | localhost         |
| nodes                  | node              |
| deployments            | pod               |
| daemonsets             | pod               |
| replicasets            | pod               |
| statefulsets           | pod               |
| replicationcontrollers | pod               |
| services               | pod               |
| pods                   | pod               |

#### Vars

* `prom_host` - URL of Prometheus host to query
* `query` - Prometheus query that returns a float or int
* `metric_name` - Name for the metric being checked
* `method` - Comparison method, one of gt, ge, lt, le, eq, ne. (Defauls to "ge")
* `accept_nan` - Accept NaN as an "OK" result
* `warning` - Warning level value (must be zero or positive)
* `critical` - Critical level value (must be zero or positive)

#### Supported Icinga2 State

* OK
* WARNING
* CRITICAL
* UNKNOWN

#### Example
###### Command
```
```
###### Output
```
```

##### Configure Alert Object
```
```
