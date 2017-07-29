# Check pod_status

This is used to check Kubernetes pod status.
In this tutorial,

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
  notifierSecretName: any-notifier
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
  6m		6m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "env-demo-0". Reason: secrets "any-notifier" not found
  6m		6m		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "env-demo-0"
```

Voila! `env` command has been synced to Icinga2. Searchlight also logged a warning event, we have not created the notifier secret `any-notifier`. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `env-demo-0`.

![Demo of check_env](/docs/images/cluster-alerts/env/demo-0.gif)

### Cleaning up
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



```console
$ kubectl create ns demo
namespace "demo" created

$ kubectl apply -f ./docs/examples/pod-alerts/pod_status/demo-0.yaml 
replicationcontroller "nginx" created
podalert "pod-status-demo-0" created

$ kubectl get podalert -n demo
NAME                KIND
pod-status-demo-0   PodAlert.v1alpha1.monitoring.appscode.com

$ kubectl describe podalert -n demo pod-status-demo-0
Name:		pod-status-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  21s		21s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-status-demo-0". Reason: secrets "any-notifier" not found
  18s		18s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-0"
  17s		17s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-0"

$ kubectl get pods -n demo
NAME          READY     STATUS    RESTARTS   AGE
nginx-sctrq   1/1       Running   0          28s
nginx-x5rm5   1/1       Running   0          28s
```


```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod_status/demo-1.yaml
pod "busybox" created
podalert "pod-status-demo-1" created

$ kubectl get pods -n demo
NAME      READY     STATUS    RESTARTS   AGE
busybox   1/1       Running   0          4s

$ kubectl get podalert -n demo
NAME                KIND
pod-status-demo-1   PodAlert.v1alpha1.monitoring.appscode.com

$ kubectl describe podalert pod-status-demo-1 -n demo
Name:		pod-status-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  36s		36s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-status-demo-1". Reason: secrets "any-notifier" not found
  36s		36s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-1"
  35s		35s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-1"
```