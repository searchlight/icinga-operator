> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Check json_path

This is used to check any API response parsing JSON using JQ queries.

ClusterAlert `env` prints the list of environment variables in searchlight-operator pods. This check command is used to test Searchlight.


## Spec
`env` check command has no variables. Execution of this command can result in following states:
- OK
- WARNING
- CRITICAL
- UNKNOWN


## Tutorial

### Before You Begin
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

### Create Alert
In this tutorial, we are going to create an alert to check `env`.
```yaml
$ cat ./docs/examples/cluster-alerts/env/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: env-demo-0
  namespace: demo
spec:
  check: env
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/env/demo-0.yaml 
clusteralert "env-demo-0" created

$ kubectl describe clusteralert env-demo-0 -n demo
Name:		env-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  6m		6m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "env-demo-0". Reason: secrets "notifier-config" not found
  6m		6m		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "env-demo-0"
```

Voila! `env` command has been synced to Icinga2. Searchlight also logged a warning event, we have not created the notifier secret `notifier-config`. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `env-demo-0`.

![Demo of check_env](/docs/images/cluster-alerts/env/demo-0.gif)

### Cleaning up
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
* `inClusterConfig` - Use InClusterConfig if hosted in Kubernetes
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
