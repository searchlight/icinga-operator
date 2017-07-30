> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Check event

Check command `event` is used to check Kubernetes events. This plugin checks for all Warning events happened in the last `spec.checkInterval` duration.


## Spec
`event` check command has the following variables:

- `clockSkew` - Clock skew in Duration. [Default: 30s]. This time is added with `spec.checkInterval` while checking events
- `involvedObjectKind` - Kind of involved object used to select events
- `involvedObjectName` - Name of involved object used to select events
- `involvedObjectNamespace` - Namespace of involved object used to select events
- `involvedObjectUID` - UID of involved object used to select events

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


### Check existence of any warning event
In this tutorial, a ClusterAlert will be used check existence of warning events occurred in the last check interval.
```yaml
$ cat ./docs/examples/cluster-alerts/event/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: event-demo-0
  namespace: demo
spec:
  check: event
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: WARNING
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/event/demo-0.yaml 
replicationcontroller "nginx" created
clusteralert "event-demo-0" created

$ kubectl describe clusteralert -n demo event-demo-0
Name:		event-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  7s		7s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "event-demo-0". Reason: secrets "notifier-config" not found
  7s		7s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "event-demo-0"

$ kubectl get events -n demo
LASTSEEN   FIRSTSEEN   COUNT     NAME           KIND           SUBOBJECT                TYPE      REASON           SOURCE                 MESSAGE
15s        15s         1         nginx-9n8z7    Pod                                     Normal    Scheduled        default-scheduler      Successfully assigned nginx-9n8z7 to minikube
15s        15s         1         nginx-9n8z7    Pod            spec.containers{nginx}   Normal    Pulling          kubelet, minikube      pulling image "nginx:bad"
12s        12s         1         nginx-9n8z7    Pod            spec.containers{nginx}   Warning   Failed           kubelet, minikube      Failed to pull image "nginx:bad": rpc error: code = 2 desc = Tag bad not found in repository docker.io/library/nginx
12s        12s         1         nginx-9n8z7    Pod                                     Warning   FailedSync       kubelet, minikube      Error syncing pod, skipping: failed to "StartContainer" for "nginx" with ErrImagePull: "rpc error: code = 2 desc = Tag bad not found in repository docker.io/library/nginx"
12s        12s         1         nginx-9n8z7    Pod            spec.containers{nginx}   Normal    BackOff          kubelet, minikube      Back-off pulling image "nginx:bad"
12s        12s         1         nginx-9n8z7    Pod                                     Warning   FailedSync       kubelet, minikube      Error syncing pod, skipping: failed to "StartContainer" for "nginx" with ImagePullBackOff: "Back-off pulling image \"nginx:bad\""
```

Voila! `event` command has been synced to Icinga2. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `event-demo-0`.

![check-all-pods](/docs/images/cluster-alerts/event/demo-0.png)


### Check existence of events for a specific object
In this tutorial, a ClusterAlert will be used check existence of events for a specific object by setting one or more `spec.vars.involvedObject*` fields.
```yaml
$ cat ./docs/examples/cluster-alerts/event/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: event-demo-1
  namespace: demo
spec:
  check: event
  vars:
    involvedObjectName: busybox
    involvedObjectNamespace: demo
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: WARNING
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/event/demo-1.yaml
pod "busybox" created
podalert "event-demo-1" created

$ kubectl get pods -n demo
NAME          READY     STATUS    RESTARTS   AGE
busybox       1/1       Running   0          5s

$ kubectl get podalert -n demo
NAME              KIND
event-demo-1   ClusterAlert.v1alpha1.monitoring.appscode.com

$ kubectl describe podalert -n demo event-demo-1
Name:		event-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  31s		31s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "event-demo-1". Reason: secrets "notifier-config" not found
  31s		31s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "event-demo-1"
  27s		27s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "event-demo-1"
```
![check-by-pod-label](/docs/images/cluster-alerts/event/demo-1.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
