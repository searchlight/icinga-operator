#!/bin/bash
set -eou pipefail

GOPATH=$(go env GOPATH)
REPO_ROOT="$GOPATH/src/github.com/appscode/searchlight"

export PROVIDER="minikube"
export STORAGE_CLASS="default"
export PROVIDED_ICINGA=""
export VERBOSITY=4

show_help() {
    echo "run.sh - run e2e test for searchlight"
    echo " "
    echo "run.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                show brief help"
    echo "    --provider            specify namespace (default: kube-system)"
    echo "    --storageclass        create RBAC roles and bindings (default: true)"
    echo "    --provided-icinga     docker registry used to pull searchlight images (default: appscode)"
    echo "    --v                   log version"

}

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        --provider*)
            export PROVIDER=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --storageclass*)
            export STORAGE_CLASS=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --provided-icinga*)
            export PROVIDED_ICINGA=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --v*)
            export VERBOSITY=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

ginkgo -r -v -progress -trace -race \
    -cover \
    -coverprofile="profile.out" \
    -outputdir="$REPO_ROOT/test/e2e" \
    test/e2e --\
    --provider="$PROVIDER" \
    --storageclass="$STORAGE_CLASS" \
    --provided-icinga="$PROVIDED_ICINGA" \
    --v="$VERBOSITY"
