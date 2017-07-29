# Check json_path

This is used to check any API response parsing JSON using JQ queries.
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

* `url` - URL to get data
* `query` - JQ query
* `secret` - Kubernetes secret name.
* `in_cluster_config` - Use InClusterConfig if hosted in Kubernetes
* `warning` - Warning JQ query which returns [true/false]
* `critical` - Critical JQ query which returns [true/false]

#### Supported Icinga2 State

* OK
* WARNING
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_json_path --url='https://api.appscode.com/health' --query='.status' --critical='.status!="OK"'
```
###### Output
```
OK: Response looks good
```

##### Configure Alert Object

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Alert
metadata:
  name: check-api-health
  namespace: demo
  labels:
    alert.appscode.com/objectType: cluster
spec:
  check: json_path
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    query: ".status"
    url: https://api.appscode.com/health
    critical: .status!="OK"
```
