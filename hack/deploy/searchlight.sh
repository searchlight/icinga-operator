#!/bin/bash
set -eou pipefail

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export SEARCHLIGHT_NAMESPACE=kube-system
export SEARCHLIGHT_SERVICE_ACCOUNT=default
export SEARCHLIGHT_ENABLE_RBAC=false
export SEARCHLIGHT_RUN_ON_MASTER=0
export SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK=false
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
    echo "    --enable-apiserver     configure admission webhook for searchlight CRDs"
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
        --enable-apiserver)
            export SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK=true
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

echo "checking kubeconfig context"
kubectl config current-context || { echo "Set a context (kubectl use-context <context>) out of the following:"; echo; kubectl config get-contexts; exit 1; }
echo ""

if [ "$SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK" = true ]; then
    # ref: https://stackoverflow.com/a/27776822/244009
    case "$(uname -s)" in
        Darwin)
            curl -fsSL -o onessl https://github.com/appscode/onessl/releases/download/0.1.0/onessl-darwin-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        Linux)
            curl -fsSL -o onessl https://github.com/appscode/onessl/releases/download/0.1.0/onessl-linux-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        CYGWIN*|MINGW32*|MSYS*)
            curl -fsSL -o onessl.exe https://github.com/appscode/onessl/releases/download/0.1.0/onessl-windows-amd64.exe
            chmod +x onessl.exe
            export ONESSL=./onessl.exe
            ;;
        *)
            echo 'other OS'
            ;;
    esac

    # create necessary TLS certificates:
    # - a local CA key and cert
    # - a webhook server key and cert signed by the local CA
    $ONESSL create ca-cert
    $ONESSL create server-cert server --domains=searchlight-operator.$SEARCHLIGHT_NAMESPACE.svc
    export SERVICE_SERVING_CERT_CA=$(cat ca.crt | $ONESSL base64)
    export TLS_SERVING_CERT=$(cat server.crt | $ONESSL base64)
    export TLS_SERVING_KEY=$(cat server.key | $ONESSL base64)
    export KUBE_CA=$($ONESSL get kube-ca | $ONESSL base64)
    rm -rf $ONESSL ca.crt ca.key server.crt server.key

    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/admission/operator.yaml | envsubst | kubectl apply -f -
else
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/operator.yaml | envsubst | kubectl apply -f -
fi

if [ "$SEARCHLIGHT_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT --namespace $SEARCHLIGHT_NAMESPACE
    kubectl label serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/rbac-list.yaml | envsubst | kubectl auth reconcile -f -

    if [ "$SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK" = true ]; then
        curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/admission/rbac-list.yaml | envsubst | kubectl auth reconcile -f -
    fi
fi

if [ "$SEARCHLIGHT_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy searchlight-operator -n $SEARCHLIGHT_NAMESPACE \
      --patch="$(curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/run-on-master.yaml)"
fi
