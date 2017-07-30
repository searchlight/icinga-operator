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

### Options inherited from parent commands

```
      --allow_verification_with_non_compliant_keys   Allow a SignatureVerifier to use keys which are technically non-compliant with RFC6962.
      --alsologtostderr                              log to standard error as well as files
      --log_backtrace_at traceLocation               when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                               If non-empty, write log files in this directory
      --logtostderr                                  log to standard error instead of files (default true)
      --stderrthreshold severity                     logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                      log level for V logs
      --vmodule moduleSpec                           comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


