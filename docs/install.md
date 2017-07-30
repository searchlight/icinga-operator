> New to Searchlight? Please start [here](/docs/tutorials/README.md).

# Installation Guide

## Using YAML
[![Install Searchlight](https://img.youtube.com/vi/Po4yXrQuHtQ/0.jpg)](https://www.youtube-nocookie.com/embed/Po4yXrQuHtQ)

Searchlight can be installed using YAML files includes in the [/hack/deploy](/hack/deploy) folder.

```console
# Install without RBAC roles
$ curl https://raw.githubusercontent.com/appscode/searchlight/3.0.0/hack/deploy/without-rbac.yaml \
  | kubectl apply -f -


# Install with RBAC roles
$ curl https://raw.githubusercontent.com/appscode/searchlight/3.0.0/hack/deploy/with-rbac.yaml \
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


## Accesing IcingaWeb2
Icinga comes with its own web dashboard called IcingaWeb. You can access IcingaWeb on your workstation by forwarding port `60006` of Searchlight operator pod.

```console
$ kubectl get pods --all-namespaces -l app=searchlight
NAME                                    READY     STATUS    RESTARTS   AGE
searchlight-operator-1987091405-ghj5b   3/3       Running   0          1m

$ kubectl port-forward searchlight-operator-1987091405-ghj5b -n kube-system 60006
Forwarding from 127.0.0.1:60006 -> 60006
E0728 04:07:28.237822   10898 portforward.go:212] Unable to create listener: Error listen tcp6 [::1]:60006: bind: cannot assign requested address
Handling connection for 60006
Handling connection for 60006
^C⏎
```

Now, open http://127.0.0.1:60006 URL on your broswer. To login, use username `admin` and password `changeit`. If you want to change the password, read the next section.


## Configuring Icinga
Searchlight installation scripts above creates a Secret called `searchlight-operator` to store icinga configuration. This following keys are supported in this Secret.

| Key                    | Default Value  | Description                                             |
|------------------------|----------------|---------------------------------------------------------|
| ICINGA_WEB_UI_PASSWORD | _**changeit**_ | Password of `admin` user in IcingaWeb2                  |
| ICINGA_API_PASSWORD    | auto-generated | Password of icinga api user `icingaapi`                 |
| ICINGA_CA_CERT         | auto-generated | PEM encoded CA certificate used for icinga api endpoint |
| ICINGA_SERVER_CERT     | auto-generated | PEM encoded certificate used for icinga api endpoint    |
| ICINGA_SERVER_KEY      | auto-generated | PEM encoded private key used for icinga api endpoint    |
| ICINGA_IDO_PASSWORD    | auto-generated | Password of postgres user `icingaido`                   |
| ICINGA_WEB_PASSWORD    | auto-generated | Password of postgres user `icingaweb`                   |

To change the `admin` user login password in IcingaWeb, change the value of `ICINGA_WEB_UI_PASSWORD` key in Secret `searchlight-operator` and restart Searchlight operator pod(s).

```console
$ kubectl edit secret searchlight-operator -n kube-system
# Update the value of ICINGA_WEB_UI_PASSWORD key

$ kubectl get pods --all-namespaces -l app=searchlight
NAME                                    READY     STATUS    RESTARTS   AGE
searchlight-operator-1987091405-ghj5b   3/3       Running   0          1m

$ kubectl delete pods -n kube-system searchlight-operator-1987091405-ghj5b
pod "searchlight-operator-1987091405-ghj5b" deleted
```
