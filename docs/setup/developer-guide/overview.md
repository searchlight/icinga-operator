---
title: Overview | Developer Guide
description: Developer Guide Overview
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: developer-guide-readme
    name: Overview
    parent: developer-guide
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

## Development Guide
This document is intended to be the canonical source of truth for things like supported toolchain versions for building Searchlight.
If you find a requirement that this doc does not capture, please submit an issue on github.

This document is intended to be relative to the branch in which it is found. It is guaranteed that requirements will change over time
for the development branch, but release branches of Searchlight should not change.

### Build Searchlight
Some of the Searchlight development helper scripts rely on a fairly up-to-date GNU tools environment, so most recent Linux distros should
work just fine out-of-the-box.

#### Setup GO
Searchlight is written in Google's GO programming language. Currently, Searchlight is developed and tested on **go 1.9.2**. If you haven't set up a GO
development environment, please follow [these instructions](https://golang.org/doc/code.html) to install GO.

#### Download Source

```console
$ go get github.com/appscode/searchlight
$ cd $(go env GOPATH)/src/github.com/appscode/searchlight
```

#### Install Dev tools
To install various dev tools for Searchlight, run the following command:

```console
$ ./hack/builddeps.sh
```

#### Build Binary
```
$ ./hack/make.py
$ searchlight version
```

#### Run Binary Locally
```console
$ searchlight run \
  --secure-port=8443 \
  --kubeconfig="$HOME/.kube/config" \
  --authorization-kubeconfig="$HOME/.kube/config" \
  --authentication-kubeconfig="$HOME/.kube/config" \
  --authentication-skip-lookup
```

#### Dependency management
Searchlight uses [Glide](https://github.com/Masterminds/glide) to manage dependencies. Dependencies are already checked in the `vendor` folder. If you want to update/add dependencies, run:

```console
$ glide slow
```


#### Build Operator Docker image
To build and push your custom Docker image, follow the steps below. To release a new version of Searchlight, please follow the [release guide](/docs/setup/developer-guide/release.md).

```console
# Build Docker image
$ ./hack/docker/searchlight/setup.sh; ./hack/docker/searchlight/setup.sh push

# Add docker tag for your repository
$ docker tag appscode/searchlight:<tag> <image>:<tag>

# Push Image
$ docker push <image>:<tag>

# Example:
docker tag appscode/searchlight:default aerokite/searchlight:default
docker push aerokite/searchlight:default
```


#### Build Icinga Docker image

Default Icinga also includes `hyperalert` plugin.
```console
gsutil cp gs://appscode-dev/binaries/hyperalert/<tag>/hyperalert-linux-amd64 plugins/hyperalert
```

We can add `hyperalert` plugin in Icinga downloaded from anywhere. We just need  to add plugin in plugins directory and name it as `hyperalert`.

```console
# Build Docker image
./hack/docker/icinga/build.sh

# This will push docker image to any repositories

# Add docker tag to image for your repository
docker tag appscode/icinga:<tag>-k8s <repository>:<tag>-k8s

# Push Image
docker push <repository>:<tag>-k8s

# Example:
docker tag appscode/icinga:default-k8s aerokite/icinga:default-k8s
docker push aerokite/icinga:default-k8s
```


#### Generate CLI Reference Docs
```console
$ ./hack/gendocs/make.sh
```

### Run e2e Test
Pass `storageclass` name as flag.
```console
$ ginkgo -v -r test/e2e -- --storageclass=standard
```
