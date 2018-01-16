[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/searchlight)](https://goreportcard.com/report/github.com/appscode/searchlight)

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
| [5.1.0](https://github.com/appscode/searchlight/releases/tag/5.1.0) (uses CRD)           | [User Guide](https://github.com/appscode/searchlight/tree/5.1.0/docs) | 1.7.x+             |
| [3.0.1](https://github.com/appscode/searchlight/releases/tag/3.0.1) (uses TPR)           | [User Guide](https://github.com/appscode/searchlight/tree/3.0.1/docs)      | 1.5.x - 1.7.x      |

## Installation
To install Searchlight, please follow the guide [here](/docs/setup/install.md).

## Using Searchlight
Want to learn how to use Searchlight? Please start [here](/docs/guides/README.md).

## Contribution guidelines
Want to help improve Searchlight? Please start [here](/docs/CONTRIBUTING.md).

## Project Status
Wondering what features are coming next? Please visit [here](/docs/roadmap.md).

---

**The searchlight operator collects anonymous usage statistics to help us learn how the software is being used and
how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Acknowledgement
 - Many thanks to [Icinga](https://www.icinga.com/) project.

## Support
We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [AppsCode Slack team](https://appscode.slack.com/messages/C8M7LT2QK/details/) channel `#searchlight_`. To sign up, use our [Slack inviter](https://slack.appscode.com/).

If you have found a bug with Searchlight or want to request for new features, please [file an issue](https://github.com/appscode/searchlight/issues/new).

