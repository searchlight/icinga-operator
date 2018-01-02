---
title: Readme
description: Readme of guides
menu:
  product_searchlight_4.0.1:
    identifier: guides-readme
    name: Readme
    parent: guides
    weight: -1
product_name: searchlight
menu_name: product_searchlight_4.0.1
section_menu_id: guides
url: /products/searchlight/4.0.1/guides/
aliases:
  - /products/searchlight/4.0.1/guides/README/
---

# Guides

This section contains guides on how to use Searchlight. Please visit the links below to learn more:

 - [ClusterAlerts](/docs/guides/cluster-alerts/README.md) - This article introduces the concept of `ClusterAlert` to periodically run various checks on a Kubernetes cluster. Also, visit the links below to learn about the available check commands for a cluster:
    - [ca_cert](/docs/guides/cluster-alerts/ca_cert.md) - To check expiration of CA certificate used by Kubernetes api server.
    - [component_status](/docs/guides/cluster-alerts/component_status.md) - To check Kubernetes component status.
    - [event](/docs/guides/cluster-alerts/event.md) - To check Kubernetes Warning events.
    - [json_path](/docs/guides/cluster-alerts/json_path.md) - To check any JSON HTTP response using [jq](https://stedolan.github.io/jq/).
    - [node_exists](/docs/guides/cluster-alerts/node_exists.md) - To check existence of Kubernetes nodes.
    - [pod_exists](/docs/guides/cluster-alerts/pod_exists.md) - To check existence of Kubernetes pods.

 - [NodeAlerts](/docs/guides/node-alerts/README.md) - This article introduces the concept of `NodeAlert` to periodically run various checks on nodes in a Kubernetes cluster. Also, visit the links below to learn about the available check commands for nodes:
    - [influx_query](/docs/guides/node-alerts/influx_query.md) - To check InfluxDB query result.
    - [node_status](/docs/guides/node-alerts/node_status.md) - To check Kubernetes Node status.
    - [node_volume](/docs/guides/node-alerts/node_volume.md) - To check Node Disk stat.

 - [PodAlerts](/docs/guides/pod-alerts/README.md) - This article introduces the concept of `PodAlert` to periodically run various checks on pods in a Kubernetes cluster. Also, visit the links below to learn about the available check commands for pods:
    - [influx_query](/docs/guides/pod-alerts/influx_query.md) - To check InfluxDB query result.
    - [pod_exec](/docs/guides/pod-alerts/pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
    - [pod_status](/docs/guides/pod-alerts/pod_status.md) - To check Kubernetes pod status.
    - [pod_volume](/docs/guides/pod-alerts/pod_volume.md) - To check Pod volume stat.

 - [Supported Notifiers](/docs/guides/notifiers.md): This article documents how to configure Searchlight to send notifications via Email, SMS or Chat.
