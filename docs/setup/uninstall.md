---
title: Uninstall
description: Searchlight Uninstall
menu:
  product_searchlight_5.1.0:
    identifier: uninstall-searchlight
    name: Uninstall
    parent: setup
    weight: 25
product_name: searchlight
menu_name: product_searchlight_5.1.0
section_menu_id: setup
---


> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Uninstall Searchlight
Please follow the steps below to uninstall Searchlight:

- Delete the various objects created for Searchlight operator.

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/5.1.0/hack/deploy/uninstall.sh | bash

+ kubectl delete deployment -l app=searchlight -n kube-system
deployment "searchlight-operator" deleted
+ kubectl delete service -l app=searchlight -n kube-system
service "searchlight-operator" deleted
+ kubectl delete secret -l app=searchlight -n kube-system
secret "searchlight-operator" deleted
+ kubectl delete serviceaccount -l app=searchlight -n kube-system
No resources found
+ kubectl delete clusterrolebindings -l app=searchlight -n kube-system
No resources found
+ kubectl delete clusterrole -l app=searchlight -n kube-system
No resources found
```

- Now, wait several seconds for Searchlight to stop running. To confirm that Searchlight operator pod(s) have stopped running, run:

```console
$ kubectl get pods --all-namespaces -l app=searchlight
```
