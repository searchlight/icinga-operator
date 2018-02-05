---
title: Check Component Status
menu:
  product_searchlight_6.0.0-alpha.0:
    identifier: hyperalert-check-component-status
    name: Check Component Status
    parent: hyperalert-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_6.0.0-alpha.0
---
## hyperalert check_component_status

Check Kubernetes Component Status

### Synopsis

Check Kubernetes Component Status

```
hyperalert check_component_status [flags]
```

### Options

```
  -n, --componentName string   Name of component which should be ready
  -h, --help                   help for check_component_status
      --kubeconfig string      Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string          The address of the Kubernetes API server (overrides any value in kubeconfig)
  -l, --selector string        Selector (label query) to filter on, supports '=', '==', and '!='.
```

### Options inherited from parent commands

```
      --allow_verification_with_non_compliant_keys   Allow a SignatureVerifier to use keys which are technically non-compliant with RFC6962.
      --alsologtostderr                              log to standard error as well as files
      --analytics                                    Send analytical events to Google Analytics (default true)
      --log_backtrace_at traceLocation               when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                               If non-empty, write log files in this directory
      --logtostderr                                  log to standard error instead of files (default true)
      --stderrthreshold severity                     logs at or above this threshold go to stderr
  -v, --v Level                                      log level for V logs
      --vmodule moduleSpec                           comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [hyperalert](/docs/reference/hyperalert/hyperalert.md)	 - AppsCode Icinga2 plugin


