---
title: Run
menu:
  product_searchlight_4.0.0:
    identifier: hostfacts-run
    name: Run
    parent: hostfacts-cli
product_name: searchlight
section_menu_id: reference
menu_name: product_searchlight_4.0.0
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
* [hostfacts](/docs/reference/hostfacts/hostfacts.md)	 - Hostfacts by AppsCode - Expose node metrics


