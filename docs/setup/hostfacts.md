---
title: Install Hostfacts
description: Install Hostfacts
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: hostfacts-searchlight
    name: Install Hostfacts
    parent: setup
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Hostfacts
[Hostfacts](/docs/reference/hostfacts/hostfacts_run.md) is a http server used to expose various [node metrics](https://github.com/appscode/searchlight/blob/29d4d2150116a284d74368931e6fdfe58efc7e6e/pkg/hostfacts/server.go#L32). This is a wrapper around the wonderful [shirou/gopsutil](https://github.com/shirou/gopsutil) library. This is used by [`check_node_volume`](/docs/guides/node-alerts/node-volume.md) and [`check_pod_volume`](/docs/guides/pod-alerts/pod-volume.md) commands to detect disk usage stats. To use these check commands, hostfacts must be installed directly on every node in the cluster. Hostfacts can't be deployed using DaemonSet. This guide will walk you through how to deploy hostfacts as a Systemd service.

## Installation Guide
First ssh into a Kubernetes node. If you are using [Minikube](https://github.com/kubernetes/minikube), run the following command:
```console
$ minikube ssh
```

### Install Hostfacts
Now, download and install a pre-built binary using the following command:
```console
curl -Lo hostfacts https://cdn.appscode.com/binaries/hostfacts/6.0.0-rc.0/hostfacts-linux-amd64 \
  && chmod +x hostfacts \
  && sudo mv hostfacts /usr/bin/
```

If you are using kube-up scripts to provision Kubernetes cluster, you can find a salt formula [here](https://github.com/appscode/kubernetes/tree/1.5.7-ac/cluster/saltbase/salt/appscode-hostfacts).


### Create Systemd Service
To run hostfacts server as a Systemd service, write `hostfacts.service` file in __systemd directory__ in your node.
```console
# Debian/Ubuntu (example, minikube)
$ sudo vi /lib/systemd/system/hostfacts.service

# RedHat
$ sudo vi /usr/lib/systemd/system/hostfacts.service
```

Hostfacts supports various types of authentication mechanism. Write the `hostfacts.service` accordingly.

#### Hostfacts without authentication
If you are running Kubernetes cluster inside a private network in AWS or GCP or just for testing in minikube, you may ignore authentication and SSL. In that case, use a `hostfacts.service` file like below:

```ini
[Unit]
Description=Provide host facts

[Service]
ExecStart=/usr/bin/hostfacts run --v=3
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Hostfacts with Basic auth
If you want to use a username/password pair with your Hostfacts binary, pass it via flag. Please note that, all nodes on your cluster must use the same username/password.

```ini
[Unit]
Description=Provide host facts

[Service]
ExecStart=/usr/bin/hostfacts run --v=3 --username="<username>" --password="<password>"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Hostfacts with Bearer token
If you want to use a bearer token with your Hostfacts binary, pass it via flag. Please note that, all nodes on your cluster must use the same token and ca certificate if any.

```ini
[Unit]
Description=Provide host facts

[Service]
ExecStart=/usr/bin/hostfacts run --v=3 --token="<token>"
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Using SSL
If your cluster is running inside a shared network (eg, DigitalOcean), you should enable SSL. If you want to set SSL certificate, do the following:

 - Generate a single pair of CA certificate and key. Then generate a separate SSL certificate pair for each node in your cluster. See the steps [here](/docs/setup/certificate.md).
 - Now, copy the ca.crt and node specific server.crt and server.key to the appropriate node. We recommend using folder `/srv/hostfacts/pki/`.
 - Use flags to pass the path to node specific certificates to its hostfact binary.

```ini
# Basic auth
ExecStart=/usr/bin/hostfacts run --v=3 --username="<username>" --password="<password>" --caCertFile="<path to ca cert file>" --certFile="<path to server cert file>" --keyFile="<path to server key file>"
```

```ini
# Bearer token
ExecStart=/usr/bin/hostfacts run --v=3 --token="<token>" --caCertFile="<path to ca cert file>" --certFile="<path to server cert file>" --keyFile="<path to server key file>"
```

### Activate Systemd service

```console
# Configure to be automatically started at boot time
$ sudo systemctl enable hostfacts

# Start service
$ sudo systemctl start hostfacts
```

### Create Hostfacts Secret
The last step is to create a Secret so that Searchlight operator can connect to Hostfacts server. This secret must be created in the same namespace where Searchlight operator is running.

| Key                    | Default | Description
|------------------------|---------|-------------|
| HOSTFACTS_PORT         | 56977   | `Required` Port used by hostfacts server. To change the default value, use `--address` flag |
| HOSTFACTS_USERNAME     |         | `Optional` Username for basic auth                                                          |
| HOSTFACTS_PASSWORD     |         | `Optional` Password for basic auth                                                          |
| HOSTFACTS_TOKEN        |         | `Optional` Token for bearer auth                                                            |
| HOSTFACTS_CA_CERT_DATA |         | `Optional` PEM encoded CA certificate used by Hostfacts server                              |

```console
$ echo -n '' > HOSTFACTS_PORT
$ echo -n '' > HOSTFACTS_USERNAME
$ echo -n '' > HOSTFACTS_PASSWORD
$ echo -n '' > HOSTFACTS_TOKEN
$ echo -n '' > HOSTFACTS_CA_CERT_DATA
$ kubectl create secret generic hostfacts -n kube-system \
    --from-file=./HOSTFACTS_PORT \
    --from-file=./HOSTFACTS_USERNAME \
    --from-file=./HOSTFACTS_PASSWORD \
    --from-file=./HOSTFACTS_TOKEN \
    --from-file=./HOSTFACTS_CA_CERT_DATA
secret "hostfacts" created
```
