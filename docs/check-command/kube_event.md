### CheckCommand `kube_event`

This is used to check Kubernetes events. This plugin checks for all Warning events happened in last `c` seconds. Icinga check_interval is provided as `c`.

#### Supported Kubernetes Objects

| Kubernetes Object | Icinga2 Host Type |
| :---:             | :---:             |
| cluster           | localhost         |

#### Vars

* `clock_skew` - Clock skew in Duration. [Default: 30s]. This time is added with check_interval while checking events

#### Supported Icinga2 State

* OK
* WARNING
* UNKNOWN

#### Example
###### Command
```sh
hyperalert check_kube_event --check_interval=1m
# --check_interval are provided by Icinga2
```
###### Output
```json
WARNING: {
   "objects":[  
      {  
         "name":"tc-1916705895-ukpfx",
         "namespace":"default",
         "kind":"Pod",
         "count":5984,
         "reason":"FailedSync",
         "message":"Error syncing pod, skipping: failed to \"StartContainer\" for \"tc\" with ImagePullBackOff: \"Back-off pulling image \\\"appscode/tillerc:765a57f\\\"\"\n"
      },
      {  
         "name":"kube-apiserver-ip-172-20-0-9.ec2.internal",
         "namespace":"kube-system",
         "kind":"Pod",
         "count":300167,
         "reason":"FailedValidation",
         "message":"Error validating pod kube-apiserver-ip-172-20-0-9.ec2.internal.kube-system from file, ignoring: metadata.name: Duplicate value: \"kube-apiserver-ip-172-20-0-9.ec2.internal\""
      }
   ],
   "message":"Found 2 Warning event(s)"
}
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-kube-event
  namespace: default
  labels:
    alert.appscode.com/objectType: cluster
spec:
  CheckCommand: kube_event
  IcingaParam:
    AlertIntervalSec: 120
    CheckIntervalSec: 60
  NotifierParams:
  - Method: EMAIL
    State: CRITICAL
    UserUid: system-admin
```
