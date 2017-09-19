# Searchlight
[Searchlight by AppsCode](https://github.com/appscode/searchlight) is an alert manager for Kubernetes built around Icinga2.

## TL;DR;

```bash
$ helm install chart/searchlight
```

## Introduction

This chart bootstraps a [Searchlight controller](https://github.com/appscode/searchlight) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.5+

## Installing the Chart
To install the chart with the release name `my-release`:
```bash
$ helm install --name my-release chart/searchlight
```
The command deploys Searchlight controller on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release`:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the Searchlight chart and their default values.


| Parameter                 | Description                                                    | Default                |
|---------------------------|----------------------------------------------------------------|------------------------|
| `operator.image`          | operator container image                                       | `appscode/searchlight` |
| `operator.tag`            | operator image tag                                             | `3.0.1`                |
| `operator.pullPolicy`     | operator container image pull policy                           | `IfNotPresent`         |
| `icinga.image`            | icinga container image                                         | `appscode/icinga`      |
| `icinga.tag`              | icinga container image tag                                     | `3.0.1-k8s`            |
| `icinga.pullPolicy`       | icinga container image pull policy                             | `IfNotPresent`         |
| `ido.image`               | ido container image                                            | `appscode/postgress`   |
| `ido.tag`                 | ido container image tag                                        | `9.5-alpine`           |
| `ido.pullPolicy`          | ido container image pull policy                                | `IfNotPresent`         |
| `rbac.create`             | install required rbac service account, roles and rolebindings  | `false`                |
| `rbac.serviceAccountName` | ServiceAccount Voyager will use (ignored if rbac.create=true)  | `default`              |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name my-release --set image.tag=v0.2.1 chart/searchlight
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install --name my-release --values values.yaml chart/searchlight
```

## RBAC
By default the chart will not install the recommended RBAC roles and rolebindings.

You need to have the flag `--authorization-mode=RBAC` on the api server. See the following document for how to enable [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/).

To determine if your cluster supports RBAC, run the the following command:

```console
$ kubectl api-versions | grep rbac
```

If the output contains "beta", you may install the chart with RBAC enabled (see below).

### Enable RBAC role/rolebinding creation

To enable the creation of RBAC resources (On clusters with RBAC). Do the following:

```console
$ helm install --name my-release chart/searchlight --set rbac.create=true
```
