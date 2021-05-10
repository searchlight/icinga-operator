[![Go Report Card](https://goreportcard.com/badge/go.searchlight.dev/icinga-operator)](https://goreportcard.com/report/go.searchlight.dev/icinga-operator)
[![Build Status](https://travis-ci.org/searchlight/searchlight.svg?branch=master)](https://travis-ci.org/searchlight/searchlight)
[![codecov](https://codecov.io/gh/searchlight/searchlight/branch/master/graph/badge.svg)](https://codecov.io/gh/searchlight/searchlight)
[![Docker Pulls](https://img.shields.io/docker/pulls/appscode/searchlight.svg)](https://hub.docker.com/r/appscode/searchlight/)
[![Slack](https://slack.appscode.com/badge.svg)](https://slack.appscode.com)
[![Twitter](https://img.shields.io/twitter/follow/appscodehq.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=AppsCodeHQ)

# Searchlight

<img src="/docs/images/cover.jpg">


Searchlight by AppsCode is a Kubernetes operator for [Icinga](https://www.icinga.com/). If you are running production workloads in Kubernetes, you probably want to be alerted when things go wrong. Icinga periodically runs various checks on a Kubernetes cluster and sends notifications if detects an issue. It also nicely supplements whitebox monitoring tools like, [Prometheus](https://prometheus.io/) with blackbox monitoring can catch problems that are otherwise invisible, and also serves as a fallback in case internal systems completely fail. Searchlight is a CRD controller for Kubernetes built around Icinga to address these issues. Searchlight can do the following things for you:

 - Periodically run various checks on a Kubernetes cluster and its nodes or pods.
 - Includes a [suite of check commands](/docs/reference/hyperalert/hyperalert.md) written specifically for Kubernetes.
 - Searchlight can send notifications via Email, SMS or Chat.
 - [Supplements](https://prometheus.io/docs/practices/alerting/#metamonitoring) the whitebox monitoring tools like [Prometheus](https://prometheus.io).

## Supported Versions
Please pick a version of Searchlight that matches your Kubernetes installation.

| Searchlight Version                                                                      | Docs                                                                       | Kubernetes Version |
|------------------------------------------------------------------------------------------|----------------------------------------------------------------------------|--------------------|
| [8.0.0](https://go.searchlight.dev/icinga-operator/releases/tag/8.0.0) (uses CRD) | [User Guide](https://appscode.com/products/searchlight/8.0.0/welcome/)| 1.9.x+ (test/qa clusters) |
| [7.0.0](https://go.searchlight.dev/icinga-operator/releases/tag/7.0.0) (uses CRD)           | [User Guide](https://appscode.com/products/searchlight/7.0.0/welcome/)     | 1.8.x              |
| [5.1.1](https://go.searchlight.dev/icinga-operator/releases/tag/5.1.1) (uses CRD)           | [User Guide](https://appscode.com/products/searchlight/5.1.1/welcome/)     | 1.7.x+             |
| [3.0.1](https://go.searchlight.dev/icinga-operator/releases/tag/3.0.1) (uses TPR)           | [User Guide](https://go.searchlight.dev/icinga-operator/tree/3.0.1/docs)      | 1.5.x - 1.7.x      |

## Installation
To install Searchlight, please follow the guide [here](https://appscode.com/products/searchlight/8.0.0/setup/install).

## Using Searchlight
Want to learn how to use Searchlight? Please start [here](https://appscode.com/products/searchlight/8.0.0).

## Searchlight API Clients
You can use Searchlight api clients to programmatically access its CRD objects. Here are the supported clients:

- Go: [https://go.searchlight.dev/icinga-operator](/client/clientset/versioned)
- Java: https://github.com/searchlight-client/java

## Contribution guidelines
Want to help improve Searchlight? Please start [here](https://appscode.com/products/searchlight/8.0.0/welcome/contributing).

---

**Searchlight binaries collects anonymous usage statistics to help us learn how the software is being used and
how we can improve it. To disable stats collection, run the operator with the flag** `--enable-analytics=false`.

---

## Acknowledgement
 - Many thanks to [Icinga](https://www.icinga.com/) project.

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8M7LT2QK/details/) channel `#searchlight_`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Searchlight or want to request for new features, please [file an issue](https://go.searchlight.dev/icinga-operator/issues/new).

