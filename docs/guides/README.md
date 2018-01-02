---
title: Table of Contents | Guides
description: Table of Contents | Guides
menu:
  product_searchlight_5.0.0:
    identifier: guides-readme
    name: Readme
    parent: guides
    weight: -1
product_name: searchlight
menu_name: product_searchlight_5.0.0
section_menu_id: guides
url: /products/searchlight/5.0.0/guides/
aliases:
  - /products/searchlight/5.0.0/guides/README/
---
# Guides

Guides show you how to perform tasks with Searchlight.

- Cluster Alerts
  - [ca_cert](/docs/guides/cluster-alerts/ca_cert.md) - To check expiration of CA certificate used by Kubernetes api server.
  - [component_status](/docs/guides/cluster-alerts/component_status.md) - To check Kubernetes component status.
  - [event](/docs/guides/cluster-alerts/event.md) - To check Kubernetes Warning events.
  - [json_path](/docs/guides/cluster-alerts/json_path.md) - To check any JSON HTTP response using [jq](https://stedolan.github.io/jq/).
  - [node_exists](/docs/guides/cluster-alerts/node_exists.md) - To check existence of Kubernetes nodes.
  - [pod_exists](/docs/guides/cluster-alerts/pod_exists.md) - To check existence of Kubernetes pods.

- Node Alerts
  - [influx_query](/docs/guides/node-alerts/influx_query.md) - To check InfluxDB query result.
  - [node_status](/docs/guides/node-alerts/node_status.md) - To check Kubernetes Node status.
  - [node_volume](/docs/guides/node-alerts/node_volume.md) - To check Node Disk stat.

- Pod Alerts
  - [influx_query](/docs/guides/pod-alerts/influx_query.md) - To check InfluxDB query result.
  - [pod_exec](/docs/guides/pod-alerts/pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
  - [pod_status](/docs/guides/pod-alerts/pod_status.md) - To check Kubernetes pod status.
  - [pod_volume](/docs/guides/pod-alerts/pod_volume.md) - To check Pod volume stat.

- [Supported Notifiers](/docs/guides/notifiers.md): This article documents how to configure Searchlight to send notifications via Email, SMS or Chat.
