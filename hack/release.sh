#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/searchlight"

export APPSCODE_ENV=prod

pushd $REPO_ROOT

rm -rf dist

./hack/docker/searchlight/setup.sh
APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release
./hack/make.py push

./hack/docker/icinga/alpine/build.sh
APPSCODE_ENV=prod ./hack/docker/icinga/alpine/build.sh release

rm dist/.tag

popd

# ./hack/docker/icinga/alpine/setup.sh
# ./hack/docker/icinga/alpine/setup.sh release

# ./hack/docker/postgres/build.sh
# APPSCODE_ENV=prod ./hack/docker/postgres/build.sh release
