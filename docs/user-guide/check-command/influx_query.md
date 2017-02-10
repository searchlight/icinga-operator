### CheckCommand `influx_query`

This is used to check InfluxDB query result.

#### Supported Kubernetes Objects

| Kubernetes Object      | Icinga2 Host Type  |
| :---:                  | :---:              |
| cluster                | node               |
| nodes                  | node               |
| pods                   | pod                |
| deployments            | pod                |
| daemonsets             | pod                |
| replicasets            | pod                |
| statefulsets           | pod                |
| replicationcontrollers | pod                |
| services               | pod                |

#### Vars

* `influx_host` - URL of InfluxDB host to query
* `secret` - Kubernetes secret name for InfluxDB authentication
* `A` - InfluxDB query (A). Query result will be assigned to variable (A)
* `B` - InfluxDB query (B). Query result will be assigned to variable (B)
* `C` - InfluxDB query (C). Query result will be assigned to variable (C)
* `D` - InfluxDB query (D). Query result will be assigned to variable (D)
* `E` - InfluxDB query (E). Query result will be assigned to variable (E)
* `R` - Equation [A+B] to get result from queries. Result will be assigned to variable (R)
* `warning` - Condition for warning, compare with result. (Example: R > 75)
* `critical` - Condition for critical, compare with result. (Example: R > 90)

> Note: `A`, `B`, `C`, `D`, `E` are parameterized variables.
> Regular expression `pod_name[ ]*=[ ]*'[?]'` is replaced with `pod_name='<pod name>'` for pod.
> And regular expression `nodename[ ]*=[ ]*'[?]'` is replaced with `nodename='<node name>'` for node.

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
