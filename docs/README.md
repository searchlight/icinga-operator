---
title: Weclome | Searchlight
description: Welcome to Searchlight
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: readme-searchlight
    name: Readme
    parent: welcome
    weight: -1
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: welcome
url: /products/searchlight/6.0.0-rc.0/welcome/
aliases:
  - /products/searchlight/6.0.0-rc.0/
  - /products/searchlight/6.0.0-rc.0/README/
---
# Searchlight

<img src="/docs/images/cover.jpg">

Searchlight by AppsCode is a Kubernetes operator for [Icinga](https://www.icinga.com/). If you are running production workloads in Kubernetes, you probably want to be alerted when things go wrong. Icinga periodically runs various checks on a Kubernetes cluster and sends notifications if detects an issue. It also nicely supplements whitebox monitoring tools like, [Prometheus](https://prometheus.io/) with blackbox monitoring can catch problems that are otherwise invisible, and also serves as a fallback in case internal systems completely fail.

From here you can learn all about Searchlight's architecture and how to deploy and use Searchlight.

- [Concepts](/docs/concepts/). Concepts explain some significant aspect of Searchlight. This is where you can learn about what Searchlight does and how it does it.

- [Setup](/docs/setup/). Setup contains instructions for installing
  the Searchlight in various cloud providers.

- [Guides](/docs/guides/). Guides show you how to perform tasks with Searchlight.

- [Reference](/docs/reference/searchlight). Detailed exhaustive lists of command-line options for various Searchlight binaries.

We're always looking for help improving our documentation, so please don't hesitate to
[file an issue](https://github.com/appscode/searchlight/issues/new) if you see some problem.
Or better yet, submit your own [contributions](/docs/CONTRIBUTING.md) to help
make our docs better.

---

**Searchlight binaries collects anonymous usage statistics to help us learn how the software is being used and how we can improve it.
To disable stats collection, run the operator with the flag** `--enable-analytics=false`.

---