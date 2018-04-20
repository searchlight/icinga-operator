---
title: CA Cert
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-ca-cert
    name: CA Cert
    parent: cluster-alert
    weight: 20
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check ca-cert

Check command `ca-cert` checks the expiration timestamp of Kubernetes api server CA certificate. No longer you have to get a surprise that the CA certificate for your cluster has expired.

## Spec
`ca-cert` check command has the following variables:

- `warning` - Condition for warning, compare with tiem left before expiration. (Default: TTL < 360h)
- `critical` - Condition for critical, compare with tiem left before expiration. (Default: TTL < 120h)

Execution of this command can result in following states:

- OK
- Warning
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

### Create Alert
In this tutorial, we are going to create an alert to check `ca-cert`.
```yaml
$ cat ./docs/examples/cluster-alerts/ca-cert/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: ca-cert-demo-0
  namespace: demo
spec:
  check: ca-cert
  vars:
    warning: 240h
    critical: 72h
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/ca-cert/demo-0.yaml 
clusteralert "ca-cert-demo-0" created

$ kubectl describe clusteralert ca-cert-demo-0 -n demo
Name:		ca-cert-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  9s		9s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "ca-cert-demo-0"
```

Voila! `ca-cert` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `ca-cert-demo-0`.

![check ca-cert](/docs/images/cluster-alerts/ca-cert/demo-0.png)

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
