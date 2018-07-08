#!/bin/bash

set -x -e

# start docker and log-in to docker-hub
entrypoint.sh
docker login --username=$DOCKER_USER --password=$DOCKER_PASS
docker run hello-world

# install python pip
apt-get update >/dev/null
apt-get install -y python python-pip gsutil git libwww-perl >/dev/null

# install kubectl
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl &>/dev/null
chmod +x ./kubectl
mv ./kubectl /bin/kubectl

# install onessl
curl -fsSL -o onessl https://github.com/kubepack/onessl/releases/download/0.3.0/onessl-linux-amd64 &&
  chmod +x onessl &&
  mv onessl /usr/local/bin/

# install pharmer
pushd /tmp
curl -LO https://cdn.appscode.com/binaries/pharmer/0.1.0-rc.4/pharmer-linux-amd64
chmod +x pharmer-linux-amd64
mv pharmer-linux-amd64 /bin/pharmer
popd

function cleanup() {
  # delete cluster on exit
  pharmer get cluster || true
  pharmer delete cluster $NAME || true
  pharmer get cluster || true
  sleep 300 || true
  pharmer apply $NAME || true
  pharmer get cluster || true

  # delete docker image on exit
  curl -LO https://raw.githubusercontent.com/appscodelabs/libbuild/master/docker.py || true
  chmod +x docker.py || true
  ./docker.py del_tag appscodeci searchlight $SEARCHLIGHT_OPERATOR_TAG || true
  ./docker.py del_tag appscodeci icinga $SEARCHLIGHT_ICINGA_TAG || true
}
trap cleanup EXIT

# copy searchlight to $GOPATH
mkdir -p $GOPATH/src/github.com/appscode
cp -r searchlight $GOPATH/src/github.com/appscode
pushd $GOPATH/src/github.com/appscode/searchlight

# name of the cluster
# nameing is based on repo+commit_hash
NAME=searchlight-$(git rev-parse --short HEAD)

./hack/builddeps.sh
go get -u golang.org/x/tools/cmd/goimports
go get github.com/Masterminds/glide
go get github.com/sgotti/glide-vc
go get github.com/onsi/ginkgo/ginkgo
go install github.com/onsi/ginkgo/ginkgo

export APPSCODE_ENV=dev
export DOCKER_REGISTRY=appscodeci

./hack/make.py build searchlight
./hack/make.py build hyperalert

./hack/docker/icinga/alpine/build.sh
./hack/docker/icinga/alpine/build.sh push

./hack/docker/searchlight/setup.sh
./hack/docker/searchlight/setup.sh push

popd

#create credential file for pharmer
cat >cred.json <<EOF
{
    "token" : "$TOKEN"
}
EOF

# create cluster using pharmer
# note: make sure the zone supports volumes, not all regions support that
pharmer create credential --from-file=cred.json --provider=DigitalOcean cred
pharmer create cluster $NAME --provider=digitalocean --zone=nyc1 --nodes=2gb=1 --credential-uid=cred --kubernetes-version=v1.10.0
pharmer apply $NAME
pharmer use cluster $NAME
#wait for cluster to be ready
sleep 300
kubectl get nodes

# create storageclass
cat >sc.yaml <<EOF
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: standard
parameters:
  zone: nyc1
provisioner: external/pharmer
EOF

# create storage-class
kubectl create -f sc.yaml
sleep 120
kubectl get storageclass

pushd $GOPATH/src/github.com/appscode/searchlight

# run tests
source ./hack/deploy/searchlight.sh --docker-registry=appscodeci --enable-validating-webhook=true --rbac=true --icinga-api-password=1234
./hack/make.py test e2e --searchlight-service=searchlight-operator@kube-system --provider=do
