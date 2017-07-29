# Check node_status

This is used to check Kubernetes Node status.
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
| cluster           | node              |
| nodes             | node              |

#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_node_status --host=ip-172-20-0-9.ec2.internal@default
# --host is provided by Icinga2
```
###### Output
```
OK: Node is Ready
```

##### Configure Alert Object
```yaml
# This alert will be set to all nodes individually
apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: check-node-status
  namespace: demo
spec:
  check: node_status
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]

# To set alert on specific node, set following labels
# labels:
#   alert.appscode.com/objectType: nodes
#   alert.appscode.com/objectName: ip-172-20-0-9.ec2.internal
```
