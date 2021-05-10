# Searchlight
[Searchlight by AppsCode](https://go.searchlight.dev/icinga-operator) is an alert manager for Kubernetes built around Icinga2.

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install appscode/searchlight --name searchlight-operator --namespace kube-system
```

## Introduction

This chart bootstraps a [Searchlight controller](https://go.searchlight.dev/icinga-operator) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart
To install the chart with the release name `searchlight-operator`:

```console
$ helm install appscode/searchlight --name searchlight-operator
```

The command deploys Searchlight operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `searchlight-operator`:

```console
$ helm delete searchlight-operator
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Searchlight chart and their default values.

| Parameter                            | Description                                                             | Default            |
| ------------------------------------ | -----------------------------------------------------------------       | ------------------ |
| `replicaCount`                       | Number of searchlight Operator replicas to create (only 1 is supported) | `1`                |
| `operator.registry`                  | Docker registry used to pull Operator image                             | `appscode`         |
| `operator.repository`                | Operator container image                                                | `searchlight`      |
| `operator.tag`                       | Operator image tag                                                      | `8.0.0`       |
| `icinga.registry`                    | Docker registry used to pull Icinga image                               | `appscode`         |
| `icinga.repository`                  | Icinga container image                                                  | `icinga`           |
| `icinga.tag`                         | icinga container image tag                                              | `8.0.0-k8s`   |
| `ido.registry`                       | Docker registry used to pull PostgreSQL image                           | `appscode`         |
| `ido.repository`                     | PostgreSQL container image                                              | `postgress`        |
| `ido.tag`                            | ido container image tag                                                 | `9.5-alpine`       |
| `imagePullSecrets`                   | Specify image pull secrets                                              | `nil` (does not add image pull secrets to deployed pods) |
| `imagePullPolicy`                    | Image pull policy                                                       | `IfNotPresent`     |
| `criticalAddon`                      | If true, installs Searchlight operator as critical addon                | `false`            |
| `logLevel`                           | Log level for operator                                                  | `3`                |
| `affinity`                           | Affinity rules for pod assignment                                       | `{}`               |
| `nodeSelector`                       | Node labels for pod assignment                                          | `{}`               |
| `nodeSelector`                       | Node labels for pod assignment                                          | `{}`               |
| `tolerations`                        | Tolerations used pod assignment                                         | `{}`               |
| `rbac.create`                        | If `true`, create and use RBAC resources                                | `true`             |
| `serviceAccount.create`              | If `true`, create a new service account                                 | `true`             |
| `serviceAccount.name`                | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `` |
| `apiserver.groupPriorityMinimum`     | The minimum priority the group should have.                             | 10000              |
| `apiserver.versionPriority`          | The ordering of this API inside of the group.                           | 15                 |
| `apiserver.enableValidatingWebhook`  | Enable validating webhooks for Searchlight CRDs                         | false              |
| `apiserver.ca`                       | CA certificate used by main Kubernetes api server                       | ``                 |
| `apiserver.disableStatusSubresource` | If true, uses status sub resource for Searchlight crds                  | `false`            |
| `enableAnalytics`                    | Send usage events to Google Analytics                                   | `true`             |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name searchlight-operator --set image.tag=v0.2.1 appscode/searchlight
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install --name searchlight-operator --values values.yaml appscode/searchlight
```

## RBAC
By default the chart will not install the recommended RBAC roles and rolebindings.

You need to have the flag `--authorization-mode=RBAC` on the api server. See the following document for how to enable [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/).

To determine if your cluster supports RBAC, run the following command:

```console
$ kubectl api-versions | grep rbac
```

If the output contains "beta", you may install the chart with RBAC enabled (see below).

### Enable RBAC role/rolebinding creation

To enable the creation of RBAC resources (On clusters with RBAC). Do the following:

```console
$ helm install --name searchlight-operator appscode/searchlight --set rbac.create=true
```
