# Install Hostfacts
Hostfacts is a http server used to expose various [node metrics](/pkg/hostfacts/server.go#L32). This is a wrapper around the wonderful [shirou/gopsutil](https://github.com/shirou/gopsutil) library. This is used by [`check_node_volume`](/docs/node-alerts/node_volume.md) and [`check_pod_volume`](/docs/pod-alerts/pod_volume.md) commands to detect available disk space. To use these check commands, hostfacts must be installed directly on every node in the cluster. Hostfacts can't be deployed using DaemonSet. This guide will walk you through how to deploy hostfacts as a Systemd service in minikube.

## Deploy Hostfacts
First ssh into a Kubernetes node. If you are using [Minikube](https://github.com/kubernetes/minikube), run the following command:
```console
$ minikube ssh
```

Now, download and install a pre-built binary using the following command:
```console
curl -Lo hostfacts https://cdn.appscode.com/binaries/hostfacts/3.0.0/hostfacts-linux-amd64 \
  && chmod +x hostfacts \
  && sudo mv hostfacts /usr/bin/
```










### Deploy Hostfacts

Write `hostfacts.service` file in __systemd directory__ in your kubernetes node.

##### systemd directory
* Ubuntu

    ```console
    /lib/systemd/system/hostfacts.service
    ```
* RedHat

    ```console
    /usr/lib/systemd/system/hostfacts.service
    ```


##### `hostfacts.service`

```ini
[Unit]
Description=Provide host facts

[Service]
ExecStart=/usr/bin/hostfacts
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

Set one of the following if you want to set authentication in `hostfacts`

* Basic Auth

    ```console
    # Use ENV
    # Add Environment in hostfacts.service under [Service] section
    Environment=HOSTFACTS_AUTH_USERNAME="<username>"
    Environment=HOSTFACTS_AUTH_PASSWORD="<password>"
    ```
    You can pass flags instead of using environment variables
    ```
    # Use Flags
    # Modify ExecStart in [Service] section
    ExecStart=/usr/bin/hostfacts --username="<username>" --password="<password>"
    ```
* Token

    ```console
    # Use ENV
    # Add Environment in hostfacts.service under [Service] section
    Environment=HOSTFACTS_AUTH_TOKEN="<token>"
    ```
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


### Add `hostfacts` binary

Download `hostfacts` and add binary in `/usr/bin`

```console
curl -fSsL  https://cdn.appscode.com/binaries/hostfacts/3.0.0/hostfacts-linux-amd64 -o /usr/bin/hostfacts

# Change access permissions for hostfacts binary
chmod +x /usr/bin/hostfacts
```

##### Start Service

```console
# Configure to be automatically started at boot time
systemctl enable hostfacts

# Start service
systemctl start hostfacts
```
