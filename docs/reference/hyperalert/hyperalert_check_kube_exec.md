## hyperalert check_kube_exec

Check exit code of exec command on kubernetes container

### Synopsis


Check exit code of exec command on kubernetes container

```
hyperalert check_kube_exec [flags]
```

### Options

```
  -a, --argv string        Arguments for exec command. [Format: 'arg; arg; arg']
  -c, --cmd string         Exec command. [Default: /bin/sh] (default "/bin/sh")
  -C, --container string   Container name in specified pod
  -h, --help               help for check_kube_exec
  -H, --host string        Icinga host name
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


