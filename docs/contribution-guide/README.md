# Release Process

The following steps must be done from a Linux x64 bit machine.

- Do a global replacement of tags so that docs point to the next release.
- Push changes to the release-x branch and apply new tag.
- Push all the changes to remote repo.
- Now, first build all the binaries:
```sh
$ cd ~/go/src/github.com/appscode/searchlight
$ ./hack/make.py build; env APPSCODE_ENV=prod ./hack/make.py release; ./hack/make.py push
```
- Build and push searchlight docker image
```sh
./hack/docker/searchlight/setup.sh; env APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release
```
- Build and push both forms of icinga image:
```sh
./hack/docker/icinga/build.sh; ./hack/docker/icinga/build.sh release
./hack/docker/icinga/setup.sh; ./hack/docker/icinga/setup.sh release
```
- Now, update the release notes in Github. See previous release notes to get an idea what to include there.


Now, you should probably also release a new version of kubed. These steps are:
- Revendor kubed so that new changes become available.
- Build kubed. Add any flags if needed.
- Push changes to release branch.
- Build and release kubed docker image.
- Now update Kubernetes salt stack files so that the new kubed image is used.

