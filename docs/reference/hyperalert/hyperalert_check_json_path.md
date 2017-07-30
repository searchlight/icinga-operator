## hyperalert check_json_path

Check Json Object

### Synopsis


Check Json Object

```
hyperalert check_json_path [flags]
```

### Options

```
  -c, --critical string     Critical JQ query which returns [true/false]
  -h, --help                help for check_json_path
  -H, --host string         Icinga host name
      --inClusterConfig     Use Kubernetes InCluserConfig
  -s, --secretName string   Kubernetes secret name
  -u, --url string          URL to get data
  -w, --warning string      Warning JQ query which returns [true/false]
```

### SEE ALSO
* [hyperalert](hyperalert.md)	 - AppsCode Icinga2 plugin


