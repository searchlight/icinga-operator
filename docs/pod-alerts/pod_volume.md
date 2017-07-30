# Check volume

PodAlert `pod_volume` checks the usage stats for of a volume of pods.


## Spec
`pod_volume` check command has the following variables:
- `volumeName` - Volume name
- `secret` - Kubernetes secret name for hostfacts authentication
- `warning` - Warning level value (usage percentage defaults to 75.0)
- `critical` - Critical level value (usage percentage defaults to 90.0)

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

To `pod_volume` command, you also have to deploy Hostfacts server in your cluster. Please follow the instructions [here](/docs/hostfacts.md) to deploy hostfacts.


### Create Alert
In this tutorial, we are going to create an alert to check `env`.
```yaml
$ cat ./docs/examples/pod-alerts/pod_volume/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: env-demo-0
  namespace: demo
spec:
  check: env
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
```
```console
$ kubectl apply -f ./docs/examples/pod-alerts/pod_volume/demo-0.yaml
clusteralert "env-demo-0" created

$ kubectl describe clusteralert env-demo-0 -n demo
Name:		env-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  6m		6m		1	Searchlight operator			Warning		BadNotifier	Bad notifier config for ClusterAlert: "env-demo-0". Reason: secrets "notifier-config" not found
  6m		6m		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "env-demo-0"
```

Voila! `env` command has been synced to Icinga2. Searchlight also logged a warning event, we have not created the notifier secret `notifier-config`. Please visit [here](/docs/tutorials/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `env-demo-0`.

![Demo of check_env](/docs/images/pod-alerts/pod_volume/demo-0.gif)

### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/uninstall.md).


## Next Steps



#### Hostfacts Secret keys

* `ca.crt`
* `hostfacts.key`
* `hostfacts.crt`
* `auth_token`
* `username`
* `password`

#### Example
###### Command
```console
hyperalert check_volume --host='monitoring-influxdb-0.12.2-n3lo2@kube-system' --volumeName=influxdb-persistent-storage --warning=70 --critical=85
# --host are provided by Icinga2
```
###### Output
```
WARNING: Disk used more than 70%
```

#### Required Hostfacts
Before using this CheckCommand, you must need to run `hostfacts` service in each Kubernetes node.
Volume stat of kubernetes pod is collected from `hostfacts` service.

See Hostfacts [deployment guide](hostfacts.md)


##### Configure Alert Object
```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: PodAlert
metadata:
  name: check-pod-volume-1
  namespace: kube-system
spec:
  check: volume
  alertInterval: 2m
  checkInterval: 1m
  receivers:
  - notifier: mailgun
    state: CRITICAL
    to: ["ops@example.com"]
  vars:
    volumeName: influxdb-persistent-storage
    warning: 70
    critical: 85
```
