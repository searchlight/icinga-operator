> New to Searchlight? Please start [here](/docs/tutorials/README.md).

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
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/cluster-alerts/README.md).
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/node-alerts/README.md).
 - See the list of supported notifiers [here](/docs/tutorials/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/ROADMAP.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/CONTRIBUTING.md).
