# Check event

This is used to check Kubernetes events. This plugin checks for all Warning events happened in last `c` seconds. Icinga check_interval is provided as `c`.
In this tutorial,


## Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

~ $ kubectl get namespaces
NAME          STATUS    AGE
default       Active    6h
kube-public   Active    6h
kube-system   Active    6h
demo          Active    4m
```



## Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps


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
```console
hyperalert check_event --check_interval=1m
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
  namespace: demo
  labels:
    alert.appscode.com/objectType: cluster
spec:
  check: event
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
