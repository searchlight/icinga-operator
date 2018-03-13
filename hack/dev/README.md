```
$ kubectl apply -f ./hack/dev/icinga.yaml

$ kubectl get pods,svc -n demo
$ minikube service list -n demo
```

- Open dashboard for icinga (port 80), login with user: `admin` and password: `changeit`

- Find out the node port for icinga api (port: 5665)
- Update the testconfig/config.ini to set `ICINGA_ADDRESS` to `$(minikube ip):(nodeport-for-5665)`

```
echo -n 'your-mailgun-api-key' > MAILGUN_API_KEY
kubectl create secret generic notifier-config -n demo \
    --from-file=./hack/dev/notifier/MAILGUN_DOMAIN \
    --from-file=./hack/dev/notifier/MAILGUN_FROM \
    --from-file=./hack/dev/notifier/MAILGUN_API_KEY
```

```
kubectl run nginx -n demo --image=nginx --labels="app=nginx"
kubectl scale deploy nginx -n demo --replicas=2
```

### Quickly update hyperalert

```
$ ./hack/make.py build hyperalert

$ eval $(minikube docker-env)
$ docker cp ./dist/hyperalert/hyperalert-alpine-amd64 e98c464a619b:/usr/lib/monitoring-plugins/hyperalert
```


### Directly call notfier inside icinga container

```
docker exec -it <icinga-container-id> sh

/usr/lib/monitoring-plugins/hyperalert notifier --alert=pod-exists-demo-0 --type=Problem --state=Critical --host=demo@cluster --output="test" --time="2006-01-02 15:04:05 +0000" -v=10
```

### Delete test incident
```
$ kubectl delete incident -n demo cluster.pod-exists-demo-0.20060102-1504
```

## direct EAS
curl -k -vv https://127.0.0.1:8443/apis/incidents.monitoring.appscode.com/v1alpha1

## minikube apiserver
curl -k -vv https://192.168.99.100:8443/apis --cert $HOME/.minikube/client.crt --key $HOME/.minikube/client.key


## inside minikube
curl -k -vv https://10.0.2.2:8443/apis/incidents.monitoring.appscode.com/v1alpha1

```
$ kubectl create -f ./hack/dev/apiregistration.yaml
$ kubectl get apiservice v1alpha1.incidents.monitoring.appscode.com -o yaml
```
```
$ kubectl create -f ./hack/dev/ack.yaml
```