# Build Instructions

## Build Binary
```
# Install/Update dependency (needs glide)
glide slow

# Build
./hack/make.py build searchlight
```

## Build Docker
```
# Build Docker image
# This will build Searchlight Controller Binary and use it in docker
./hack/docker/searchlight/setup.sh
```

### Upload

#### For AppsCode dev team only

###### Release Docker Image
```
# The release refers to a public repository [docker.io/appscode/searchlight]
./hack/docker/searchlight/setup.sh release

# This script uses "docker push appscode/searchlight:<tag>"
```


###### Push Docker Image
```
# The push refers to two private Repositories
# 1. Google Container Registry [gcr.io/tigerworks-kube/searchlight]
# 2. Appscode Artifactory [docker.appscode.com/searchlight]

./hack/docker/searchlight/setup.sh push

# This script uses following two commands for gcr:
# 1. "docker tag appscode/searchlight:<tag> gcr.io/tigerworks-kube/searchlight:<tag>"
# 2. "gcloud docker -- push gcr.io/tigerworks-kube/searchlight:<tag>"

# And following two commands for artifactory:
# 1. "docker tag appscode/searchlight:<tag> docker.appscode.com/searchlight:<tag>"
# 2. "docker push docker.appscode.com/searchlight:<tag>"

```


#### For public use

###### Push Docker Image
```
# This will push docker image to other repositories

# Add docker tag for your repository
docker tag appscode/searchlight:<tag> <image>:<tag>

# Push Image
docker push <image>:<tag>

# Example:
docker tag appscode/searchlight:default aerokite/searchlight:default
docker push aerokite/searchlight:default
```
