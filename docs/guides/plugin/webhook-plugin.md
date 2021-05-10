---
title: Webhook SearchlightPlugin
menu:
  product_searchlight_{{ .version }}:
    identifier: guides-webhook-searchlight-plugin
    name: Webhook SearchlightPlugin
    parent: searchlight-plugin
    weight: 20
product_name: searchlight
menu_name: product_searchlight_{{ .version }}
section_menu_id: guides
---

> New to SearchlightPlugin? Please start [here](/docs/setup/developer-guide/webhook-plugin.md).

# Custom Webhook Check Command

Since 8.0.0 release, Searchlight supports adding custom check commands using http webhook server. No longer you have to build binary and attach it inside Icinga container. Simply you can write a HTTP server and register your check command with Searchlight using `SearchlightPlugin` CRD.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: SearchlightPlugin
metadata:
  name: check-pod-count
spec:
  webhook:
    namespace: default
    name: searchlight-plugin
  alertKinds:
  - ClusterAlert
  arguments:
    vars:
      fields:
        warning:
          type: integer
        critical:
          type: integer
  states:
  - OK
  - Critical
  - Unknown
```

The `.spec` section has following parts:

## Spec

The `.spec` section determines how the webhook server will be used by Searchlight.

**spec.command**

`spec.command` defines the check command which will be called by Icinga. To use the webhook server, skip the command field.

**spec.webhook**

`spec.webhook` provides information of Kubernetes `Service` for the webhook server.

- `spec.webhook.namespace` represents the namespace of Service.
- `spec.webhook.name` represents the name of Service.

**spec.alertKinds**

`spec.alertKinds` is a required field that specifies which kinds of alerts will support this `CheckCommand`.
Possible values are: ClusterAlert, NodeAlert and PodAlert.

**spec.arguments**

`spec.arguments` defines arguments which will be passed to check command by Icinga.

- `spec.arguments.vars` defines user-defined arguments. These arguments can be provided to create alerts.

    - `spec.arguments.vars.fields` is the required field which provides the list of arguments with their `description` and `type`. Here,

          arguments:
            vars:
              fields:
                warning:
                  type: integer
                critical:
                  type: integer

        `warning` and `critical` are registered as user-defined variables. User can provide values for these variables while creating alerts.

        - `spec.arguments.vars.fields[].type` is required field used to define variable's data type. Possible values are `integer`, `number`, `boolean`, `string`, `duration`.
        > Note: duration will be converted to millisecond and will be passed to webhook as Int64. Example: 5m -> 300000. 
        - `spec.arguments.vars.fields[].description` describes the variable.

    - `spec.arguments.vars.required` represents the list of user-defined arguments required to create Alert. If any of these required arguments is not provided, Searchlight will give validation error.

    From above example, none of these variables are required. To make variable `critical` required, need to add following

              arguments:
                vars:
                  required:
                  - critical

- `spec.arguments.host` represents the list of Icinga host variables which will be passed to check command. [Here is the list](https://www.icinga.com/docs/icinga2/latest/doc/03-monitoring-basics/#host-runtime-macros) of available host variables.

Suppose, you need `host.check_attempt` to be forwarded to check command, you can add like this

```yaml
spec:
  arguments:
    host:
      attempt: check_attempt
```

Here,

- Icinga host variable `check_attempt` will be forwarded to webhook as variable `attempt`.

> Note: User can't provided value for these variables.

**spec.states**

`spec.state` are the supported states for this command. Different notification receivers can be set for each state.

Lets create above SearchlightPlugin

```console
$ kubectl apply -f ./docs/examples/plugins/webhook/demo-0.yaml
searchlightplugin "check-pod-count" created
```

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/add-plugin.svg" width="581" height="362">
</p>

CheckCommand `check-pod-count` is added in Icinga2 configuration. Here, `vars.fields` from `spec.arguments` are added as arguments in CheckCommand.


> Note: Webhook will be called with URL formatted as bellow:

`http://<spec.webhook.name>.<spec.webhook.namespace>.svc/<metadata.name>`

### Create ClusterAlert

Lets create a ClusterAlert for this CheckCommand.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: ClusterAlert
metadata:
  name: count-all-pods-demo-0
  namespace: demo
spec:
  check: check-pod-count
  vars:
    warning: 10
    critical: 15
  notifierSecretName: notifier-config
  receivers:
  - notifier: Mailgun
    state: Critical
    to: ["ops@example.com"]
```

Here,

- `spec.check` is the name of your custom check you added as SearchlightPlugin
- `spec.vars` are variables those are registered when SearchlightPlugin is created with `spec.arguments.vars`

```console
$ kubectl apply -f ./docs/examples/cluster-alerts/count-all-pods/demo-0.yaml
clusteralert "count-all-pods-demo-0" created
```

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/add-alert.svg" width="581" height="362">
</p>

Now periodically, Icinga will call `check_webhook` plugin under `hyperalert`.
And this plugin will call your webhook you have registered in your SearchlightPlugin. According to the response from webhook, Service State will be determined.

<p align="center">
  <img alt="lifecycle"  src="/docs/images/plugin/call-webhook.svg" width="581" height="362">
</p>

In the example above, Service State will be **Warning**.


## Writing Webhook Server

Command `check_webhook` calls a HTTP server with user provided variables and receives response to determine service state. In this tutorial, we will see how we can write a webhook server for Searchlight. The most important part for this webhook is its `Response` type.

```go
package main

type State int32

const (
  OK       State = iota // 0
  Warning               // 1
  Critical              // 2
  Unknown               // 3
)

type Response struct {
  Code    State  `json:"code"`
  Message string `json:"message,omitempty"`
}
```

Icinga2 Service State is determined according to `Code` in `Response`.

> Note: Webhook may not have any `Request` option.

Add HTTP handler to serve request.

```go
package main

import (
  "fmt"
  "net/http"
)

func main() {
  http.HandleFunc("/check-pod-count", func(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodPost:
      fmt.Println("do your stuff")
      fmt.Println("write response with code")
    default:
      http.Error(w, "", http.StatusNotImplemented)
      return
    }
  })
  http.ListenAndServe(":80", http.DefaultServeMux)
}
```

Here,

- Path `/check-pod-count` only serves POST request. And return Response according to its check.

> Note: This webhook should listen on `80` port and serve POST request.


Now build this server code.

```bash
go build -o webhook main.go
```

And build docker image and push to your registry.

```dockerfile
FROM ubuntu

RUN set -x \
  && apt-get update
RUN set -x \
  && apt-get install -y ca-certificates

COPY webhook /usr/bin/webhook

ENTRYPOINT ["webhook"]
EXPOSE 80
```

Now deploy this server in Kubernetes

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: searchlight-plugin
  labels:
    app: searchlight-plugin
spec:
  replicas: 1
  selector:
    matchLabels:
      app: searchlight-plugin
  template:
    metadata:
      labels:
        app: searchlight-plugin
    spec:
      containers:
      - name: webhook
        image: appscode/searchlight-plugin-go
        imagePullPolicy: Always
        ports:
        - containerPort: 80
---
apiVersion: v1
kind: Service
metadata:
  name: searchlight-plugin
  labels:
    app: searchlight-plugin
spec:
  ports:
  - name: http
    port: 80
    targetPort: 80
  selector:
    app: searchlight-plugin
```

Here,

- Service `searchlight-plugin` in Namespace `default` will be used in SearchlightPlugin

For a fully working version of this tutorial, visit [this](https://go.searchlight.dev/icinga-operator-plugin). Please note that the webhook server can be written in any programming language your prefer as long as the `Response` json format is maintained.
