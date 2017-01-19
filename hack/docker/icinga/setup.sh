#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

LIB_ROOT=$(dirname "${BASH_SOURCE}")/../../..
source "$LIB_ROOT/hack/libbuild/common/lib.sh"
source "$LIB_ROOT/hack/libbuild/common/public_image.sh"

GOPATH=$(go env GOPATH)
IMG=icinga
ICINGA_VER=2.4.8
K8S_VER=1.5
ICINGAWEB_VER=2.1.2

DIST=$GOPATH/src/github.com/appscode/searchlight/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
	export $(cat $DIST/.tag | xargs)
fi

clean() {
    pushd $GOPATH/src/github.com/appscode/searchlight/hack/docker/icinga
	rm -rf icingaweb2 plugins
	popd
}

build() {
    pushd $GOPATH/src/github.com/appscode/searchlight/hack/docker/icinga
    detect_tag $DIST/.tag

	rm -rf icingaweb2
	clone git@diffusion.appscode.com:appscode/79/icingaweb.git icingaweb2
	cd icingaweb2
	checkout apicss
	cd ..

	rm -rf plugins; mkdir -p plugins
	gsutil cp gs://appscode-dev/binaries/hyperalert/$TAG/hyperalert-linux-amd64 plugins/hyperalert
	chmod 755 plugins/*

	local cmd="docker build -t appscode/$IMG:$TAG-ac ."
	echo $cmd; $cmd
	popd
}

docker_push() {
	docker_up $IMG:$TAG-ac
}

docker_release() {
	local cmd="docker push appscode/$IMG:$TAG-ac"
	echo $cmd; $cmd
}

binary_repo $@
