#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
pushd $GOPATH/src/github.com/appscode/searchlight

./hack/make.py build
APPSCODE_ENV=prod ./hack/make.py push
./hack/make.py push

./hack/docker/searchlight/setup.sh
APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release

./hack/docker/icinga/alpine/build.sh
APPSCODE_ENV=prod ./hack/docker/icinga/alpine/build.sh release

# ./hack/docker/icinga/alpine/setup.sh
# ./hack/docker/icinga/alpine/setup.sh release

# ./hack/docker/postgres/build.sh
# APPSCODE_ENV=prod ./hack/docker/postgres/build.sh release
