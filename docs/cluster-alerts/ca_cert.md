# Check ca_cert

ClusterAlert `ca_cert` checks the expiration timestamp of Kubernetes api server CA certificate. No longer you have to get a surprise that the CA certificate for your cluster has expired.

## Spec

#### Vars
`ca_cert` check command has the following variables:
- `warning` - Condition for warning, compare with tiem left before expiration. (Default: TTL < 360h)
- `critical` - Condition for critical, compare with tiem left before expiration. (Default: TTL < 120h)

Execution of this command can result in following states:
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
In this tutorial, we are going to create an alert to check `ca_cert`.
```yaml
$ cat ./docs/examples/cluster-alerts/ca_cert/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: ca-cert-demo-0
  namespace: demo
spec:
  check: ca_cert
  vars:
    warning: 240h
    critical: 72h
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: any-notifier
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/cluster-alerts/ca_cert/demo-0.yaml 
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

Voila! `ca_cert` command has been synced to Icinga2. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `ca-cert-demo-0`.

![Demo of check_ca_cert](/docs/images/cluster-alerts/ca_cert/demo-0.png)

### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps

