#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/searchlight"

pushd $REPO_ROOT

rm -rf dist

./hack/docker/searchlight/setup.sh
./hack/docker/searchlight/setup.sh push

./hack/docker/icinga/alpine/build.sh
./hack/docker/icinga/alpine/build.sh push

rm dist/.tag

popd
