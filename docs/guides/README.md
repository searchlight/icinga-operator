---
title: Table of Contents | Guides
description: Table of Contents | Guides
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-readme
    name: Readme
    parent: guides
    weight: -1
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
url: /products/searchlight/6.0.0-rc.0/guides/
aliases:
  - /products/searchlight/6.0.0-rc.0/guides/README/
---
# Guides

Guides show you how to perform tasks with Searchlight.

- Cluster Alerts
  - [ca-cert](/docs/guides/cluster-alerts/ca-cert.md) - To check expiration of CA certificate used by Kubernetes api server.
  - [component-status](/docs/guides/cluster-alerts/component-status.md) - To check Kubernetes component status.
  - [event](/docs/guides/cluster-alerts/event.md) - To check Kubernetes Warning events.
  - [json-path](/docs/guides/cluster-alerts/json-path.md) - To check any JSON HTTP response using [jsonpath](https://kubernetes.io/docs/reference/kubectl/jsonpath/).
  - [node-exists](/docs/guides/cluster-alerts/node-exists.md) - To check existence of Kubernetes nodes.
  - [pod-exists](/docs/guides/cluster-alerts/pod-exists.md) - To check existence of Kubernetes pods.

- Node Alerts
  - [node-status](/docs/guides/node-alerts/node-status.md) - To check Kubernetes Node status.
  - [node-volume](/docs/guides/node-alerts/node-volume.md) - To check Node Disk stat.

- Pod Alerts
  - [pod-exec](/docs/guides/pod-alerts/pod-exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns Critical
  - [pod-status](/docs/guides/pod-alerts/pod-status.md) - To check Kubernetes pod status.
  - [pod-volume](/docs/guides/pod-alerts/pod-volume.md) - To check Pod volume stat.

- [Supported Notifiers](/docs/guides/notifiers.md): This article documents how to configure Searchlight to send notifications via Email, SMS or Chat.
