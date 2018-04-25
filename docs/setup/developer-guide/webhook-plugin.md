---
title: Searchlight Webhook Plugin | Icinga2
description: How to write Webhook Plugin for Searchlight 
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: write-searchlight-webhook-plugin
    name: SearchlightPlugin
    parent: developer-guide
    weight: 15
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

## Webhook Plugin

Command `check_webhook` calls a HTTP server with user provided variables and receives response to determine Service State.

In this tutorial, we will see how we can write a webhook for Searchlight.

The most important part for this webhook is its `Response` type.

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

## Next Steps

- Learn how to use this [webhook with Searchlight](/docs/guides/plugin/webhook-plugin.md).
- For some demo searchlight-plugin, visit [this](https://github.com/appscode/searchlight-plugin).