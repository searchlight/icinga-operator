---
title: Install
description: Searchlight Install
menu:
  product_searchlight_{{ .version }}:
    identifier: install-searchlight
    name: Install
    parent: setup
    weight: 10
product_name: searchlight
menu_name: product_searchlight_{{ .version }}
section_menu_id: setup
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Installation Guide

Searchlight operator can be installed via a script or as a Helm chart.
[![Install Searchlight](https://img.youtube.com/vi/Po4yXrQuHtQ/0.jpg)](https://www.youtube-nocookie.com/embed/Po4yXrQuHtQ)

<ul class="nav nav-tabs" id="installerTab" role="tablist">
  <li class="nav-item">
    <a class="nav-link active" id="script-tab" data-toggle="tab" href="#script" role="tab" aria-controls="script" aria-selected="true">Script</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="helm-tab" data-toggle="tab" href="#helm" role="tab" aria-controls="helm" aria-selected="false">Helm</a>
  </li>
</ul>
<div class="tab-content" id="installerTabContent">
  <div class="tab-pane fade show active" id="script" role="tabpanel" aria-labelledby="script-tab">

## Using Script

To install Searchlight in your Kubernetes cluster, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh | bash
```

After successful installation, you should have a `searchlight-operator-***` pod running in the `kube-system` namespace.

```console
$ kubectl get pods -n kube-system | grep searchlight-operator
searchlight-operator-6945bcd777-4jdv7   3/3       Running   0          2m
```

### Customizing Installer

The installer script and associated yaml files can be found in the [/hack/deploy](https://go.searchlight.dev/icinga-operator/tree/{{< param "info.version" >}}/hack/deploy) folder. You can see the full list of flags available to installer using `-h` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh | bash -s -- -h
searchlight.sh - install searchlight operator

searchlight.sh [options]

options:
-h, --help                         show brief help
-n, --namespace=NAMESPACE          specify namespace (default: kube-system)
    --rbac                         create RBAC roles and bindings (default: true)
    --docker-registry              docker registry used to pull searchlight images (default: appscode)
    --image-pull-secret            name of secret used to pull searchlight operator images
    --run-on-master                run searchlight operator on master
    --enable-validating-webhook    enable/disable validating webhooks for Searchlight crds
    --icinga-api-password          password used by icinga2 api (if unset, a random password will be generated and used)
    --enable-status-subresource    If enabled, uses status sub resource for Searchlight crds
    --enable-analytics             send usage events to Google Analytics (default: true)
    --uninstall                    uninstall searchlight
    --purge                        purges searchlight crd objects and crds
```

If you would like to run Searchlight operator pod in `master` instances, pass the `--run-on-master` flag:

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh \
    | bash -s -- --run-on-master [--rbac]
```

Searchlight operator will be installed in a `kube-system` namespace by default. If you would like to run Searchlight operator pod in `searchlight` namespace, pass the `--namespace=searchlight` flag:

```console
$ kubectl create namespace searchlight
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh \
    | bash -s -- --namespace=searchlight [--run-on-master] [--rbac]
```

If you are using a private Docker registry, you need to pull the following 3 docker images:

 - [appscode/searchlight](https://hub.docker.com/r/appscode/searchlight)
 - [appscode/icinga](https://hub.docker.com/r/appscode/icinga)
 - [appscode/postgres](https://hub.docker.com/r/appscode/postgres)

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace searchlight
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME] [--rbac]
```

Searchlight implements a [validating admission webhook](https://kubernetes.io/docs/admin/admission-controllers/#validatingadmissionwebhook-alpha-in-18-beta-in-19) to validate Searchlight CRDs. This is enabled by default for Kubernetes 1.9.0 or later releases. To disable this feature, pass the `--enable-validating-webhook=false` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/{{< param "info.version" >}}/hack/deploy/searchlight.sh \
    | bash -s -- --enable-admission-webhook [--rbac]
```

Searchlight {{< param "info.version" >}} or later releases can use status sub resource for CustomResourceDefintions. This is enabled by default for Kubernetes 1.11.0 or later releases. To disable this feature, pass the `--enable-status-subresource=false` flag.

</div>
<div class="tab-pane fade" id="helm" role="tabpanel" aria-labelledby="helm-tab">

## Using Helm
Searchlight can be installed via [Helm](https://helm.sh/) using the [chart](https://go.searchlight.dev/icinga-operator/blob/master/chart/searchlight) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search appscode/searchlight
NAME            CHART VERSION       APP VERSION DESCRIPTION
appscode/searchlight  {{< param "info.version" >}}    {{< param "info.version" >}}  Searchlight by AppsCode - Alerts for Kubernetes

$ helm install appscode/searchlight --name searchlight-operator --version {{< param "info.version" >}} --namespace kube-system
```

To see the detailed configuration options, visit [here](https://go.searchlight.dev/icinga-operator/tree/master/chart/searchlight).

</div>
</div>

### Installing in GKE Cluster

If you are installing Searchlight on a GKE cluster, you will need cluster admin permissions to install Searchlight operator. Run the following command to grant admin permision to the cluster.

```console
$ kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
  --clusterrole=cluster-admin \
  --user="$(gcloud config get-value core/account)"
```


## Verify installation
To check if Searchlight operator pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l app=searchlight --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.

Now, to confirm CRD groups have been registered by the operator, run the following command:
```console
$ kubectl get crd -l app=searchlight
```


## Accesing IcingaWeb2
Icinga comes with its own web dashboard called IcingaWeb. You can access IcingaWeb on your workstation by forwarding port `60006` of Searchlight operator pod.

```console
$ kubectl get pods --all-namespaces -l app=searchlight
NAME                                    READY     STATUS    RESTARTS   AGE
searchlight-operator-1987091405-ghj5b   3/3       Running   0          1m

$ kubectl port-forward searchlight-operator-1987091405-ghj5b -n kube-system 60006
Forwarding from 127.0.0.1:60006 -> 60006
E0728 04:07:28.237822   10898 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:60006: bind: cannot assign requested address
Handling connection for 60006
Handling connection for 60006
^C⏎
```

Now, open URL http://127.0.0.1:60006 on your browser. To login, use username `admin` and password `changeit`. If you want to change the password, read the next section.


## Configuring Icinga
Searchlight installation scripts above creates a Secret called `searchlight-operator` to store icinga configuration. This following keys are supported in this Secret.

| Key                    | Default Value  | Description                                             |
|------------------------|----------------|---------------------------------------------------------|
| ICINGA_WEB_UI_PASSWORD | _**changeit**_ | Password of `admin` user in IcingaWeb2                  |
| ICINGA_API_PASSWORD    | auto-generated | Password of icinga api user `icingaapi`                 |
| ICINGA_CA_CERT         | auto-generated | PEM encoded CA certificate used for icinga api endpoint |
| ICINGA_SERVER_CERT     | auto-generated | PEM encoded certificate used for icinga api endpoint    |
| ICINGA_SERVER_KEY      | auto-generated | PEM encoded private key used for icinga api endpoint    |
| ICINGA_IDO_PASSWORD    | auto-generated | Password of postgres user `icingaido`                   |
| ICINGA_WEB_PASSWORD    | auto-generated | Password of postgres user `icingaweb`                   |

To change the `admin` user login password in IcingaWeb, change the value of `ICINGA_WEB_UI_PASSWORD` key in Secret `searchlight-operator` and restart Searchlight operator pod(s).

```console
$ kubectl edit secret searchlight-operator -n kube-system
# Update the value of ICINGA_WEB_UI_PASSWORD key

$ kubectl get pods --all-namespaces -l app=searchlight
NAME                                    READY     STATUS    RESTARTS   AGE
searchlight-operator-1987091405-ghj5b   3/3       Running   0          1m

$ kubectl delete pods -n kube-system searchlight-operator-1987091405-ghj5b
pod "searchlight-operator-1987091405-ghj5b" deleted
```

## Configuring RBAC
Searchlight introduces the following Kubernetes objects:

| API Group                         | Kinds             |
|-----------------------------------|-------------------|
| monitoring.appscode.com           | `ClusterAlert`<br/>`NodeAlert`<br/>`PodAlert`<br/>`Incident` |
| incidents.monitoring.appscode.com | `Acknowledgement` |

Searchlight installer will create 3 user facing cluster roles:

| ClusterRole               | Aggregates To | Desription                            |
|---------------------------|---------------|---------------------------------------|
| appscode:searchlight:edit | admin         | Allows admin access to Searchlight objects, intended to be granted within a namespace using a RoleBinding. This grants ability to create incidents manually.|
| appscode:searchlight:edit | edit          | Allows edit access to Searchlight objects, intended to be granted within a namespace using a RoleBinding.      |
| appscode:searchlight:view | view          | Allows read-only access to Searchlight objects, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.

## Using kubectl
```console
# List all Searchlight objects
$ kubectl get clusteralerts,nodealerts,podalerts --all-namespaces
$ kubectl get ca,noa,poa --all-namespaces

# List Searchlight objects for a namespace
$ kubectl get clusteralerts,nodealerts,podalerts -n <namespace>
$ kubectl get ca,noa,poa -n <namespace>

# Get Searchlight object YAML
$ kubectl get podalert -n <namespace> <name> -o yaml
$ kubectl get poa -n <namespace> <name> -o yaml

# Describe Searchlight object. Very useful to debug problems.
$ kubectl describe podalert -n <namespace> <name>
$ kubectl describe poa -n <namespace> <name>
```

## Detect Searchlight version
To detect Searchlight version, exec into the operator pod and run `searchlight version` command.

```console
$ POD_NAMESPACE=kube-system
$ POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app=searchlight -o jsonpath={.items[0].metadata.name})
$ kubectl exec -it $POD_NAME -c operator -n $POD_NAMESPACE searchlight version

Version = {{< param "info.version" >}}
VersionStrategy = tag
Os = alpine
Arch = amd64
CommitHash = 9442863beb09a50a2c3818ab586fa5b1541fddf1
GitBranch = release-4.0
GitTag = {{< param "info.version" >}}
CommitTimestamp = 2017-09-26T03:00:58
```
