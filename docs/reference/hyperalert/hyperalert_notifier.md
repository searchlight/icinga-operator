---
title: Notifier
menu:
  product_searchlight_5.1.0:
    identifier: hyperalert-notifier
    name: Notifier
    parent: hyperalert-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_5.1.0
---
## hyperalert notifier

AppsCode Icinga2 Notifier

### Synopsis

AppsCode Icinga2 Notifier

```
hyperalert notifier [flags]
```

### Options

```
  -A, --alert string        Kubernetes alert object name
  -a, --author string       Event author name
  -c, --comment string      Event comment
  -h, --help                help for notifier
  -H, --host string         Icinga host name
      --kubeconfig string   Path to kubeconfig file with authorization information (the master location is set by the master flag).
      --master string       The address of the Kubernetes API server (overrides any value in kubeconfig)
      --output string       Service output
      --state string        Service state (OK | WARNING | CRITICAL)
      --time string         Event time
      --type string         Notification type (PROBLEM | ACKNOWLEDGEMENT | RECOVERY)
```

### Options inherited from parent commands

```
      --allow_verification_with_non_compliant_keys   Allow a SignatureVerifier to use keys which are technically non-compliant with RFC6962.
      --alsologtostderr                              log to standard error as well as files
      --analytics                                    Send analytical events to Google Analytics (default true)
      --log_backtrace_at traceLocation               when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                               If non-empty, write log files in this directory
      --logtostderr                                  log to standard error instead of files
      --stderrthreshold severity                     logs at or above this threshold go to stderr
  -v, --v Level                                      log level for V logs
      --vmodule moduleSpec                           comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [hyperalert](/docs/reference/hyperalert/hyperalert.md)	 - AppsCode Icinga2 plugin


