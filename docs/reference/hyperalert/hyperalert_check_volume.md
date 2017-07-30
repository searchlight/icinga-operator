## hyperalert check_volume

Check kubernetes volume

### Synopsis


Check kubernetes volume

```
hyperalert check_volume [flags]
```

### Options

```
  -c, --critical float      Critical level value (usage percentage) (default 95)
  -h, --help                help for check_volume
  -H, --host string         Icinga host name
      --nodeStat            Checking Node disk size
  -s, --secretName string   Kubernetes secret name
  -N, --volumeName string   Volume name
  -w, --warning float       Warning level value (usage percentage) (default 80)
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


