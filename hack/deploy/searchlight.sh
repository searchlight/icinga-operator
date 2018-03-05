#!/bin/bash
set -eou pipefail

echo "checking kubeconfig context"
kubectl config current-context || { echo "Set a context (kubectl use-context <context>) out of the following:"; echo; kubectl config get-contexts; exit 1; }
echo ""

# https://stackoverflow.com/a/677212/244009
if [ -x "$(command -v onessl >/dev/null 2>&1)" ]; then
    export ONESSL=onessl
else
    # ref: https://stackoverflow.com/a/27776822/244009
    case "$(uname -s)" in
        Darwin)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-darwin-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        Linux)
            curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-linux-amd64
            chmod +x onessl
            export ONESSL=./onessl
            ;;

        CYGWIN*|MINGW32*|MSYS*)
            curl -fsSL -o onessl.exe https://github.com/kubepack/onessl/releases/download/0.1.0/onessl-windows-amd64.exe
            chmod +x onessl.exe
            export ONESSL=./onessl.exe
            ;;
        *)
            echo 'other OS'
            ;;
    esac
fi

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup {
    rm -rf $ONESSL ca.crt ca.key server.crt server.key
}
trap cleanup EXIT

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export SEARCHLIGHT_NAMESPACE=kube-system
export SEARCHLIGHT_SERVICE_ACCOUNT=searchlight-operator
export SEARCHLIGHT_ENABLE_RBAC=true
export SEARCHLIGHT_RUN_ON_MASTER=0
export SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK=false
export SEARCHLIGHT_DOCKER_REGISTRY=appscode
export SEARCHLIGHT_IMAGE_PULL_SECRET=

KUBE_APISERVER_VERSION=$(kubectl version -o=json | $ONESSL jsonpath '{.serverVersion.gitVersion}')
$ONESSL semver --check='>=1.9.0' $KUBE_APISERVER_VERSION
if [ $? -eq 0 ]; then
    export SEARCHLIGH_ENABLE_ADMISSION_WEBHOOK=true
fi

show_help() {
    echo "searchlight.sh - install searchlight operator"
    echo " "
    echo "searchlight.sh [options]"
    echo " "
    echo "options:"
    echo "-h, --help                         show brief help"
    echo "-n, --namespace=NAMESPACE          specify namespace (default: kube-system)"
    echo "    --rbac                         create RBAC roles and bindings (default: true)"
    echo "    --docker-registry              docker registry used to pull searchlight images (default: appscode)"
    echo "    --image-pull-secret            name of secret used to pull searchlight operator images"
    echo "    --run-on-master                run searchlight operator on master"
    echo "    --enable-admission-webhook     configure admission webhook for searchlight CRDs"
    echo "    --uninstall                    uninstall searchlight"
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
        --enable-admission-webhook*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK=false
            else
                export SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK=true
            fi
            shift
            ;;
        --rbac*)
            val=`echo $1 | sed -e 's/^[^=]*=//g'`
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_SERVICE_ACCOUNT=default
                export SEARCHLIGHT_ENABLE_RBAC=false
            fi
            shift
            ;;
        --run-on-master)
            export SEARCHLIGHT_RUN_ON_MASTER=1
            shift
            ;;
        --uninstall)
            export SEARCHLIGHT_UNINSTALL=1
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

if [ "$SEARCHLIGHT_UNINSTALL" -eq 1 ]; then
    kubectl delete deployment -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete service -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete secret -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete apiservice -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete validatingwebhookconfiguration -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete mutatingwebhookconfiguration -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    # Delete RBAC objects, if --rbac flag was used.
    kubectl delete serviceaccount -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete clusterrolebindings -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete clusterrole -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete rolebindings -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete role -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE

    exit 0
fi

env | sort | grep SEARCHLIGHT*
echo ""

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

curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/operator.yaml | $ONESSL envsubst | kubectl apply -f -

if [ "$SEARCHLIGHT_ENABLE_RBAC" = true ]; then
    kubectl create serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT --namespace $SEARCHLIGHT_NAMESPACE
    kubectl label serviceaccount $SEARCHLIGHT_SERVICE_ACCOUNT app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/rbac-list.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/user-roles.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
fi

if [ "$SEARCHLIGHT_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy searchlight-operator -n $SEARCHLIGHT_NAMESPACE \
      --patch="$(curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/run-on-master.yaml)"
fi

if [ "$SEARCHLIGHT_ENABLE_ADMISSION_WEBHOOK" = true ]; then
    curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-alpha.0/hack/deploy/admission.yaml | $ONESSL envsubst | kubectl apply -f -
fi

echo
echo "Successfully installed Searchlight!"
