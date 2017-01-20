# Deployment Guide

This guide will walk you through deploying the icinga.

### Build Icinga

You can build Docker Image by yourself and use it.

See build instruction [here](docs/icinga/k8s/build.md).

### Deploy Icinga

###### Create Secret
```
# Encode Icinga Secret data in base64 format
cat deployments/icinga/k8s/secret-data.ini | base64

# Use encoded data in deployments/icinga/k8s/secret.yaml as .env value
# Create Secret
kubectl create -f deployments/icinga/k8s/secret.yaml
```

###### Create Service
```
# Create Service
kubectl create -f deployments/icinga/k8s/service.yaml
```

###### Create Deployment
```
# Create Deployment
kubectl create -f deployments/icinga/k8s/deployment.yaml
```

### Login

To login into `Icingaweb2`, use following authentication information:
```
Username: admin
Password: <ICINGA_WEB_ADMIN_PASSWORD>
```
Password will be set from Icinga secret. See following example:
> `ICINGA_WEB_ADMIN_PASSWORD`=admin
