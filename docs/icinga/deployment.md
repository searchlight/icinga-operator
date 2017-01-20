# Deployment Guide

This guide will walk you through deploying the icinga.

### Resource

Following resources are used in this deployment

| Resource            | Version  |
| :---:               | :---:    |
| Icinga2             | 2.4.8    |
| Icingaweb2          | 2.1.2    |
| Monitoring Plugins  | 2.1.2    |
 AppsCode Custom Plugin ||


### Build Appscode Custom Plugin

We add a custom plugin named `hyperalert` to handle various CheckCommands. This plugin includes following commands:

* check_component_status
* check_influx_query
* check_json_path
* check_node_count
* check_node_status
* check_pod_exists
* check_pod_status
* check_prometheus_metric
* check_node_disk
* check_volume
* check_kube_event
* check_kube_exec
* notifier

See build instruction [here](docs/plugin/build.md).

Learn about usage of CheckCommands [here]().

## Deploy default Icinga

See deployment instruction [here](docs/icinga/k8s/deployment.md).

## Deploy Appscode Icinga

See deployment instruction [here](docs/icinga/appscode/deployment.md).
