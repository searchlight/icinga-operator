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
To run hostfacts server as a System service, write `hostfacts.service` file in __systemd directory__ in your node.
```console
# Ubuntu (example, minikube)
$ sudo vi /lib/systemd/system/hostfacts.service

# RedHat
$ sudo vi /usr/lib/systemd/system/hostfacts.service
```




Below is a systemd service file for Hostfacts without any authentication.
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

Set one of the following if you want to set authentication in `hostfacts`

* Basic Auth

    You can pass flags instead of using environment variables
    ```
    # Use Flags
    # Modify ExecStart in [Service] section
    ExecStart=/usr/bin/hostfacts --username="<username>" --password="<password>"
    ```
* Token

    You can pass flag instead of using environment variable
    ```
    # Use Flags
    # Modify ExecStart in [Service] section
    ExecStart=/usr/bin/hostfacts --token="<token>"
    ```

If you want to set SSL certificate, do following

1. Generate certificates and key. See process [here](../icinga2/certificate.md).
2. Use flags to pass file directory

    ```console
    # Modify ExecStart in [Service] section
    ExecStart=/usr/bin/hostfacts --caCertFile="<path to ca cert file>" --certFile="<path to server cert file>" --keyFile="<path to server key file>"
    ```

You can ignore SSL when Kubernetes is running in private network like GCE, AWS.

> __Note:__ Modify `ExecStart` in `hostfacts.service`


### Activate Systemd service

```console
# Configure to be automatically started at boot time
$ sudo systemctl enable hostfacts

# Start service
$ sudo systemctl start hostfacts
```
