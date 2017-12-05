---
title: Check Pod Exec
menu:
  product_searchlight_4.0.0:
    identifier: hyperalert-check-pod-exec
    name: Check Pod Exec
    parent: hyperalert-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_4.0.0
---
## hyperalert check_pod_exec

Check exit code of exec command on Kubernetes container

### Synopsis

Check exit code of exec command on Kubernetes container

```
hyperalert check_pod_exec [flags]
```

### Options

```
  -a, --argv string         Arguments for exec command. [Format: 'arg; arg; arg']
  -c, --cmd string          Exec command. [Default: /bin/sh] (default "/bin/sh")
  -C, --container string    Container name in specified pod
  -h, --help                help for check_pod_exec
  -H, --host string         Icinga host name
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

* [hyperalert](/docs/reference/hyperalert/hyperalert.md)	 - AppsCode Icinga2 plugin


