# Build Instructions


### Build default Icinga

Default Icinga also includes `hyperalert` plugin.
```sh
gsutil cp gs://appscode-dev/binaries/hyperalert/<tag>/hyperalert-linux-amd64 plugins/hyperalert
```

We can add `hyperalert` plugin in Icinga downloaded from anywhere. We just need  to add plugin in plugins directory and name it as `hyperalert`.
We can also modified `hack/docker/icinga/build.sh` script to do this.

##### Build

```sh
# Build Docker image
./hack/docker/icinga/build.sh
```

###### Push Docker Image
```sh
# This will push docker image to any repositories

# Add docker tag to image for your repository
docker tag appscode/icinga:<tag>-k8s <repository>:<tag>-k8s

# Push Image
docker push <repository>:<tag>-k8s

# Example:
docker tag appscode/icinga:default-k8s aerokite/icinga:default-k8s
docker push aerokite/icinga:default-k8s
```
