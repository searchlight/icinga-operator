# Check pod_exec

This is used to check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL.
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

| Kubernetes Object      | Icinga2 Host Type |
| :---:                  | :---:             |
| deployments            | pod               |
| daemonsets             | pod               |
| replicasets            | pod               |
| statefulsets           | pod               |
| replicationcontrollers | pod               |
| services               | pod               |
| pods                   | pod               |

#### Vars

* `container` - Container name in a Kubernetes Pod
* `cmd` - Exec command. [Default: '/bin/sh']
* `argv` - Exec command arguments. [Format: 'arg; arg; arg']

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_pod_exec --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --argv="ls /var/influxdb/token.ini"
# --host are provided by Icinga2
```
###### Output
```
CRITICAL: Exit Code: 2
```

##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-kube-exec
  namespace: kube-system
  labels:
    alert.appscode.com/objectType: pods
    alert.appscode.com/objectName: monitoring-influxdb-0.12.2-n3lo2
spec:
  check: pod_exec
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    argv: ls /var/influxdb/token.ini
```
