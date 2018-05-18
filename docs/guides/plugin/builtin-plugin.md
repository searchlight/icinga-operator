---
title: Searchlight Plugin
menu:
  product_searchlight_6.0.0-rc.0:
    identifier: guides-searchlight-plugin
    name: SearchlightPlugin
    parent: searchlight-plugin
    weight: 20
product_name: searchlight
menu_name: product_searchlight_6.0.0-rc.0
section_menu_id: guides
---

> New to SearchlightPlugin? Please start [here](/docs/setup/developer-guide/webhook-plugin.md).

# Builtin Searchlight Plugin

Searchlight Icinga container includes a plugin `hyperalert`. This plugin has following Available Commands:

- [check_ca_cert](./docs/guides/cluster-alerts/ca-cert.md) - Check Certificate expire date
- [check_cert](./docs/guides/cluster-alerts/cert.md) - Check Certificate expire date
- [check_component_status](./docs/guides/cluster-alerts/component-status.md) - Check Kubernetes Component Status
- [check_event](./docs/guides/cluster-alerts/event.md) - Check kubernetes events for all namespaces
- [check_json_path](./docs/guides/cluster-alerts/node-exists.md) - Check Json Object
- [check_node_exists](./docs/guides/cluster-alerts/node-exists.md) - Count Kubernetes Nodes
- [check_node_status](./docs/guides/node-alerts/node-status.md) - Check Kubernetes Node
- [check_pod_exec](./docs/guides/pod-alerts/pod-exec.md) - Check exit code of exec command on Kubernetes container
- [check_pod_exists](./docs/guides/cluster-alerts/pod-exists.md) - Check Kubernetes Pod(s)
- [check_pod_status](./docs/guides/pod-alerts/pod-status.md) - Check Kubernetes Pod(s) status
- [check_volume](./docs/guides/pod-alerts/pod-volume.md) - Check kubernetes volume
- [notifier](./docs/guides/notifiers.md) - AppsCode Icinga2 Notifier

To use these commands, you need to register the CheckCommand first by creating SearchlightPlugin.

And also, to unregister the CheckCommand, delete the SearchlightPlugin.
