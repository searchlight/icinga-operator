# Deployment Guide

This guide will walk you through deploying the Searchlight which includes Icinga2 and Controller.

## High Level Tasks

* Create the Alert Third Party Resource
* Create the Icinga2 Deployment
* Create the Searchlight Controller Deployment

## Deploying the Searchlight

#### Create the Third Party Resource

The `Searchlight` is driven by [Kubernetes Alert Objects](alert-resource/objects.md). Alert is not a core Kubernetes kind, but can be enabled with the Third Party Resource [Alert](alert-resource/third-party-resource.md):

#### Deploy Icinga2

Icinga2 is used as monitoring system which uses various plugins to check resources on Kubernetes. It also notifies users of outages and generates performance data for reporting.

See Icinga2 [Deployment Guide](icinga2/deployment.md).

#### Deploy Searchlight Controller

Searchlight Controller is used to communicate with Icinga2 API. To set an alert, create [Kubernetes Alert Objects](alert-resource/objects.md) with relevant information. Controller will consume that alert object. 
 
```sh
# Create Deployment
kubectl apply -f https://raw.githubusercontent.com/appscode/searchlight/master/hack/kubernetes/searchlight/deployment.yaml
```
