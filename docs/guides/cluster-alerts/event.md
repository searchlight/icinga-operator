---
title: Event
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-event
    name: Event
    parent: cluster-alert
    weight: 30
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

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
- Warning
- Unknown

## Tutorial

### Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install Searchlight operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

$ kubectl get namespaces
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
  - notifier: Mailgun
    state: Warning
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

Voila! `event` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `event-demo-0`.

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
  - notifier: Mailgun
    state: Warning
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/event/demo-1.yaml
pod "busybox" created
clusteralert "event-demo-1" created

$ kubectl describe clusteralert -n demo event-demo-1
Name:		event-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  13s		13s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "event-demo-1". Reason: secrets "notifier-config" not found
  13s		13s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "event-demo-1"

$ kubectl get events -n demo
LASTSEEN   FIRSTSEEN   COUNT     NAME      KIND      SUBOBJECT                  TYPE      REASON       SOURCE              MESSAGE
19s        19s         1         busybox   Pod                                  Normal    Scheduled    default-scheduler   Successfully assigned busybox to minikube
3s         18s         5         busybox   Pod       spec.containers{busybox}   Normal    Pulled       kubelet, minikube   Container image "busybox" already present on machine
18s        18s         1         busybox   Pod       spec.containers{busybox}   Normal    Created      kubelet, minikube   Created container with id fa40f0698ed44706774a504be7b0eb6bd776b67082e76c4376a432b04c6e4f26
18s        18s         1         busybox   Pod       spec.containers{busybox}   Warning   Failed       kubelet, minikube   Failed to start container with id fa40f0698ed44706774a504be7b0eb6bd776b67082e76c4376a432b04c6e4f26 with error: rpc error: code = 2 desc = failed to start container "fa40f0698ed44706774a504be7b0eb6bd776b67082e76c4376a432b04c6e4f26": Error response from daemon: Container command 'bad' not found or does not exist.
18s        18s         1         busybox   Pod                                  Warning   FailedSync   kubelet, minikube   Error syncing pod, skipping: failed to "StartContainer" for "busybox" with rpc error: code = 2 desc = failed to start container "fa40f0698ed44706774a504be7b0eb6bd776b67082e76c4376a432b04c6e4f26": Error response from daemon: Container command 'bad' not found or does not exist.: "Start Container Failed"
17s        17s         1         busybox   Pod       spec.containers{busybox}   Normal    Created      kubelet, minikube   Created container with id 6774cff0d7e0b9487ad87fd2dd712028102d0f2854be47561898cbe72cf10e4d
17s        17s         1         busybox   Pod       spec.containers{busybox}   Warning   Failed       kubelet, minikube   Failed to start container with id 6774cff0d7e0b9487ad87fd2dd712028102d0f2854be47561898cbe72cf10e4d with error: rpc error: code = 2 desc = failed to start container "6774cff0d7e0b9487ad87fd2dd712028102d0f2854be47561898cbe72cf10e4d": Error response from daemon: Container command 'bad' not found or does not exist.
17s        17s         1         busybox   Pod                                  Warning   FailedSync   kubelet, minikube   Error syncing pod, skipping: failed to "StartContainer" for "busybox" with rpc error: code = 2 desc = failed to start container "6774cff0d7e0b9487ad87fd2dd712028102d0f2854be47561898cbe72cf10e4d": Error response from daemon: Container command 'bad' not found or does not exist.: "Start Container Failed"
6s         6s          1         busybox   Pod       spec.containers{busybox}   Normal    Created      kubelet, minikube   Created container with id 49455114b6bc1a626eb217fcf23cd1172dfd03d75d7f4650fbe52dd5940a1b24
5s         5s          1         busybox   Pod       spec.containers{busybox}   Warning   Failed       kubelet, minikube   Failed to start container with id 49455114b6bc1a626eb217fcf23cd1172dfd03d75d7f4650fbe52dd5940a1b24 with error: rpc error: code = 2 desc = failed to start container "49455114b6bc1a626eb217fcf23cd1172dfd03d75d7f4650fbe52dd5940a1b24": Error response from daemon: Container command 'bad' not found or does not exist.
5s         5s          1         busybox   Pod                                  Warning   FailedSync   kubelet, minikube   Error syncing pod, skipping: failed to "StartContainer" for "busybox" with rpc error: code = 2 desc = failed to start container "49455114b6bc1a626eb217fcf23cd1172dfd03d75d7f4650fbe52dd5940a1b24": Error response from daemon: Container command 'bad' not found or does not exist.: "Start Container Failed"
4s         4s          1         busybox   Pod       spec.containers{busybox}   Normal    Created      kubelet, minikube   Created container with id a6e4a7734965f039faf72d60c131fe312a974b08100e787295ea1a5bb6cc3806
4s         4s          1         busybox   Pod       spec.containers{busybox}   Warning   Failed       kubelet, minikube   Failed to start container with id a6e4a7734965f039faf72d60c131fe312a974b08100e787295ea1a5bb6cc3806 with error: rpc error: code = 2 desc = failed to start container "a6e4a7734965f039faf72d60c131fe312a974b08100e787295ea1a5bb6cc3806": Error response from daemon: Container command 'bad' not found or does not exist.
4s         4s          1         busybox   Pod                                  Warning   FailedSync   kubelet, minikube   Error syncing pod, skipping: failed to "StartContainer" for "busybox" with rpc error: code = 2 desc = failed to start container "a6e4a7734965f039faf72d60c131fe312a974b08100e787295ea1a5bb6cc3806": Error response from daemon: Container command 'bad' not found or does not exist.: "Start Container Failed"
3s         3s          1         busybox   Pod       spec.containers{busybox}   Normal    Created      kubelet, minikube   Created container with id 8ef27cb9fd83b61a6a99b838bc55fb61b1f76c33f0a55b25b104ccb08e743e28
3s         3s          1         busybox   Pod       spec.containers{busybox}   Warning   Failed       kubelet, minikube   Failed to start container with id 8ef27cb9fd83b61a6a99b838bc55fb61b1f76c33f0a55b25b104ccb08e743e28 with error: rpc error: code = 2 desc = failed to start container "8ef27cb9fd83b61a6a99b838bc55fb61b1f76c33f0a55b25b104ccb08e743e28": Error response from daemon: Container command 'bad' not found or does not exist.
3s         3s          1         busybox   Pod                                  Warning   FailedSync   kubelet, minikube   Error syncing pod, skipping: failed to "StartContainer" for "busybox" with rpc error: code = 2 desc = failed to start container "8ef27cb9fd83b61a6a99b838bc55fb61b1f76c33f0a55b25b104ccb08e743e28": Error response from daemon: Container command 'bad' not found or does not exist.: "Start Container Failed"
```
![check-by-pod-label](/docs/images/cluster-alerts/event/demo-1.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - To periodically run various checks on pods in a Kubernetes cluster, use [PodAlerts](/docs/concepts/alert-types/pod-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
