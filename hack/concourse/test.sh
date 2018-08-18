#!/usr/bin/env bash

set -eoux pipefail

ORG_NAME=appscode
REPO_NAME=searchlight
OPERATOR_NAME=searchlight
APP_LABEL=searchlight #required for `kubectl describe deploy -n kube-system -l app=$APP_LABEL`

export APPSCODE_ENV=test-concourse
export DOCKER_REGISTRY=appscodeci

# docker images to delete after running tests
DOCKER_IMG_NAMES=()

# get concourse-common
#pushd $REPO_NAME
#git status # required, otherwise you'll get error `Working tree has modifications.  Cannot add.`. why?
#git subtree pull --prefix hack/libbuild https://github.com/appscodelabs/libbuild.git master --squash -m 'concourse'
#popd

source $REPO_NAME/hack/libbuild/concourse/init.sh

pushd $GOPATH/src/github.com/$ORG_NAME/$REPO_NAME

# install dependencies
./hack/builddeps.sh
go get -u golang.org/x/tools/cmd/goimports
go get github.com/Masterminds/glide
go get github.com/sgotti/glide-vc
go get github.com/onsi/ginkgo/ginkgo
go install github.com/onsi/ginkgo/ginkgo

./hack/make.py build searchlight
./hack/make.py build hyperalert

./hack/docker/icinga/alpine/build.sh
./hack/docker/icinga/alpine/build.sh push

DOCKER_IMG_NAMES+=(icinga)

./hack/docker/searchlight/setup.sh
./hack/docker/searchlight/setup.sh push

DOCKER_IMG_NAMES+=(searchlight)

source ./hack/deploy/searchlight.sh --docker-registry=$DOCKER_REGISTRY --enable-validating-webhook=true --rbac=true --icinga-api-password=1234
./hack/make.py test e2e --searchlight-service=searchlight-operator@kube-system --provider=do
#./hack/make.py test e2e --searchlight-service=searchlight-operator@kube-system --provider=$ClusterProvider
popd
