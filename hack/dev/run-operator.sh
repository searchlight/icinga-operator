#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/searchlight"

pushd $REPO_ROOT

searchlight run \
  --v=6 \
  --secure-port=8443 \
  --kubeconfig="$HOME/.kube/config" \
  --authorization-kubeconfig="$HOME/.kube/config" \
  --authentication-kubeconfig="$HOME/.kube/config" \
  --authentication-skip-lookup \
  --config-dir="$GOPATH/src/github.com/appscode/icinga-testconfig

popd
