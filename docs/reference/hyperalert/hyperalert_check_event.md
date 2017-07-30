## hyperalert check_event

Check kubernetes events for all namespaces

### Synopsis


Check kubernetes events for all namespaces

```
hyperalert check_event [flags]
```

### Options

```
  -c, --checkInterval duration           Icinga check_interval in duration. [Format: 30s, 5m]
  -s, --clockSkew duration               Add skew with check_interval in duration. [Default: 30s] (default 30s)
  -h, --help                             help for check_event
  -H, --host string                      Icinga host name
      --involvedObjectKind string        Involved object kind used to select events
      --involvedObjectName string        Involved object name used to select events
      --involvedObjectNamespace string   Involved object namespace used to select events
      --involvedObjectUID string         Involved object uid used to select events
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


