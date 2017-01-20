# Deployment Guide


Searchlight Controller consumes [Kubernetes Alert Objects](docs/alert-resource/objects.md) to create Icinga2 hosts, services and notifications. 

## Build Searchlight Controller

See build instruction [here](docs/searchlight/build.md).

### Deploy Searchlight Controller

###### Create Deployment
```
# Create Deployment
kubectl create -f deployments/searchlight/deployment.yaml
```

