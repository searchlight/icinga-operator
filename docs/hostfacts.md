# Hostfacts
[Hostfacts](/docs/reference/hostfacts/hostfacts_run.md) is a http server used to expose various [node metrics](/pkg/hostfacts/server.go#L32). This is a wrapper around the wonderful [shirou/gopsutil](https://github.com/shirou/gopsutil) library. This is used by [`check_node_volume`](/docs/node-alerts/node_volume.md) and [`check_pod_volume`](/docs/pod-alerts/pod_volume.md) commands to detect available disk space. To use these check commands, hostfacts must be installed directly on every node in the cluster. Hostfacts can't be deployed using DaemonSet. This guide will walk you through how to deploy hostfacts as a Systemd service.

## Installation Guide
First ssh into a Kubernetes node. If you are using [Minikube](https://github.com/kubernetes/minikube), run the following command:
```console
$ minikube ssh
```

### Install Hostfacts
Now, download and install a pre-built binary using the following command:
```console
curl -Lo hostfacts https://cdn.appscode.com/binaries/hostfacts/3.0.0/hostfacts-linux-amd64 \
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

 - Generate a single pair of CA certificate and key. Then generate a separate SSL certificate pair for each node in your cluster. See the steps [here](/docs/certificate.md).
 - Now, copy the ca.crt and node specific server.crt and server.key to the appropriate node. We recommend using folder `/srv/hostfacts/pki/`.
 - Use flags to pass the path to node specific certificates to its hostfact binary.

```ini
# Basic auth
ExecStart=/usr/bin/hostfacts run --v=3 --username="<username>" --password="<password>" --caCertFile="<path to ca cert file>" --certFile="<path to server cert file>" --keyFile="<path to server key file>"


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
