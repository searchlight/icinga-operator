# Build Instructions


### Build default Icinga

Default Icinga also includes `hyperalert` plugin.
```
gsutil cp gs://appscode-dev/binaries/hyperalert/<tag>/hyperalert-linux-amd64 plugins/hyperalert
```

We can add `hyperalert` plugin in Icinga downloaded from anywhere. We just need  to add plugin in plugins directory and name it as `hyperalert`.
We can also modified `hack/docker/icinga/build.sh` script to do this.

##### Build

```
# Build Docker image
./hack/docker/icinga/build.sh
```

### Upload

#### For AppsCode dev team only

###### Release Docker Image
```
# The release refers to a public repository [docker.io/appscode/icinga]
./hack/docker/icinga/build.sh release

# This script uses "docker push appscode/icinga:<tag>-k8s"
```


###### Push Docker Image
```
# The push refers to two private Repositories
# 1. Google Container Registry [gcr.io/tigerworks-kube/icinga]
# 2. Appscode Artifactory [docker.appscode.com/icinga]

./hack/docker/icinga/build.sh push

# This script uses following two commands for gcr:
# 1. "docker tag appscode/icinga:<tag>-k8s gcr.io/tigerworks-kube/icinga:<tag>-k8s"
# 2. "gcloud docker -- push gcr.io/tigerworks-kube/icinga:<tag>-k8s"

# And following two commands for artifactory:
# 1. "docker tag appscode/icinga:<tag>-k8s docker.appscode.com/icinga:<tag>-k8s"
# 2. "docker push docker.appscode.com/icinga:<tag>-k8s"
```

#### For public use

###### Push Docker Image
```
# This will push docker image to any repositories

# Add docker tag to image for your repository
docker tag appscode/icinga:<tag>-k8s <repository>:<tag>-k8s

# Push Image
docker push <repository>:<tag>-k8s

# Example:
docker tag appscode/icinga:default-k8s aerokite/icinga:default-k8s
docker push aerokite/icinga:default-k8s
```
