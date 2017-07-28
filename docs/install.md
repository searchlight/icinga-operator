> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Installation Guide

[![Install Searchlight](https://img.youtube.com/vi/Po4yXrQuHtQ/0.jpg)](https://www.youtube-nocookie.com/embed/Po4yXrQuHtQ)

## Create Cluster Config
Before you can install Searchlight, you need a cluster config for Searchlight. Cluster config is defined in YAML format. You find an example config in [./hack/deploy/config.yaml](/hack/deploy/config.yaml).

```yaml
$ cat ./hack/deploy/config.yaml

apiServer:
  address: :8080
  enableReverseIndex: true
  enableSearchIndex: true
enableConfigSyncer: true
eventForwarder:
  ingressAdded:
    handle: true
  nodeAdded:
    handle: true
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
  storageAdded:
    handle: true
  warningEvents:
    handle: true
    namespaces:
    - kube-system
janitors:
- elasticsearch:
    endpoint: http://elasticsearch-logging.kube-system:9200
    logIndexPrefix: logstash-
  kind: Elasticsearch
  ttl: 2160h0m0s
- influxdb:
    endpoint: https://monitoring-influxdb.kube-system:8086
  kind: InfluxDB
  ttl: 2160h0m0s
notifierSecretName: searchlight-notifier
recycleBin:
  handleUpdates: false
  path: /tmp/searchlight/trash
  receiver:
    notifier: mailgun
    to:
    - ops@example.com
  ttl: 168h0m0s
snapshotter:
  Storage:
    gcs:
      bucket: restic
      prefix: minikube
    storageSecretName: snap-secret
  sanitize: true
  schedule: '@every 6h'
```

To understand the various configuration options, check Searchlight [tutorials](/docs/tutorials/README.md). Once you are satisfied with the configuration, create a Secret with the Searchlight cluster config under `config.yaml` key.

You may have to create another [Secret for notifiers](/docs/tutorials/notifiers.md). If you are [storing cluster snapshots](/docs/tutorials/cluster-snapshot.md) in cloud storage, you have to create a Secret appropriately.

### Generate Config using script
If you are familiar with GO, you can use the [./hack/config/main.go](/hack/config/main.go) script to generate a cluster config. Open this file in your favorite editor, update the config returned from `#CreateClusterConfig()` method. Then run the script to generate updated config in [./hack/deploy/config.yaml](/hack/deploy/config.yaml).

```console
go run ./hack/config/main.go
```

### Verifying Cluster Config
Searchlight includes a check command to verify a cluster config. Download the pre-built binary from [appscode/searchlight Github releases](https://github.com/appscode/searchlight/releases) and put the binary to some directory in your `PATH`.

```console
$ searchlight check --clusterconfig=./hack/deploy/config.yaml
Cluster config was parsed successfully.
```

## Using YAML
Searchlight can be installed using YAML files includes in the [/hack/deploy](/hack/deploy) folder.

```console
# Install without RBAC roles
$ curl https://raw.githubusercontent.com/appscode/searchlight/0.1.0/hack/deploy/without-rbac.yaml \
  | kubectl apply -f -


# Install with RBAC roles
$ curl https://raw.githubusercontent.com/appscode/searchlight/0.1.0/hack/deploy/with-rbac.yaml \
  | kubectl apply -f -
```

## Using Helm
Searchlight can be installed via [Helm](https://helm.sh/) using the [chart](/chart/searchlight) included in this repository. To install the chart with the release name `my-release`:
```bash
$ helm install chart/searchlight --name my-release
```
To see the detailed configuration options, visit [here](/chart/searchlight/README.md).


## Verify installation
To check if Searchlight operator pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l app=searchlight --watch
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.


## Update Cluster Config
If you would like to update cluster config, update the `searchlight-config` Secret and restart Searchlight operator pod(s).










































$ kubectl apply -f ./hack/deploy/without-rbac.yaml
secret "searchlight-operator" created
deployment "searchlight-operator" created
service "searchlight-operator" created

$ kubectl get pods -n kube-system -w
NAME                          READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube   1/1       Running   0          14m
kube-dns-1301475494-p8pcr   3/3       Running   0         14m
kubernetes-dashboard-psl27   1/1       Running   0         14m
searchlight-operator-1987091405-ghj5b   0/3       ContainerCreating   0         5s
searchlight-operator-1987091405-ghj5b   3/3       Running   0         43s
^C⏎

$ kubectl get pods -n kube-system
NAME                                    READY     STATUS    RESTARTS   AGE
kube-addon-manager-minikube             1/1       Running   0          15m
kube-dns-1301475494-p8pcr               3/3       Running   0          15m
kubernetes-dashboard-psl27              1/1       Running   0          15m
searchlight-operator-1987091405-ghj5b   3/3       Running   0          1m

$ kubectl port-forward searchlight-operator-1987091405-ghj5b -n kube-system 60006
Forwarding from 127.0.0.1:60006 -> 60006
E0728 04:07:28.237822   10898 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:60006: bind: cannot assign requested address
Handling connection for 60006
Handling connection for 60006
^C⏎

Open http://127.0.0.1:60006
