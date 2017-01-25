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
