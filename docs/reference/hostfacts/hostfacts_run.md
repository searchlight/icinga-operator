---
title: Run
menu:
  product_searchlight_7.0.0:
    identifier: hostfacts-run
    name: Run
    parent: hostfacts-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_7.0.0
---
## hostfacts run

Run server

### Synopsis

Run server

```
hostfacts run [flags]
```

### Options

```
      --address string      Http server address (default ":56977")
      --caCertFile string   File containing CA certificate
      --certFile string     File container server TLS certificate
  -h, --help                help for run
      --keyFile string      File containing server TLS private key
      --password string     Password used for basic authentication
      --token string        Token used for bearer authentication
      --username string     Username used for basic authentication
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

* [hostfacts](/docs/reference/hostfacts/hostfacts.md)	 - Hostfacts by AppsCode - Expose node metrics


