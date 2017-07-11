#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/searchlight
source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

APPSCODE_ENV=${APPSCODE_ENV:-dev}
IMG=postgres
TAG=9.5-alpine

DIST=$REPO_ROOT/dist

clean() {
    pushd $REPO_ROOT/hack/docker/postgres
	popd
}

build() {
	pushd $REPO_ROOT/hack/docker/postgres
	local cmd="docker build -t appscode/$IMG:$TAG ."
	echo $cmd; $cmd
	popd
}

docker_push() {
    if [ "$APPSCODE_ENV" = "prod" ]; then
        echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"
        exit 1
    fi
    hub_canary
}

docker_release() {
    if [ "$APPSCODE_ENV" != "prod" ]; then
        echo "'release' only works in PROD env."
        exit 1
    fi
    TAG=$TAG-k8s hub_up
}

binary_repo $@
