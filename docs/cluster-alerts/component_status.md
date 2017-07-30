> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Check component_status

Check command `component_status` is used to check status of Kubernetes components. Returns OK if components are `Healthy`, otherwise, returns CRITICAL.


## Spec
`component_status` has the following variables:
- `selector` - Label selector for components whose existence are checked.
- `componentName` - Name of Kubernetes component whose existence is checked.

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

### Check status of all components
In this tutorial, we are going to create a ClusterAlert to check status of all components.
```yaml
$ cat ./docs/examples/cluster-alerts/component_status/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: component-status-demo-0
  namespace: demo
spec:
  check: component_status
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/component_status/demo-0.yaml
clusteralert "component-status-demo-0" created

$ kubectl describe clusteralert -n demo component-status-demo-0
Name:		component-status-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  6s		6s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "component-status-demo-0"
```

Voila! `component_status` command has been synced to Icinga2. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `component-status-demo-0`.

![check-all-components](/docs/images/cluster-alerts/component_status/demo-0.png)


### Check status of a specific component
In this tutorial, a ClusterAlert will be used check status of a component by name by setting `spec.componentName` field.

```yaml
$ cat ./docs/examples/cluster-alerts/component_status/demo-1.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: component-status-demo-1
  namespace: demo
spec:
  check: component_status
  vars:
    componentName: etcd-0
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```

```console
$ kubectl apply -f ./docs/examples/cluster-alerts/component_status/demo-1.yaml
clusteralert "component-status-demo-1" created

$ kubectl describe clusteralert -n demo component-status-demo-1
Name:		component-status-demo-1
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  22s		22s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "component-status-demo-1"
```
![check-by-component-name](/docs/images/cluster-alerts/component_status/demo-1.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps
