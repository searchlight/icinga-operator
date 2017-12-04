---
title: Overview | Searchlight
description: Overview of Searchlight
menu:
  product_searchlight_4.0.0:
    identifier: overview-searchlight
    name: Overview
    parent: getting-started
    weight: 20
product_name: searchlight
menu_name: product_searchlight_4.0.0
section_menu_id: getting-started
url: /products/searchlight/4.0.0/getting-started/
aliases:
  - /products/searchlight/4.0.0/
  - /products/searchlight/4.0.0/README/
---

[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/searchlight)](https://goreportcard.com/report/github.com/appscode/searchlight)

# Searchlight

<img src="/cover.jpg">


Searchlight by AppsCode is a Kubernetes operator for [Icinga](https://www.icinga.com/). If you are running production workloads in Kubernetes, you probably want to be alerted when things go wrong. Icinga periodically runs various checks on a Kubernetes cluster and sends notifications if detects an issue. It also nicely supplements whitebox monitoring tools like, [Prometheus](https://prometheus.io/) with blackbox monitoring can catch problems that are otherwise invisible, and also serves as a fallback in case internal systems completely fail. Searchlight is a CRD controller for Kubernetes built around Icinga to address these issues. Searchlight can do the following things for you:

 - Periodically run various checks on a Kubernetes cluster and its nodes or pods.
 - Includes a [suite of check commands](/docs/reference/hyperalert/hyperalert.md) written specifically for Kubernetes.
 - Searchlight can send notifications via Email, SMS or Chat.
 - [Supplements](https://prometheus.io/docs/practices/alerting/#metamonitoring) the whitebox monitoring tools like [Prometheus](https://prometheus.io).

## Supported Versions
Please pick a version of Searchlight that matches your Kubernetes installation.

| Searchlight Version                                                                      | Docs                                                                       | Kubernetes Version |
|------------------------------------------------------------------------------------------|----------------------------------------------------------------------------|--------------------|
| [4.0.0](https://github.com/appscode/searchlight/releases/tag/4.0.0) (uses CRD) | [User Guide](https://github.com/appscode/searchlight/tree/4.0.0/docs) | 1.7.x+             |
| [3.0.1](https://github.com/appscode/searchlight/releases/tag/3.0.1) (uses TPR)           | [User Guide](https://github.com/appscode/searchlight/tree/3.0.1/docs)      | 1.5.x - 1.7.x      |

## Installation
To install Searchlight, please follow the guide [here](/docs/install.md).

## Using Searchlight
Want to learn how to use Searchlight? Please start [here](/docs/tutorials/README.md).

## Contribution guidelines
Want to help improve Searchlight? Please start [here](/CONTRIBUTING.md).

## Project Status
Wondering what features are coming next? Please visit [here](/ROADMAP.md).

---

**The searchlight operator collects anonymous usage statistics to help us learn how the software is being used and
how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Acknowledgement
 - Many thanks to [Icinga](https://www.icinga.com/) project.

## Support
If you have any questions, you can reach out to us.
* [Slack](https://slack.appscode.com)
* [Twitter](https://twitter.com/AppsCodeHQ)
* [Website](https://appscode.com)
