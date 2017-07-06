[![Go Report Card](https://goreportcard.com/badge/github.com/appscode/searchlight)](https://goreportcard.com/report/github.com/appscode/searchlight)

# searchlight

<img src="cover.jpg">

Searchlight is an Alert Management project.
It has a Controller to watch Kubernetes Objects. Alert objects are consumed by Searchlight Controller to create Icinga2 hosts, services and notifications.

### Resource

Following resources are used in Searchlight

| Resource               | Version   |
| :---                   | :---      |
| Icinga2                | 2.4.8     |
| Icingaweb2             | 2.1.2     |
| Monitoring Plugins     | 2.1.2     |
| Postgres               | 9.5       |
| Searchlight Controller | 1.5.9     |

## Features

Searchlight supports additional custom plugins. Followings are currently added

| Check Command                                                           | Plugin                  | Details                                                                                       |
| :---                                                                    | :---                    | :---                                                                                          |
| [component_status](docs/check-command/component_status.md)   | check_component_status  | To check Kubernetes components                                                                |
| [influx_query](docs/check-command/influx_query.md)           | check_influx_query      | To check InfluxDB query result                                                                |
| [json_path](docs/check-command/json_path.md)                 | check_json_path         | To check any API response by parsing JSON using JQ queries                                    |
| [node_count](docs/check-command/node_count.md)               | check_node_count        | To check total number of Kubernetes node                                                      |
| [node_status](docs/check-command/node_status.md)             | check_node_status       | To check Kubernetes Node status                                                               |
| [pod_exists](docs/check-command/pod_exists.md)               | check_pod_exists        | To check Kubernetes pod existence                                                             |
| [pod_status](docs/check-command/pod_status.md)               | check_pod_status        | To check Kubernetes pod status                                                                |
| [prometheus_metric](docs/check-command/prometheus_metric.md) | check_prometheus_metric | To check Prometheus query result                                                              |
| [node_disk](docs/check-command/node_disk.md)                 | check_node_disk         | To check Node Disk stat                                                                       |
| [volume](docs/check-command/volume.md)                       | check_volume            | To check Pod volume stat                                                                      |
| [kube_event](docs/check-command/kube_event.md)               | check_kube_event        | To check all Kubernetes Warning events happened in last `c` seconds                           |
| [kube_exec](docs/check-command/kube_exec.md)                 | check_kube_exec         | To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL|

> Note: All of these plugins are combined into a single plugin called `hyperalert`

#### Supported Notifiers
Searchlight can send alert notification via following notifiers:

1. [Hipchat](docs/notifier/hipchat.md)
2. [Mailgun](docs/notifier/mailgun.md)
3. [SMTP](docs/notifier/smtp.md)
4. [Twilio](docs/notifier/twilio.md)
5. [Slack](docs/notifier/slack.md)
6. [Plivo](docs/notifier/plivo.md)

## User Guide

To deploy Searchlight in Kubernetes cluster, follow this [guide](docs/deployment-guide.md).
This guide will walk you through following three steps:

1. Creating Third Party Resource
2. Deploying Icinga2
3. Deploying Searchlight Controller

## Architectural Design

If you want to know how Searchlight Controller is working, read this [doc](docs/architecture-guide/controller.md).


## Contribution

If you're interested in being a contributor, read following guides:

* Build guides
    
    1. [Icinga2](docs/contribution-guide/icinga2/build.md)
    2. [Searchlight Controller](docs/contribution-guide/controller/build.md)
   
## Versioning Policy
There are 2 parts to versioning policy:
 - Operator version: Searchlight __does not follow semver__, rather the _major_ version of operator points to the
Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version.
You can verify this from the `glide.yaml` file. This means there might be breaking changes
between point releases of the operator. This generally manifests as changed annotation keys or their meaning.
Please always check the release notes for upgrade instructions.
 - TPR version: monitoring.appscode.com/v1alpha1 is considered in alpha. This means breaking changes to the YAML format
might happen among different releases of the operator.

---

**The searchlight operator collects anonymous usage statistics to help us learn how the software is being used and
how we can improve it. To disable stats collection, run the operator with the flag** `--analytics=false`.

---

## Support
If you have any questions, you can reach out to us.

* [Slack](https://slack.appscode.com)
* [Forum](https://discuss.appscode.com)
* [Twitter](https://twitter.com/AppsCodeHQ)
* [Website](https://appscode.com)
