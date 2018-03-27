---
title: Uninstall
description: Searchlight Uninstall
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: uninstall-searchlight
    name: Uninstall
    parent: setup
    weight: 25
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: setup
---

> New to Searchlight? Please start [here](/docs/concepts/README.md).

# Uninstall Searchlight

To uninstall Searchlight operator, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/appscode/searchlight/6.0.0-rc.0/hack/deploy/searchlight.sh \
    | bash -s -- --uninstall [--namespace=NAMESPACE]

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

The above command will leave the Searchlight crd objects as-is. If you wish to **nuke** all Searchlight crd objects, also pass the `--purge` flag. This will keep a copy of Searchlight crd objects in your current directory.
