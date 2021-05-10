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

set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/go.searchlight.dev/icinga-operator"

export APPSCODE_ENV=prod

pushd $REPO_ROOT

rm -rf dist

./hack/docker/searchlight/setup.sh
APPSCODE_ENV=prod ./hack/docker/searchlight/setup.sh release
./hack/make.py push

./hack/docker/icinga/alpine/build.sh
APPSCODE_ENV=prod ./hack/docker/icinga/alpine/build.sh release

rm dist/.tag

popd

# ./hack/docker/icinga/alpine/setup.sh
# ./hack/docker/icinga/alpine/setup.sh release

# ./hack/docker/postgres/build.sh
# APPSCODE_ENV=prod ./hack/docker/postgres/build.sh release
