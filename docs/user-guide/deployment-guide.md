# Deployment Guide

This guide will walk you through deploying the Searchlight which includes Icinga2 and Controller.

## High Level Tasks

* Create the Alert Third Party Resource
* Create the Icinga2 Deployment
* Create the Searchlight Controller Deployment

## Deploying the Searchlight

#### Create the Third Party Resource

The `Searchlight` is driven by [Kubernetes Alert Objects](alert-object.md). `Alert` is not a core Kubernetes kind, but can be enabled with following Third Party Resource.
```yaml
# Third Party Resource `Alert`
metadata:
  name: alert.appscode.com
apiVersion: extensions/v1beta1
kind: ThirdPartyResource
description: "Alert support for Kubernetes by appscode.com"
versions:
  - name: v1beta1
```

```sh
# Create Third Party Resource
kubectl apply -f https://raw.githubusercontent.com/appscode/k8s-addons/master/api/extensions/alert.yaml
```

#### Deploy Icinga2

Icinga2 is used as monitoring system which uses various plugins to check resources on Kubernetes. It also notifies users of outages and generates performance data for reporting.

See Icinga2 [Deployment Guide](icinga2/deployment.md).

Run following command to deploy Icinga2
```sh
curl https://raw.githubusercontent.com/appscode/searchlight/1.5.2/hack/deploy/icinga2/run.sh | bash
```

> Make sure you have set notifier to send notifications. Check [this](icinga2/deployment.md#create-deployment).

#### Deploy Searchlight Controller

Searchlight Controller is used to communicate with Icinga2 API. To set an alert, create [Kubernetes Alert Objects](alert-object.md) with relevant information. Controller will consume that alert object. 
 
```sh
# Create Deployment
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/1.5.2/hack/deploy/searchlight/deployment.yaml
```
