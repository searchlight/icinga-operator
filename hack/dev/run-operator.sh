#!/bin/bash
set -xeou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/searchlight"

pushd $REPO_ROOT

rm -rf ./hack/dev/testconfig/searchlight/pki

KUBE_NAMESPACE=demo searchlight run \
  --v=6 \
  --secure-port=8443 \
  --kubeconfig="$HOME/.kube/config" \
  --authorization-kubeconfig="$HOME/.kube/config" \
  --authentication-kubeconfig="$HOME/.kube/config" \
  --authentication-skip-lookup \
  --config-dir="$REPO_ROOT/hack/dev/testconfig"

popd
