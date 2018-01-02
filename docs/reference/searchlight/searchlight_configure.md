---
title: Configure
menu:
  product_searchlight_5.0.0:
    identifier: searchlight-configure
    name: Configure
    parent: searchlight-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_5.0.0
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
      --allow_verification_with_non_compliant_keys   Allow a SignatureVerifier to use keys which are technically non-compliant with RFC6962.
      --alsologtostderr                              log to standard error as well as files
      --analytics                                    Send analytical events to Google Analytics (default true)
      --log_backtrace_at traceLocation               when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                               If non-empty, write log files in this directory
      --logtostderr                                  log to standard error instead of files (default true)
      --stderrthreshold severity                     logs at or above this threshold go to stderr (default 2)
  -v, --v Level                                      log level for V logs
      --vmodule moduleSpec                           comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [searchlight](/docs/reference/searchlight/searchlight.md)	 - Searchlight by AppsCode - Alerts for Kubernetes


