---
title: Cert
menu:
  product_searchlight_8.0.0-rc.0:
    identifier: guides-cert
    name: Cert
    parent: cluster-alert
    weight: 20
product_name: searchlight
menu_name: product_searchlight_8.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check cert

Check command `cert` checks the expiration timestamp of any certificate from Secrets. No longer you have to get a surprise that your certificates have expired.

## Spec
`cert` check command has the following variables:

- `selector` - Selector (label query) to filter on, supports '=', '==', and '!='
- `secretName` - Name of secret from where certificates are checked
- `secretKey` - Name of secret key where certificates are kept
- `warning` - Remaining duration for Warning state. [Default: 360h]
- `critical` - Remaining duration for Critical state. [Default: 120h]

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

### Create a Secret

In this tutorial, we are going to use `onessl` to issue certificates. Download `onessl` from [kubepack/onessl](https://github.com/kubepack/onessl/releases).

```bash
$ onessl create ca-cert
$ onessl create server-cert
```

Now, we have two certificates `ca.crt` and `server.crt`.

Lets create a Secret with these Certificates.

```bash
$ kubectl create secret generic server-cert -n demo \
        --from-file=./ca.crt --from-file=./server.crt

secret "server-cert" created
```

```bash
$ kubectl get secret -n demo server-cert -o yaml
```

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: server-cert
  namespace: demo
type: Opaque
data:
  ca.crt: Y2EuY3J0Cg==
  server.crt: c2VydmVyLmNydAo=
```

### Create Alert
In this tutorial, we are going to create an alert to check certificates in Secret.

```yaml
$ cat ./docs/examples/cluster-alerts/cert/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: cert-demo-0
  namespace: demo
spec:
  check: cert
  vars:
    secretName: server-cert
    secretKey: "ca.crt,server.crt"
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

Here,

- `spec.check` provides check command name. In this case, it is `cert`.
- `spec.vars` supports following variables

    - `selector` - Label selector for secrets where certificates are stored. Supports '=', '==', and '!='
    - `secretName` - Name of secret from where certificates are checked.
    - `secretKey` - List of secret keys where certificates are kept
    - `warning` - Remaining duration for Warning state. [Default: 360h]
    - `critical` - Remaining duration for Critical state. [Default: 120h]

```console
$ kubectl apply -f ./docs/examples/cluster-alerts/cert/demo-0.yaml
clusteralert "cert-demo-0" created

$ kubectl describe clusteralert cert-demo-0 -n demo
Name:		cert-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  9s		9s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "cert-demo-0"
```

Voila! `cert` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `ca-cert-demo-0`.

Following notes are important:

- If `secretName` and `selector` both are not provided, all secrets in same namespace will be checked.
- If `secretKey` is not provided in the alert, and SecretType of a secret is `SecretTypeTLS`, TLS certificate in `tls.crt"` will be checked.

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
