## hyperalert check_volume

Check kubernetes volume

### Synopsis


Check kubernetes volume

```
hyperalert check_volume [flags]
```

### Options

```
  -c, --critical float       Critical level value (usage percentage) (default 90)
  -h, --help                 help for check_volume
  -H, --host string          Icinga host name
      --node_stat            Checking Node disk size
  -s, --secret string        Kubernetes secret name
  -N, --volume_name string   Volume name
  -w, --warning float        Warning level value (usage percentage) (default 75)
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


