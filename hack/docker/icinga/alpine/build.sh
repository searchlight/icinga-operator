#!/bin/bash

# Copyright AppsCode Inc. and Contributors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/go.searchlight.dev/icinga-operator
source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

APPSCODE_ENV=${APPSCODE_ENV:-dev}
IMG=icinga
ICINGAWEB_VER=2.4.1

DIST=$REPO_ROOT/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
    export $(cat $DIST/.tag | xargs)
fi

clean() {
    pushd $REPO_ROOT/hack/docker/icinga/alpine
    rm -rf icingaweb2 plugins
    popd
}

build() {
    pushd $REPO_ROOT/hack/docker/icinga/alpine
    detect_tag $DIST/.tag

    rm -rf icingaweb2
    clone https://github.com/Icinga/icingaweb2.git
    cd icingaweb2
    git checkout tags/v$ICINGAWEB_VER
    cd ..

    rm -rf plugins
    mkdir -p plugins
    if [ -f $DIST/hyperalert/hyperalert-alpine-amd64 ]; then
        cp $DIST/hyperalert/hyperalert-alpine-amd64 plugins/hyperalert
    else
        gsutil cp gs://appscode-dev/binaries/hyperalert/$TAG/hyperalert-alpine-amd64 plugins/hyperalert
    fi
    chmod 755 plugins/*

    local cmd="docker build -t $DOCKER_REGISTRY/$IMG:$TAG-k8s ."
    echo $cmd
    $cmd

    rm -rf icingaweb2 plugins
    popd
}

docker_push() {
    if [ "$APPSCODE_ENV" = "prod" ]; then
        echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"
        exit 1
    fi
    if [ "$TAG_STRATEGY" = "git_tag" ]; then
        echo "Are you trying to 'release' binaries to prod?"
        exit 1
    fi
    TAG=$TAG-k8s hub_canary
}

docker_release() {
    if [ "$APPSCODE_ENV" != "prod" ]; then
        echo "'release' only works in PROD env."
        exit 1
    fi
    if [ "$TAG_STRATEGY" != "git_tag" ]; then
        echo "'apply_tag' to release binaries and/or docker images."
        exit 1
    fi
    TAG=$TAG-k8s hub_up
}

binary_repo $@
