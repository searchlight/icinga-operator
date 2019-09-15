---
title: Concepts | Searchlight
menu:
  product_searchlight_{{ .version }}:
    identifier: concepts-readme
    name: Readme
    parent: concepts
    weight: -1
product_name: searchlight
menu_name: product_searchlight_{{ .version }}
section_menu_id: concepts
url: /products/searchlight/{{ .version }}/concepts/
aliases:
  - /products/searchlight/{{ .version }}/concepts/README/
---
# Concepts

Concepts help you learn about the different parts of the Searchlight and the abstractions it uses.

- What is Searchlight?
  - [Overview](/docs/concepts/what-is-searhclight/overview.md). Provides a conceptual introduction to Searchlight, including the problems it solves and its high-level architecture.
- Types of Alerts
  - [ClusterAlerts](/docs/concepts/alert-types/cluster-alert.md). Introduces the concept of `ClusterAlert` to periodically run various checks on a Kubernetes cluster.
  - [NodeAlerts](/docs/concepts/alert-types/node-alert.md). Introduces the concept of `NodeAlert` to periodically run various checks on nodes in a Kubernetes cluster.
  - [PodAlerts](/docs/concepts/alert-types/pod-alert.md). Introduces the concept of `PodAlert` to periodically run various checks on pods in a Kubernetes cluster.
