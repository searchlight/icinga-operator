#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT=$GOPATH/src/github.com/appscode/searchlight
source "$REPO_ROOT/hack/libbuild/common/lib.sh"
source "$REPO_ROOT/hack/libbuild/common/public_image.sh"

IMG=icinga
ICINGAWEB_VER=2.1.2

DIST=$REPO_ROOT/dist
mkdir -p $DIST
if [ -f "$DIST/.tag" ]; then
	export $(cat $DIST/.tag | xargs)
fi

clean() {
    pushd $REPO_ROOT/hack/docker/icinga/jessie
	rm -rf icingaweb2 plugins
	popd
}

build() {
    pushd $REPO_ROOT/hack/docker/icinga/jessie
    detect_tag $DIST/.tag

	rm -rf icingaweb2
	clone https://github.com/Icinga/icingaweb2.git
	cd icingaweb2
	git checkout tags/v$ICINGAWEB_VER
	cd ..

	rm -rf plugins; mkdir -p plugins
	gsutil cp gs://appscode-dev/binaries/hyperalert/$TAG/hyperalert-linux-amd64 plugins/hyperalert
	chmod 755 plugins/*

	local cmd="docker build -t appscode/$IMG:$TAG-k8s ."
	echo $cmd; $cmd
	popd
}

docker_push() {
	docker_up $IMG:$TAG-k8s
}

docker_release() {
	local cmd="docker push appscode/$IMG:$TAG-k8s"
	echo $cmd; $cmd
}

binary_repo $@
