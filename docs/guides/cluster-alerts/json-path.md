---
title: JSON Path
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-json-path
    name: JSON Path
    parent: cluster-alert
    weight: 40
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Check json-path

Check command `json-path` is used to check JSON HTTP response using [jsonpath](https://kubernetes.io/docs/reference/kubectl/jsonpath/) queries.

## Spec
`json-path` check command has the following variables:

- `url` - URL to get data
- `secretName` - Name of Kubernetes Secret used to call HTTP api.
- `warning` - Query for warning which returns [true/false].
- `critical` - Query for critical which returns [true/false].

### Query

A query used in `warning` and `critical` variable must return boolean [true/false].
In this query, you can use following operators:

* Modifiers: `+` `-` `/` `*` `&` `|` `^` `**` `%` `>>` `<<`
* Comparators: `>` `>=` `<` `<=` `==` `!=` `=~` `!~`
* Logical ops: `||` `&&`

And also you can use [jsonpath](https://kubernetes.io/docs/reference/kubectl/jsonpath/) queries to get values from JSON data.

#### Examples

Lets assume, we get following JSON from provided URL.
```json
{
   "Book":[
      {
         "Category":"reference",
         "Author":"Nigel Rees",
         "Title":"Sayings of the Centurey",
         "Price":8.95
      }
   ],
   "Bicycle":[
      {
         "Color":"red",
         "Price":19.95,
         "IsNew":true
      },
      {
         "Color":"green",
         "Price":20.01,
         "IsNew":false
      }
   ]
}
```

Supported queries look like:

* `{.Book[0].Category}==novel`
* `{.Book[0].Category}!=reference`
* `{.Book[0].Price} > 10`
* `{.Book[0].Price} < 10`
* `{.Bicycle[0].IsNew} != true`
* `{.Bicycle[0].Color} == red && {.Bicycle[0].Price} < 20`
* `{.Bicycle[0].Color} != {.Bicycle[1].Color}`

and so many others.

### Secret

The following keys are supported for Secret passed via `secretName` flag.

| Key                    | Description                                                 |
-------------------------|-------------------------------------------------------------|
| `USERNAME`             | `Optional` Username used with Basic auth for HTTP URL.      |
| `PASSWORD`             | `Optional` Password used with Basic auth for HTTP URL.      |
| `TOKEN`                | `Optional` Token used as Bearer auth for HTTP URL.          |
| `CA_CERT_DATA`         | `Optional` PEM encoded CA certificate used by HTTP URL.     |
| `CLIENT_CERT_DATA`     | `Optional` PEM encoded Client certificate used by HTTP URL. |
| `CLIENT_KEY_DATA`      | `Optional` PEM encoded Client private key used by HTTP URL. |
| `INSECURE_SKIP_VERIFY` | `Optional` If set to `true`, skip certificate verification. |

Execution of this command can result in following states:

- OK
- Warning
- Critical
- Unknown


## Tutorial

### Before You Begin
At first, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube).

Now, install Searchlight operator in your cluster following the steps [here](/docs/setup/install.md).

To keep things isolated, this tutorial uses a separate namespace called `demo` throughout this tutorial. Run the following command to prepare your cluster for this tutorial:

```console
$ kubectl create namespace demo
namespace "demo" created

$ kubectl get namespaces
NAME          STATUS    AGE
default       Active    6h
kube-public   Active    6h
kube-system   Active    6h
demo          Active    4m
```


### Check JSON response of HTTP api
In this tutorial, a ClusterAlert will be used check JSON response of a HTTP api.

```yaml
$ cat ./docs/examples/cluster-alerts/json-path/demo-0.yaml

apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: json-path-demo-0
  namespace: demo
spec:
  check: json-path
  vars:
    url: http://echo.jsontest.com/key/value/one/two
    critical: '{.one} != "one"'
  checkInterval: 30s
  alertInterval: 2m
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```

```console
$ kubectl apply -f ./docs/examples/cluster-alerts/json-path/demo-0.yaml
clusteralert "json-path-demo-0" created

$ kubectl describe clusteralert -n demo json-path-demo-0
Name:		json-path-demo-0
Namespace:	demo
Labels:		<none>
Events:
  FirstSeen	LastSeen	Count	From			SubObjectPath	Type		Reason		Message
  ---------	--------	-----	----			-------------	--------	------		-------
  16s		16s		1	Searchlight operator			Normal		SuccessfulSync	Applied ClusterAlert: "json-path-demo-0"
```

Voila! `json-path` command has been synced to Icinga2. Please visit [here](/docs/guides/notifiers.md) to learn how to configure notifier secret. Now, open IcingaWeb2 in your browser. You should see a Icinga host `demo@cluster` and Icinga service `json-path-demo-0`.

![check-all-pods](/docs/images/cluster-alerts/json-path/demo-0.png)


### Cleaning up
To cleanup the Kubernetes resources created by this tutorial, run:
```console
$ kubectl delete ns demo
```

If you would like to uninstall Searchlight operator, please follow the steps [here](/docs/setup/uninstall.md).


## Next Steps
 - To periodically run various checks on nodes in a Kubernetes cluster, use [NodeAlerts](/docs/concepts/alert-types/node-alert.md).
 - To periodically run various checks on pods in a Kubernetes cluster, use [PodAlerts](/docs/concepts/alert-types/pod-alert.md).
 - See the list of supported notifiers [here](/docs/guides/notifiers.md).
 - Wondering what features are coming next? Please visit [here](/docs/roadmap.md).
 - Want to hack on Searchlight? Check our [contribution guidelines](/docs/CONTRIBUTING.md).
