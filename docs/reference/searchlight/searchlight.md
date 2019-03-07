---
title: Searchlight
menu:
  product_searchlight_8.0.0-rc.0:
    identifier: searchlight
    name: Searchlight
    parent: searchlight-cli
    weight: 0

product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_8.0.0-rc.0
url: /products/searchlight/8.0.0-rc.0/reference/searchlight/
aliases:
  - products/searchlight/8.0.0-rc.0/reference/searchlight/searchlight/

---
## searchlight

Searchlight by AppsCode - Alerts for Kubernetes

### Synopsis

Searchlight by AppsCode - Alerts for Kubernetes

### Options

```
      --alsologtostderr                  log to standard error as well as files
      --bypass-validating-webhook-xray   if true, bypasses validating webhook xray checks
      --enable-analytics                 send usage events to Google Analytics (default true)
  -h, --help                             help for searchlight
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

* [searchlight configure](/docs/reference/searchlight/searchlight_configure.md)	 - Generate icinga configuration
* [searchlight run](/docs/reference/searchlight/searchlight_run.md)	 - Launch Searchlight operator
* [searchlight version](/docs/reference/searchlight/searchlight_version.md)	 - Prints binary version number.


