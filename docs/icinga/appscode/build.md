# Build Instructions


### Build default Icinga

AppsCode Icinga also includes `hyperalert` plugin by default.

##### Build

```
# Build Docker image
./hack/docker/icinga/setup.sh
```

### Upload

#### For AppsCode dev team only

###### Release Docker Image
```
# The release refers to a public repository [docker.io/appscode/icinga]
./hack/docker/icinga/setup.sh release

# This script uses "docker push appscode/icinga:<tag>-ac"
```


###### Push Docker Image
```
# The push refers to two private Repositories
# 1. Google Container Registry [gcr.io/tigerworks-kube/icinga]
# 2. Appscode Artifactory [docker.appscode.com/icinga]

./hack/docker/icinga/setup.sh push

# This script uses following two commands for gcr:
# 1. "docker tag appscode/icinga:<tag>-ac gcr.io/tigerworks-kube/icinga:<tag>-ac"
# 2. "gcloud docker -- push gcr.io/tigerworks-kube/icinga:<tag>-ac"

# And following two commands for artifactory:
# 1. "docker tag appscode/icinga:<tag>-ac docker.appscode.com/icinga:<tag>-ac"
# 2. "docker push docker.appscode.com/icinga:<tag>-ac"
```

#### For public use

###### Push Docker Image
```
# This will push docker image to any repositories

# Add docker tag to image for your repository
docker tag appscode/icinga:<tag>-ac <repository>:<tag>-ac

# Push Image
docker push <repository>:<tag>-ac

# Example:
docker tag appscode/icinga:default-ac aerokite/icinga:default-ac
docker push aerokite/icinga:default-ac
```
