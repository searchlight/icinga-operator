---
title: Acknowledgement Concepts
description: Acknowledgement Concepts
menu:
  product_searchlight_7.0.0-rc.0:
    identifier: acknowledgement-concepts
    parent: incident
    name: Acknowledgement Concepts
    weight: 15
menu_name: product_searchlight_7.0.0-rc.0
---

# Acknowledgement

Kubernetes Extended Api Server resource **Acknowledgement** is used to acknowledge an Incident with type **Problem**. 

Following is the example of Acknowledgement object

```yaml
apiVersion: incidents.monitoring.appscode.com
kind: Acknowledgement
metadata:
  name: cluster.pod-exists-demo-0.20180428-1109
  namespace: demo
request:
    comment: working on fix
```

When user creates this Acknowledgement object, Searchlight operator gets Incident with same name of Acknowledgement object.
Operator then acknowledges Icinga notification with provided `comment`.

To remove acknowledgement, you just need to delete Acknowledgement object.

> Note: Acknowledgement object name should be similar as Incident object name