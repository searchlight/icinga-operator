## hyperalert check_prometheus_metric

Check prometheus metric

### Synopsis


Check prometheus metric

```
hyperalert check_prometheus_metric [flags]
```

### Options

```
  -O, --accept_nan           Accept NaN as an "OK" result
  -c, --critical int         Critical level value (must be zero or positive)
  -h, --help                 help for check_prometheus_metric
  -m, --method string        Comparison method, one of gt, ge, lt, le, eq, ne
	(defaults to ge unless otherwise specified) (default "ge")
  -n, --metric_name string   A name for the metric being checked
  -H, --prom_host string     URL of Prometheus host to query
  -q, --query string         Prometheus query that returns a float or int
  -w, --warning int          Warning level value (must be zero or positive)
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


