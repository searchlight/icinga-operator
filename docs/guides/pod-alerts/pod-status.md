---
title: Pod Status
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: pod-pod-status
    name: Pod Status
    parent: pod-alert
    weight: 30
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check pod-status

Check command `pod-status` is used to check status of Kubernetes pods. Returns OK if `status.phase` of a pod is `Succeeded` or `Running`, otherwise, returns Critical.


## Spec
`pod-status` check command has no variables. Execution of this command can result in following states:

- OK
- Critical
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

### Check status of pods with matching labels
In this tutorial, a PodAlert will be used check status of pods with matching labels by setting `spec.selector` field.
```yaml
$ cat ./docs/examples/pod-alerts/pod-status/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: pod-status-demo-0
  namespace: demo
spec:
  selector:
    matchLabels:
      app: nginx
  check: pod-status
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod-status/demo-0.yaml
replicationcontroller "nginx" created
podalert "pod-status-demo-0" created

$ kubectl get pods -n demo
NAME          READY     STATUS    RESTARTS   AGE
nginx-c0v51   1/1       Running   0          53s
nginx-vqhzv   1/1       Running   0          53s

$ kubectl describe podalert -n demo pod-status-demo-0
Name:		pod-status-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  3m		3m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-status-demo-0". Reason: secrets "notifier-config" not found
  3m		3m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-0"
  3m		3m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-0"
```

Voila! `pod-status` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@pod@minikube` and Icinga service `pod-status-demo-0`.

![check-all-pods](/docs/images/pod-alerts/pod-status/demo-0.png)


### Check status of a specific pod
In this tutorial, a PodAlert will be used check status of a pod by name by setting `spec.podName` field.
```yaml
$ cat ./docs/examples/pod-alerts/pod-status/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: pod-status-demo-1
  namespace: demo
spec:
  podName: busybox
  check: pod-status
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod-status/demo-1.yaml
pod "busybox" created
podalert "pod-status-demo-1" created

$ kubectl get pods -n demo
NAME          READY     STATUS    RESTARTS   AGE
busybox       1/1       Running   0          5s

$ kubectl describe podalert -n demo pod-status-demo-1
Name:		pod-status-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  31s		31s		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-status-demo-1". Reason: secrets "notifier-config" not found
  31s		31s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-1"
  27s		27s		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-status-demo-1"
```
![check-by-pod-label](/docs/images/pod-alerts/pod-status/demo-1.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md).
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
