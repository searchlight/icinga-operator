---
title: Version
menu:
  product_searchlight_7.0.0:
    identifier: searchlight-version
    name: Version
    parent: searchlight-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_7.0.0
---
## searchlight version

Prints binary version number.

### Synopsis

Prints binary version number.

```
searchlight version [flags]
```

### Options

```
  -h, --help    help for version
      --short   Print just the version number.
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --bypass-validating-webhook-xray   if true, bypasses validating webhook xray checks
      --enable-analytics                 send usage events to Google Analytics (default true)
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
      --use-kubeapiserver-fqdn-for-aks   if true, uses kube-apiserver FQDN for AKS cluster to workaround https://github.com/Azure/AKS/issues/522 (default true)
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [searchlight](/docs/reference/searchlight/searchlight.md)	 - Searchlight by AppsCode - Alerts for Kubernetes

