---
title: Release | Searchlight
description: Searchlight Release
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: release
    name: Release
    parent: developer-guide
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---
# Release Process

The following steps must be done from a Linux x64 bit machine.

- Do a global replacement of tags so that docs point to the next release.
- Push changes to the release-x branch and apply new tag.
- Push all the changes to remote repo.
- Now, first build all the binaries:
```console
$ cd ~/go/src/github.com/appscode/searchlight
$ ./hack/make.py build; env APPSCODE_ENV=prod ./hack/make.py push; ./hack/make.py push
```
- Build and push searchlight docker image
```console
./hack/docker/searchlight/setup.sh; env APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release
```
- Build and push both forms of icinga image:
```console
./hack/docker/icinga/alpine/build.sh; env APPSCODE_ENV=prod ./hack/docker/icinga/alpine/build.sh release
./hack/docker/icinga/alpine/setup.sh; env APPSCODE_ENV=prod ./hack/docker/icinga/alpine/setup.sh release
```
- Now, update the release notes in Github. See previous release notes to get an idea what to include there.
