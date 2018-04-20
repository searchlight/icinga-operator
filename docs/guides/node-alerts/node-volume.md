---
title: Node Volume
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: node-volume
    name: Node Volume
    parent: node-alert
    weight: 30
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check node-volume

Check command `node-volume` is used to check percentage of available space in Kubernetes Nodes.

## Spec
`node-volume` check command has the following variables:

- `mountpoint` - Mountpoint of volume whose usage stats will be checked
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

Now, install Searchlight operator in your cluster following the steps [here](/docs/setup/install.md). To use `node-volume` command, please also [install Hostfacts](/docs/setup/hostfacts.md) server in your cluster.

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

### Check volume stats of all nodes
In this tutorial, we are going to create a NodeAlert to check volume stats of all nodes.
```yaml
$ cat ./docs/examples/node-alerts/node-volume/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: node-volume-demo-0
  namespace: demo
spec:
  check: node-volume
  vars:
    mountpoint: /mnt/sda1
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
$ kubectl apply -f ./docs/examples/node-alerts/node-volume/demo-0.yaml
nodealert "node-volume-demo-0" created

$ kubectl describe nodealert -n demo node-volume-demo-0
Name:		node-volume-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  6s		6s		1	Searchlight operator			Normal		SuccessfulSync	Applied NodeAlert: "node-volume-demo-0"
```

Voila! `node-volume` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@node@minikube` and Icinga service `node-volume-demo-0`.

![check-all-nodes](/docs/images/node-alerts/node-volume/demo-0.png)


### Check volume stats of nodes with matching labels
In this tutorial, a NodeAlert will be used check volume stats of nodes with matching labels by setting `spec.selector` field.

```yaml
$ cat ./docs/examples/node-alerts/node-volume/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: node-volume-demo-1
  namespace: demo
spec:
  selector:
    beta.kubernetes.io/os: linux
  check: node-volume
  vars:
    mountpoint: /mnt/sda1
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
$ kubectl apply -f ./docs/examples/node-alerts/node-volume/demo-1.yaml
nodealert "node-volume-demo-1" created

$ kubectl describe nodealert -n demo node-volume-demo-1
Name:		node-volume-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  33s		33s		1	Searchlight operator			Normal		SuccessfulSync	Applied NodeAlert: "node-volume-demo-1"
```
![check-by-node-label](/docs/images/node-alerts/node-volume/demo-1.png)


### Check volume stats of a specific node
In this tutorial, a NodeAlert will be used check volume stats of a node by name by setting `spec.nodeName` field.

```yaml
$ cat ./docs/examples/node-alerts/node-volume/demo-2.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: NodeAlert
metadata:
  name: node-volume-demo-2
  namespace: demo
spec:
  nodeName: minikube
  check: node-volume
  vars:
    mountpoint: /mnt/sda1
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
$ kubectl apply -f ./docs/examples/node-alerts/node-volume/demo-2.yaml
nodealert "node-volume-demo-2" created

$ kubectl describe nodealert -n demo node-volume-demo-2
Name:		node-volume-demo-2
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  22s		22s		1	Searchlight operator			Normal		SuccessfulSync	Applied NodeAlert: "node-volume-demo-2"
```
![check-by-node-name](/docs/images/node-alerts/node-volume/demo-2.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - To periodically run various checks on a Kubernetes cluster, use [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md).
 - To periodically run various checks on pods in a Kubernetes cluster, use [PodAlerts](/docs/concepts/alert-types/pod-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
