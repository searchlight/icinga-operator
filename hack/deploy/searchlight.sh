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

set -eou pipefail

crds=(clusteralerts nodealerts podalerts incidents searchlightplugins)
apiversions=(v1alpha1.admission v1alpha1.incidents)

echo "checking kubeconfig context"
kubectl config current-context || {
    echo "Set a context (kubectl use-context <context>) out of the following:"
    echo
    kubectl config get-contexts
    exit 1
}
echo ""

OS=""
ARCH=""
DOWNLOAD_URL=""
DOWNLOAD_DIR=""
TEMP_DIRS=()
ONESSL=""
ONESSL_VERSION=v0.13.1

# http://redsymbol.net/articles/bash-exit-traps/
function cleanup() {
    rm -rf ca.crt ca.key server.crt server.key
    # remove temporary directories
    for dir in "${TEMP_DIRS[@]}"; do
        rm -rf "${dir}"
    done
}

# detect operating system
# ref: https://raw.githubusercontent.com/helm/helm/master/scripts/get
function detectOS() {
    OS=$(echo $(uname) | tr '[:upper:]' '[:lower:]')

    case "$OS" in
        # Minimalist GNU for Windows
        cygwin* | mingw* | msys*) OS='windows' ;;
    esac
}

# detect machine architecture
function detectArch() {
    ARCH=$(uname -m)
    case $ARCH in
        armv7*) ARCH="arm" ;;
        aarch64) ARCH="arm64" ;;
        x86) ARCH="386" ;;
        x86_64) ARCH="amd64" ;;
        i686) ARCH="386" ;;
        i386) ARCH="386" ;;
    esac
}

detectOS
detectArch

# download file pointed by DOWNLOAD_URL variable
# store download file to the directory pointed by DOWNLOAD_DIR variable
# you have to sent the output file name as argument. i.e. downloadFile myfile.tar.gz
function downloadFile() {
    if curl --output /dev/null --silent --head --fail "$DOWNLOAD_URL"; then
        curl -fsSL ${DOWNLOAD_URL} -o $DOWNLOAD_DIR/$1
    else
        echo "File does not exist"
        exit 1
    fi
}

export APPSCODE_ENV=${APPSCODE_ENV:-prod}
trap cleanup EXIT

# ref: https://github.com/appscodelabs/libbuild/blob/master/common/lib.sh#L55
inside_git_repo() {
    git rev-parse --is-inside-work-tree >/dev/null 2>&1
    inside_git=$?
    if [ "$inside_git" -ne 0 ]; then
        echo "Not inside a git repository"
        exit 1
    fi
}

detect_tag() {
    inside_git_repo

    # http://stackoverflow.com/a/1404862/3476121
    git_tag=$(git describe --exact-match --abbrev=0 2>/dev/null || echo '')

    commit_hash=$(git rev-parse --verify HEAD)
    git_branch=$(git rev-parse --abbrev-ref HEAD)
    commit_timestamp=$(git show -s --format=%ct)

    if [ "$git_tag" != '' ]; then
        TAG=$git_tag
        TAG_STRATEGY='git_tag'
    elif [ "$git_branch" != 'master' ] && [ "$git_branch" != 'HEAD' ] && [[ "$git_branch" != release-* ]]; then
        TAG=$git_branch
        TAG_STRATEGY='git_branch'
    else
        hash_ver=$(git describe --tags --always --dirty)
        TAG="${hash_ver}"
        TAG_STRATEGY='commit_hash'
    fi

    export TAG
    export TAG_STRATEGY
    export git_tag
    export git_branch
    export commit_hash
    export commit_timestamp
}

onessl_found() {
    # https://stackoverflow.com/a/677212/244009
    if [ -x "$(command -v onessl)" ]; then
        onessl version --check=">=${ONESSL_VERSION}" >/dev/null 2>&1 || {
            # old version of onessl found
            echo "Found outdated onessl"
            return 1
        }
        export ONESSL=onessl
        return 0
    fi
    return 1
}

# download onessl if it does not exist
onessl_found || {
    echo "Downloading onessl ..."

    ARTIFACT="https://github.com/kubepack/onessl/releases/download/${ONESSL_VERSION}"
    ONESSL_BIN=onessl-${OS}-${ARCH}
    case "$OS" in
        cygwin* | mingw* | msys*)
            ONESSL_BIN=${ONESSL_BIN}.exe
            ;;
    esac

    DOWNLOAD_URL=${ARTIFACT}/${ONESSL_BIN}
    DOWNLOAD_DIR="$(mktemp -dt onessl-XXXXXX)"
    TEMP_DIRS+=($DOWNLOAD_DIR) # store DOWNLOAD_DIR to cleanup later

    downloadFile $ONESSL_BIN # downloaded file name will be saved as the value of ONESSL_BIN variable

    export ONESSL=${DOWNLOAD_DIR}/${ONESSL_BIN}
    chmod +x $ONESSL
}

# ref: https://stackoverflow.com/a/7069755/244009
# ref: https://jonalmeida.com/posts/2013/05/26/different-ways-to-implement-flags-in-bash/
# ref: http://tldp.org/LDP/abs/html/comparison-ops.html

export SEARCHLIGHT_NAMESPACE=kube-system
export SEARCHLIGHT_SERVICE_ACCOUNT=searchlight-operator
export SEARCHLIGHT_ENABLE_RBAC=true
export SEARCHLIGHT_RUN_ON_MASTER=0
export SEARCHLIGHT_ICINGA_API_PASSWORD=
export SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK=false
export SEARCHLIGHT_DOCKER_REGISTRY=${DOCKER_REGISTRY:-appscode}
export SEARCHLIGHT_OPERATOR_TAG=8.0.0
export SEARCHLIGHT_ICINGA_TAG=8.0.0-k8s
export SEARCHLIGHT_IMAGE_PULL_SECRET=
export SEARCHLIGHT_IMAGE_PULL_POLICY=IfNotPresent
export SEARCHLIGHT_ENABLE_ANALYTICS=true
export SEARCHLIGHT_UNINSTALL=0
export SEARCHLIGHT_PURGE=0
export SEARCHLIGHT_ENABLE_STATUS_SUBRESOURCE=false
export SEARCHLIGHT_PRIORITY_CLASS=system-cluster-critical

export SCRIPT_LOCATION="curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/8.0.0/"
if [[ "$APPSCODE_ENV" == "dev" ]]; then
    detect_tag
    export SCRIPT_LOCATION="cat "
    export SEARCHLIGHT_OPERATOR_TAG=$TAG
    export SEARCHLIGHT_ICINGA_TAG=$TAG-k8s
    export SEARCHLIGHT_IMAGE_PULL_POLICY=Always
fi

KUBE_APISERVER_VERSION=$(kubectl version -o=json | $ONESSL jsonpath '{.serverVersion.gitVersion}')
$ONESSL semver --check='<1.9.0' $KUBE_APISERVER_VERSION || { export SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK=true; }
$ONESSL semver --check='<1.11.0' $KUBE_APISERVER_VERSION || { export SEARCHLIGHT_ENABLE_STATUS_SUBRESOURCE=true; }

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
    echo "    --enable-validating-webhook    enable/disable validating webhooks for Searchlight crds"
    echo "    --icinga-api-password          password used by icinga2 api (if unset, a random password will be generated and used)"
    echo "    --enable-status-subresource    If enabled, uses status sub resource for Searchlight crds"
    echo "    --enable-analytics             send usage events to Google Analytics (default: true)"
    echo "    --uninstall                    uninstall searchlight"
    echo "    --purge                        purges searchlight crd objects and crds"
}

while test $# -gt 0; do
    case "$1" in
        -h | --help)
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
            export SEARCHLIGHT_NAMESPACE=$(echo $1 | sed -e 's/^[^=]*=//g')
            shift
            ;;
        --docker-registry*)
            export SEARCHLIGHT_DOCKER_REGISTRY=$(echo $1 | sed -e 's/^[^=]*=//g')
            shift
            ;;
        --image-pull-secret*)
            secret=$(echo $1 | sed -e 's/^[^=]*=//g')
            export SEARCHLIGHT_IMAGE_PULL_SECRET="name: '$secret'"
            shift
            ;;
        --enable-validating-webhook*)
            val=$(echo $1 | sed -e 's/^[^=]*=//g')
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK=false
            fi
            shift
            ;;
        --enable-status-subresource*)
            val=$(echo $1 | sed -e 's/^[^=]*=//g')
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_ENABLE_STATUS_SUBRESOURCE=false
            fi
            shift
            ;;
        --enable-analytics*)
            val=$(echo $1 | sed -e 's/^[^=]*=//g')
            if [ "$val" = "false" ]; then
                export SEARCHLIGHT_ENABLE_ANALYTICS=false
            fi
            shift
            ;;
        --rbac*)
            val=$(echo $1 | sed -e 's/^[^=]*=//g')
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
        --icinga-api-password*)
            pass=$(echo $1 | sed -e 's/^[^=]*=//g')
            pass_b64=$(echo -n $pass | $ONESSL base64)
            export SEARCHLIGHT_ICINGA_API_PASSWORD="ICINGA_API_PASSWORD: '$pass_b64'"
            shift
            ;;
        --uninstall)
            export SEARCHLIGHT_UNINSTALL=1
            shift
            ;;
        --purge)
            export SEARCHLIGHT_PURGE=1
            shift
            ;;
        *)
            show_help
            exit 1
            ;;
    esac
done

if [ "$SEARCHLIGHT_NAMESPACE" != "kube-system" ]; then
    export SEARCHLIGHT_PRIORITY_CLASS=""
fi

if [ "$SEARCHLIGHT_UNINSTALL" -eq 1 ]; then
    # delete webhooks and apiservices
    kubectl delete validatingwebhookconfiguration -l app=searchlight || true
    kubectl delete mutatingwebhookconfiguration -l app=searchlight || true
    kubectl delete apiservice -l app=searchlight
    # delete searchlight operator
    kubectl delete deployment -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete service -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete secret -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    # delete RBAC objects, if --rbac flag was used.
    kubectl delete serviceaccount -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete clusterrolebindings -l app=searchlight
    kubectl delete clusterrole -l app=searchlight
    kubectl delete rolebindings -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE
    kubectl delete role -l app=searchlight --namespace $SEARCHLIGHT_NAMESPACE

    echo "waiting for searchlight operator pod to stop running"
    for (( ; ; )); do
        pods=($(kubectl get pods -n $SEARCHLIGHT_NAMESPACE -l app=searchlight -o jsonpath='{range .items[*]}{.metadata.name} {end}'))
        total=${#pods[*]}
        if [ $total -eq 0 ]; then
            break
        fi
        sleep 2
    done

    # https://github.com/kubernetes/kubernetes/issues/60538
    if [ "$SEARCHLIGHT_PURGE" -eq 1 ]; then
        for crd in "${crds[@]}"; do
            pairs=($(kubectl get ${crd}.monitoring.appscode.com --all-namespaces -o jsonpath='{range .items[*]}{.metadata.name} {.metadata.namespace} {end}' || true))
            total=${#pairs[*]}

            # save objects
            if [ $total -gt 0 ]; then
                echo "dumping ${crd} objects into ${crd}.yaml"
                kubectl get ${crd}.monitoring.appscode.com --all-namespaces -o yaml >${crd}.yaml
            fi

            for ((i = 0; i < $total; i += 2)); do
                name=${pairs[$i]}
                namespace=${pairs[$i + 1]}
                # delete crd object
                echo "deleting ${crd} $namespace/$name"
                kubectl delete ${crd}.monitoring.appscode.com $name -n $namespace
            done

            # delete crd
            kubectl delete crd ${crd}.monitoring.appscode.com || true
        done

        # delete user roles
        kubectl delete clusterroles appscode:searchlight:admin appscode:searchlight:edit appscode:searchlight:view
    fi

    echo
    echo "Successfully uninstalled Searchlight!"
    exit 0
fi

echo "checking whether extended apiserver feature is enabled"
$ONESSL has-keys configmap --namespace=kube-system --keys=requestheader-client-ca-file extension-apiserver-authentication || {
    echo "Set --requestheader-client-ca-file flag on Kubernetes apiserver"
    exit 1
}
echo ""

export KUBE_CA=
if [ "$SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK" = true ]; then
    $ONESSL get kube-ca >/dev/null 2>&1 || {
        echo "Admission webhooks can't be used when kube apiserver is accesible without verifying its TLS certificate (insecure-skip-tls-verify : true)."
        echo
        exit 1
    }
    export KUBE_CA=$($ONESSL get kube-ca | $ONESSL base64)
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

${SCRIPT_LOCATION}hack/deploy/operator.yaml | $ONESSL envsubst | kubectl apply -f -

if [ "$SEARCHLIGHT_ENABLE_RBAC" = true ]; then
    ${SCRIPT_LOCATION}hack/deploy/service-account.yaml | $ONESSL envsubst | kubectl apply -f -
    ${SCRIPT_LOCATION}hack/deploy/rbac-list.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
    ${SCRIPT_LOCATION}hack/deploy/user-roles.yaml | $ONESSL envsubst | kubectl auth reconcile -f -
fi

if [ "$SEARCHLIGHT_RUN_ON_MASTER" -eq 1 ]; then
    kubectl patch deploy searchlight-operator -n $SEARCHLIGHT_NAMESPACE \
        --patch="$(${SCRIPT_LOCATION}hack/deploy/run-on-master.yaml)"
fi

if [ "$SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK" = true ]; then
    ${SCRIPT_LOCATION}hack/deploy/validating-webhook.yaml | $ONESSL envsubst | kubectl apply -f -
fi

echo
echo "waiting until searchlight operator deployment is ready"
$ONESSL wait-until-ready deployment searchlight-operator --timeout=10m --namespace $SEARCHLIGHT_NAMESPACE || {
    echo "Searchlight operator deployment failed to be ready"
    exit 1
}

if [ "$SEARCHLIGHT_ENABLE_VALIDATING_WEBHOOK" = true ]; then
    echo "waiting until searchlight apiservice is available"
    for gv in "${apiversions[@]}"; do
        $ONESSL wait-until-ready apiservice ${gv}.monitoring.appscode.com || {
            echo "${gv}.monitoring.appscode.com apiservice failed to be ready"
            exit 1
        }
    done
fi

echo "waiting until searchlight crds are ready"
for crd in "${crds[@]}"; do
    $ONESSL wait-until-ready crd ${crd}.monitoring.appscode.com || {
        echo "$crd crd failed to be ready"
        exit 1
    }
done

echo
echo "Successfully installed Searchlight in $SEARCHLIGHT_NAMESPACE namespace!"
