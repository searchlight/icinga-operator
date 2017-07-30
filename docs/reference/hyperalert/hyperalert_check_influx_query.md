## hyperalert check_influx_query

Check InfluxDB Query Data

### Synopsis


Check InfluxDB Query Data

```
hyperalert check_influx_query [flags]
```

### Options

```
      --A string            InfluxDB query A
      --B string            InfluxDB query B
      --C string            InfluxDB query C
      --D string            InfluxDB query D
      --E string            InfluxDB query E
      --R string            Equation to evaluate result
  -c, --critical string     Critical query which returns [true/false]
  -h, --help                help for check_influx_query
  -H, --host string         Icinga host name
      --influxHost string   URL of InfluxDB host to query
  -s, --secretName string   Kubernetes secret name
  -w, --warning string      Warning query which returns [true/false]
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


