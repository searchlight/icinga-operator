# Deployment Guide

This guide will walk you through deploying the icinga.

### Build Icinga

You can build Docker Image by yourself and use it.

See build instruction [here](docs/icinga/appscode/build.md).

### Deploy Icinga

###### Create Secret
```
# Encode Icinga Secret data in base64 format
cat deployments/icinga/appscode/secret-data.ini | base64

# Use encoded data in deployments/icinga/appscode/secret.yaml as .env value
# Create Secret
kubectl create -f deployments/icinga/appscode/secret.yaml
```

###### Create Service
```
# Create Service
kubectl create -f deployments/icinga/appscode/service.yaml
```

###### Create Deployment
```
# Create Deployment
kubectl create -f deployments/icinga/appscode/deployment.yaml
```

### Login

To login into `Icingaweb2`, use `username` & `password` from your AppsCode organization.
