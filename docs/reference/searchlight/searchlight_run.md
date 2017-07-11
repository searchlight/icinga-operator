## searchlight run

Run operator

### Synopsis


Run operator

```
searchlight run [flags]
```

### Options

```
      --address string      Address to listen on for web interface and telemetry. (default ":56790")
      --analytics           Send analytical event to Google Analytics (default true)
  -s, --config-dir string   Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.
  -h, --help                help for run
      --kubeconfig string   Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string       The address of the Kubernetes API server (overrides any value in kubeconfig)
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
* [searchlight](searchlight.md)	 - Searchlight by AppsCode - Alerts for Kubernetes


