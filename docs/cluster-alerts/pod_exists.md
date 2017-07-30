> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Check pod_exists

Check command `pod_exists` is used to check existence of pods in a Kubernetes cluster.


## Spec
`pod_exists` has the following variables:
- `selector` - Label selector for pods whose existence are checked.
- `podName` - Name of Kubernetes pod whose existence is checked.
- `count` - Number of expected Kubernetes pods

Execution of this command can result in following states:
- OK
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


### Check existence of pods with matching labels
In this tutorial, a ClusterAlert will be used check existence of pods with matching labels by setting `spec.vars.selector` field.
```yaml
$ cat ./docs/examples/cluster-alerts/pod_exists/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: pod-exists-demo-0
  namespace: demo
spec:
  check: pod_exists
  vars:
    selector: app=nginx
    count: 2
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/pod_exists/demo-0.yaml
replicationcontroller "nginx" created
clusteralert "pod-exists-demo-0" created

$ kubectl get clusteralert -n demo
NAME                KIND
pod-exists-demo-0   ClusterAlert.v1alpha1.monitoring.appscode.com
~/g/s/g/a/searchlight (d7) $ kubectl describe clusteralert -n demo pod-exists-demo-0
Name:		pod-exists-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  19s		19s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "pod-exists-demo-0". Reason: secrets "notifier-config" not found
  19s		19s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "pod-exists-demo-0"
```

Voila! `pod_exists` command has been synced to Icinga2. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `pod-exists-demo-0`.

![check-all-pods](/docs/images/cluster-alerts/pod_exists/demo-0.png)


### Check existence of a specific pod
In this tutorial, a ClusterAlert will be used check existence of a pod by name by setting `spec.vars.podName` field.
```yaml
$ cat ./docs/examples/cluster-alerts/pod_exists/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: pod-exists-demo-1
  namespace: demo
spec:
  check: pod_exists
  vars:
    podName: busybox
    count: 1
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/pod_exists/demo-1.yaml
pod "busybox" created
podalert "pod-exists-demo-1" created

$ kubectl get pods -n demo
NAME          READY     STATUS    RESTARTS   AGE
busybox       1/1       Running   0          5s

$ kubectl get podalert -n demo
NAME              KIND
pod-exists-demo-1   ClusterAlert.v1alpha1.monitoring.appscode.com

$ kubectl describe podalert -n demo pod-exists-demo-1
Name:		pod-exists-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  31s		31s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "pod-exists-demo-1". Reason: secrets "notifier-config" not found
  31s		31s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "pod-exists-demo-1"
  27s		27s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "pod-exists-demo-1"
```
![check-by-pod-label](/docs/images/cluster-alerts/pod_exists/demo-1.png)


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

#### Vars


#### Supported Icinga2 State

* OK
* CRITICAL
* UNKNOWN

#### Example
###### Command
```console
hyperalert check_pod_exists --host='pod_exists@default' --count=7
# --host is provided by Icinga2
```
###### Output
```
OK: Found all pods
```

##### Configure Alert Object
```yaml
# This will check if any pod exists in default namespace
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: check-pod-exist-1
  namespace: demo
spec:
  check: pod_exists
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]

# To check with expected pod number, suppose 8, add following in spec.vars
# vars:
#   count: 8

# To check for others kubernetes objects, set following labels
# labels:
#   alert.appscode.com/objectType: services
#   alert.appscode.com/objectName: elasticsearch-logging
```
