#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/searchlight

source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

IMG=icinga
ICINGA_VER=2.4.8
K8S_VER=1.5
ICINGAWEB_VER=2.1.2

mkdir -p $REPO_ROOT/dist
if [ -f "$REPO_ROOT/dist/.tag" ]; then
  export $(cat $REPO_ROOT/dist/.tag | xargs)
fi

clean() {
  pushd $REPO_ROOT/hack/docker/icinga/alpine
  rm -rf icingaweb2 plugins
  popd
}

build() {
  pushd $REPO_ROOT/hack/docker/icinga/alpine
  detect_tag $REPO_ROOT/dist/.tag

  rm -rf icingaweb2
  clone git@diffusion.appscode.com:appscode/79/icingaweb.git icingaweb2
  cd icingaweb2
  checkout 2.1.2-ac
  cd ..

  rm -rf plugins
  mkdir -p plugins
  gsutil cp gs://appscode-dev/binaries/hyperalert/$TAG/hyperalert-linux-amd64 plugins/hyperalert
  chmod 755 plugins/*

  local cmd="docker build -t appscode/$IMG:$TAG-ac ."
  echo $cmd; $cmd

  rm -rf icingaweb2 plugins
  popd
}

docker_push() {
  TAG="$TAG-ac" attic_up
}

docker_release() {
  TAG="$TAG-ac" hub_up
}

binary_repo $@
