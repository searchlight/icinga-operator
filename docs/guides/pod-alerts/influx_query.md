---
title: Pod Influx Query
menu:
  product_searchlight_5.0.0:
    identifier: pod-influx-query
    name: Influx Query
    parent: pod-alert
    weight: 20
product_name: searchlight
menu_name: product_searchlight_5.0.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check influx_query

Check command `influx_query` is used to check InfluxDB query result.


## Spec
`influx_query` check command has the following variables.

- `influxHost` - URL of InfluxDB host to query
- `secretName` - Name of Secret used for InfluxDB authentication
- `A` - InfluxDB query (A). Query result will be assigned to variable (A)
- `B` - InfluxDB query (B). Query result will be assigned to variable (B)
- `C` - InfluxDB query (C). Query result will be assigned to variable (C)
- `D` - InfluxDB query (D). Query result will be assigned to variable (D)
- `E` - InfluxDB query (E). Query result will be assigned to variable (E)
- `R` - Equation [A+B] to get result from queries. Result will be assigned to variable (R)
- `warning` - Condition for warning, compare with result. (Example: R > 75)
- `critical` - Condition for critical, compare with result. (Example: R > 90)

Here `A`, `B`, `C`, `D`, `E` are processed as a GO template generating the final InfluxDB query. The available template variables are: PodName, PodIP and Namespace.

Execution of this command can result in following states:

- OK
- WARNING
- CRITICAL
- UNKNOWN


## Next Steps
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md).
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
