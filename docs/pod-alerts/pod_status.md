# Check pod_status

This is used to check Kubernetes pod status.
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
| cluster                | localhost         |
| deployments            | localhost         |
| daemonsets             | localhost         |
| replicasets            | localhost         |
| statefulsets           | localhost         |
| replicationcontrollers | localhost         |
| services               | localhost         |
| pods                   | pod               |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_pod_status --host='pod_status@default'
# --host is provided by Icinga2
```
###### Output
```json
CRITICAL: {
  "objects": [
    {
      "name": "test-pod-0",
      "namespace": "default",
      "status": "Pending"
    },
    {
      "name": "test-pod-1",
      "namespace": "default",
      "status": "Pending"
    }
  ],
  "message": "Found 2 not running pods(s)"
}
```

##### Configure Alert Object
```yaml
# This will check all pod status in default namespace
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-pod-status
  namespace: demo
spec:
  check: pod_status
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]


# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
