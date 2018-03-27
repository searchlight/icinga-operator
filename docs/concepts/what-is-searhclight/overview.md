---
title: Searchlight Overview
description: Searchlight Overview
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: overview-concepts
    name: Overview
    parent: what-is-searchlight
    weight: 10
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: concepts
---

# Searchlight

<img src="/docs/images/cover.jpg">


Searchlight by AppsCode is a Kubernetes operator for [Icinga](https://www.icinga.com/). If you are running production workloads in Kubernetes, you probably want to be alerted when things go wrong. Icinga periodically runs various checks on a Kubernetes cluster and sends notifications if detects an issue. It also nicely supplements whitebox monitoring tools like, [Prometheus](https://prometheus.io/) with blackbox monitoring can catch problems that are otherwise invisible, and also serves as a fallback in case internal systems completely fail. Searchlight is a CRD controller for Kubernetes built around Icinga to address these issues. Searchlight can do the following things for you:

 - Periodically run various checks on a Kubernetes cluster and its nodes or pods.
 - Includes a [suite of check commands](/docs/reference/hyperalert/hyperalert.md) written specifically for Kubernetes.
 - Searchlight can send notifications via Email, SMS or Chat.
 - [Supplements](https://prometheus.io/docs/practices/alerting/#metamonitoring) the whitebox monitoring tools like [Prometheus](https://prometheus.io).
