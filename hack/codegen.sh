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

set -x

GOPATH=$(go env GOPATH)
PACKAGE_NAME=go.searchlight.dev/icinga-operator
REPO_ROOT="$GOPATH/src/$PACKAGE_NAME"
DOCKER_REPO_ROOT="/go/src/$PACKAGE_NAME"
DOCKER_CODEGEN_PKG="/go/src/k8s.io/code-generator"
apiGroups=(incidents/v1alpha1 monitoring/v1alpha1)

pushd $REPO_ROOT

## Generate ugorji stuff
rm "$REPO_ROOT"/apis/monitoring/v1alpha1/*.generated.go
mkdir -p "$REPO_ROOT"/api/api-rules

# for EAS types
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.14 "$DOCKER_CODEGEN_PKG"/generate-internal-groups.sh "deepcopy,defaulter,conversion" \
    go.searchlight.dev/icinga-operator/client \
    go.searchlight.dev/icinga-operator/apis \
    go.searchlight.dev/icinga-operator/apis \
    incidents:v1alpha1 \
    --go-header-file "$DOCKER_REPO_ROOT/hack/gengo/boilerplate.go.txt"

# for both CRD and EAS types
docker run --rm -ti -u $(id -u):$(id -g) \
    -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
    -w "$DOCKER_REPO_ROOT" \
    appscode/gengo:release-1.14 "$DOCKER_CODEGEN_PKG"/generate-groups.sh all \
    go.searchlight.dev/icinga-operator/client \
    go.searchlight.dev/icinga-operator/apis \
    "incidents:v1alpha1 monitoring:v1alpha1" \
    --go-header-file "$DOCKER_REPO_ROOT/hack/gengo/boilerplate.go.txt"

# Generate openapi
for gv in "${apiGroups[@]}"; do
    docker run --rm -ti -u $(id -u):$(id -g) \
        -v "$REPO_ROOT":"$DOCKER_REPO_ROOT" \
        -w "$DOCKER_REPO_ROOT" \
        appscode/gengo:release-1.14 openapi-gen \
        --v 1 --logtostderr \
        --go-header-file "hack/gengo/boilerplate.go.txt" \
        --input-dirs "$PACKAGE_NAME/apis/${gv},k8s.io/apimachinery/pkg/apis/meta/v1,k8s.io/apimachinery/pkg/api/resource,k8s.io/apimachinery/pkg/runtime,k8s.io/apimachinery/pkg/version" \
        --output-package "$PACKAGE_NAME/apis/${gv}" \
        --report-filename api/api-rules/violation_exceptions.list
done

# Generate crds.yaml, plugins.yaml and swagger.json
go run ./hack/gencrd/main.go
go run ./hack/genplugin/main.go

popd
