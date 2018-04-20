---
title: Pod Volume
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: pod-pod-volume
    name: Pod Volume
    parent: pod-alert
    weight: 40
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check pod-volume

Check command `pod-volume` is used to check percentage of available space in Kubernetes Pods.

## Spec
`pod-volume` check command has the following variables:

- `volumeName` - Name of volume whose usage stats will be checked
- `secretName` - Name of Kubernetes Secret used to pass [hostfacts auth info](/docs/setup/hostfacts.md#create-hostfacts-secret)
- `warning` - Warning level value (usage percentage defaults to 80.0)
- `critical` - Critical level value (usage percentage defaults to 95.0)

Execution of this command can result in following states:

- OK
- Warning
- Critical
- Unknown


## Tutorial

### Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install Searchlight operator in your cluster following the steps [here](/docs/setup/install.md). To use `pod-volume` command, please also [install Hostfacts](/docs/setup/hostfacts.md) server in your cluster.

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

### Check volume of pods with matching labels
In this tutorial, a PodAlert will be used check volume stats of pods with matching labels by setting `spec.selector` field.
```yaml
$ cat ./docs/examples/pod-alerts/pod-volume/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: pod-volume-demo-0
  namespace: demo
spec:
  selector:
    matchLabels:
      app: nginx
  check: pod-volume
  vars:
    volumeName: www
    warning: '70'
    critical: '95'
  checkInterval: 5m
  alertInterval: 3m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod-volume/demo-1.yaml
persistentvolumeclaim "boxclaim" created
pod "busybox" created
podalert "pod-volume-demo-1" created

$ kubectl get pods -n demo
NAME      READY     STATUS    RESTARTS   AGE
web-0     1/1       Running   0          10m
web-1     1/1       Running   0          10m

$ kubectl describe podalert -n demo pod-volume-demo-0
Name:		pod-volume-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  11m		11m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-volume-demo-0". Reason: secrets "notifier-config" not found
  11m		11m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-volume-demo-0"
  11m		11m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-volume-demo-0"
```

Voila! `pod-volume` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@pod@minikube` and Icinga service `pod-volume-demo-0`.

![check-pods-by-label](/docs/images/pod-alerts/pod-volume/demo-0.png)


### Check volume stats of a specific pod
In this tutorial, a PodAlert will be used check volume stats of a pod by name by setting `spec.podName` field.

```yaml
$ cat ./docs/examples/pod-alerts/pod-volume/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: pod-volume-demo-1
  namespace: demo
spec:
  podName: busybox
  check: pod-volume
  vars:
    volumeName: mypd
    warning: '70'
    critical: '95'
  checkInterval: 5m
  alertInterval: 3m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod-volume/demo-1.yaml
persistentvolumeclaim "boxclaim" created
pod "busybox" created
podalert "pod-volume-demo-1" created

$ kubectl describe podalert -n demo pod-volume-demo-1
Name:		pod-volume-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  1m		1m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for PodAlert: "pod-volume-demo-1". Reason: secrets "notifier-config" not found
  1m		1m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-volume-demo-1"
  1m		1m		1	Searchlight operator			Normal		SuccessfulSync	Applied PodAlert: "pod-volume-demo-1"
```
![check-by-pod-name](/docs/images/pod-alerts/pod-volume/demo-1.png)


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
