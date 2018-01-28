#!/bin/bash
set -eou pipefail

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export SEARCHLIGHT_NAMESPACE=kube-system
export SEARCHLIGHT_SERVICE_ACCOUNT=default
export SEARCHLIGHT_ENABLE_RBAC=false
export SEARCHLIGHT_RUN_ON_MASTER=0
export SEARCHLIGHT_DOCKER_REGISTRY=appscode
export SEARCHLIGHT_IMAGE_PULL_SECRET=

show_help() {
    echo "searchlight.sh - install searchlight operator"
    echo " "
    echo "searchlight.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                         show brief help"
    echo "-n, --namespace=NAMESPACE          specify namespace (default: kube-system)"
    echo "    --rbac                         create RBAC roles and bindings"
    echo "    --docker-registry              docker registry used to pull searchlight images (default: appscode)"
    echo "    --image-pull-secret            name of secret used to pull searchlight operator images"
    echo "    --run-on-master                run searchlight operator on master"
}

while test $# -gt 0; do
    case "$1" in
        -h|--help)
            show_help
            exit 0
            ;;
        -n)
            shift
            if test $# -gt 0; then
                export SEARCHLIGHT_NAMESPACE=$1
            else
                echo "no namespace specified"
                exit 1
            fi
            shift
            ;;
        --namespace*)
            export SEARCHLIGHT_NAMESPACE=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --docker-registry*)
            export SEARCHLIGHT_DOCKER_REGISTRY=`echo $1 | sed -e 's/^[^=]*=//g'`
            shift
            ;;
        --image-pull-secret*)
            secret=`echo $1 | sed -e 's/^[^=]*=//g'`
            export SEARCHLIGHT_IMAGE_PULL_SECRET="name: '$secret'"
            shift
            ;;
        --rbac)
            export SEARCHLIGHT_SERVICE_ACCOUNT=searchlight-operator
            export SEARCHLIGHT_ENABLE_RBAC=true
            shift
            ;;
        --run-on-master)
            export SEARCHLIGHT_RUN_ON_MASTER=1
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

env | sort | grep SEARCHLIGHT*
echo ""

curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/5.1.0/hack/deploy/operator.yaml | envsubst | kubectl apply -f -

if [ "$SEARCHLIGHT_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT --namespace $SEARCHLIGHT_NAMESPACE
    kubectl label serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/5.1.0/hack/deploy/rbac-list.yaml | envsubst | kubectl auth reconcile -f -
fi

if [ "$SEARCHLIGHT_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy searchlight-operator -n $SEARCHLIGHT_NAMESPACE \
      --patch="$(curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/5.1.0/hack/deploy/run-on-master.yaml)"
fi
