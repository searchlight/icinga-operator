## hyperalert check_pod_exec

Check exit code of exec command on Kubernetes container

### Synopsis


Check exit code of exec command on Kubernetes container

```
hyperalert check_pod_exec [flags]
```

### Options

```
  -a, --argv string        Arguments for exec command. [Format: 'arg; arg; arg']
  -c, --cmd string         Exec command. [Default: /bin/sh] (default "/bin/sh")
  -C, --container string   Container name in specified pod
  -h, --help               help for check_pod_exec
  -H, --host string        Icinga host name
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


