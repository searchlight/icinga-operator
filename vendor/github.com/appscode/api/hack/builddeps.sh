#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

RETVAL=0
ROOT=$PWD

setup_protoc() {
    script=$(mktemp /tmp/protoc-installer.XXXXXX)
	cat >$script <<EOF
#!/bin/bash -e
set -o errexit
set -o nounset
set -o pipefail

# ref: http://stackoverflow.com/a/17072017
if [ "$(uname)" == "Darwin" ]; then
    xcode-select --install
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
    apt-get -y install curl unzip git build-essential automake libtool
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ]; then
    echo "Windows NT platform is not supported."
fi

rm -rf /opt/grpc
mkdir -p /opt
pushd /opt
# git clone -b $(curl -L http://grpc.io/release) https://github.com/grpc/grpc
git clone https://github.com/grpc/grpc
cd grpc
git submodule update --init
echo "setting up protoc"
cd /opt/grpc/third_party/protobuf
./autogen.sh && ./configure && make
make install
ldconfig
echo "setting up grpc"
cd /opt/grpc
make
make install
popd
EOF
    if sudo -n true 2>/dev/null; then
        sudo bash "$script"
    else
        bash "$script"
    fi
	rm "$script"
}

setup_proxy() {
	echo "Setting up grpc proxy"
	rm -rf $GOPATH/src/github.com/googleapis/googleapis
	go get -u github.com/googleapis/googleapis || true
	rm -rf $GOPATH/src/github.com/googleapis/googleapis/third_party
	rm -rf $GOPATH/src/google.golang.org/genproto
	go get google.golang.org/genproto || true
	rm -rf $GOPATH/src/google.golang.org/grpc
	go get -u google.golang.org/grpc
	pushd $GOPATH/src/google.golang.org/grpc
	git checkout v1.3.0
	popd
	rm -rf $GOPATH/src/github.com/golang/protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go
	mkdir -p $GOPATH/src/github.com/grpc-ecosystem
	pushd $GOPATH/src/github.com/grpc-ecosystem
	if [ ! -d grpc-gateway ]; then
		git clone git@github.com:appscode/grpc-gateway.git
	fi
	cd grpc-gateway
	git fetch origin
	git checkout v1.2.2-ac
	git reset --soft HEAD~20
	git reset HEAD --hard
	git pull origin v1.2.2-ac
	go install ./protoc-gen-grpc-gateway/...
	go install ./protoc-gen-grpc-gateway-cors/...
	go install ./protoc-gen-grpc-js-client/...
	go install ./protoc-gen-swagger/...
	popd
}

setup() {
	setup_protoc
	setup_proxy
}

if [ $# -eq 0 ]; then
	setup_proxy
	exit $RETVAL
fi

case "$1" in
	protoc)
		setup_protoc
		;;
	gateway)
		setup_proxy
		;;
	all)
		setup
		;;
	*)  echo $"Usage: $0 {protoc|gateway|all}"
		RETVAL=1
		;;
esac
exit $RETVAL
