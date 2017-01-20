# Deployment Guide

This guide will walk you through deploying the Searchlight which includes Controller & Icinga2.

## High Level Tasks

* Create the Alert Third Party Resource
* Create the Icinga2 Deployment
* Create the Searchlight Controller Deployment

## Deploying the Searchlight

### Create the Alert Third Party Resource

The `Searchlight` is driven by [Kubernetes Alert Objects](docs/alert-resource/objects.md). Alert is not a core Kubernetes kind, but can be enabled with the [Alert Third Party Resource](docs/alert-resource/third-party-resource.md):

### Deploy Icinga2

Icinga2 is used as monitoring system which uses various plugins to check resources. It also notifies users of outages and generates performance data for reporting.

See Icinga2 [Deployment Guide](docs/icinga/deployment.md).

### Deploy Searchlight Controller

Searchlight Controller is used to communicate with Icinga2 API. To set an alert, create [Kubernetes Alert Objects](docs/alert-resource/objects.md) with relevant information.
 
See Searchlight Controller [Deployment Guide](docs/searchlight/deployment.md).