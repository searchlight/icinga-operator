#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

RETVAL=0

deps() {
    go get github.com/appscode/api/...
    local -r owd=$PWD
    cd $GOPATH/src/github.com/appscode/api
    ./hack/gen.sh
    cd $owd
}

build() {
    local -r owd=$PWD
    cd $GOPATH/src/github.com/appscode/client
    goimports -w cli credential *.go
    gofmt -s -w cli credential *.go
    go build ./...
    cd $owd
}

if [ $# -eq 0 ]; then
	build
	exit $RETVAL
fi

case "$1" in
	deps)
		deps
		;;
	build)
		build
		;;
	*)	(10)
		echo $"Usage: $0 {deps|build}"
		RETVAL=1
esac
exit $RETVAL
