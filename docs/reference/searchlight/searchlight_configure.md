---
title: Configure
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: searchlight-configure
    name: Configure
    parent: searchlight-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_6.0.0-rc.0
---
## searchlight configure

Generate icinga configuration

### Synopsis

Generate icinga configuration

```
searchlight configure [flags]
```

### Options

```
  -s, --config-dir string   Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.
  -h, --help                help for configure
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --enable-analytics                 send usage events to Google Analytics (default true)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [searchlight](/docs/reference/searchlight/searchlight.md)	 - Searchlight by AppsCode - Alerts for Kubernetes


