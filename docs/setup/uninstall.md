---
title: Uninstall
description: Searchlight Uninstall
menu:
  product_searchlight_5.0.0:
    identifier: uninstall-searchlight
    name: Uninstall
    parent: setup
    weight: 25
product_name: searchlight
menu_name: product_searchlight_5.0.0
section_menu_id: setup
---


> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Uninstall Searchlight
Please follow the steps below to uninstall Searchlight:

1. Delete the various objects created for Searchlight operator.
```console
$ ./hack/deploy/uninstall.sh 
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

2. Now, wait several seconds for Searchlight to stop running. To confirm that Searchlight operator pod(s) have stopped running, run:
```console
$ kubectl get pods --all-namespaces -l app=searchlight
```
