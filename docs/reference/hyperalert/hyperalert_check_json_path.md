## hyperalert check_json_path

Check Json Object

### Synopsis


Check Json Object

```
hyperalert check_json_path [flags]
```

### Options

```
  -c, --critical string     Critical JQ query which returns [true/false]
  -h, --help                help for check_json_path
  -H, --host string         Icinga host name
      --inClusterConfig     Use Kubernetes InCluserConfig
  -s, --secretName string   Kubernetes secret name
  -u, --url string          URL to get data
  -w, --warning string      Warning JQ query which returns [true/false]
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


