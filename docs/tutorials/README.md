# Tutorials

This section contains tutorials on how to use Searchlight. Please visit the links below to learn more:

 - [Cluster Snapshots](/docs/tutorials/cluster-snapshot.md): This tutorial will show you how to use Kubed to take periodic snapshots of a Kubernetes cluster objects.

 - [Kubernetes Recycle Bin](/docs/tutorials/recycle-bin.md): Kubed provides a recycle bin for deleted and/or updated Kubernetes objects. This tutorial will show you how to use Kubed to setup a recycle bin for Kubernetes cluster objects.

 - [Synchronize Configuration across Namespaces](/docs/tutorials/config-syncer.md): This tutorial will show you how Kubed can sync ConfigMaps/Secrets across Kubernetes namespaces.

 - [Forward Cluster Events](/docs/tutorials/event-forwarder.md): This tutorial will show you how to use Kubed to send notifications via Email, SMS or Chat for various cluster events.

 - [Using Janitors](/docs/tutorials/janitors.md): Kubernetes supports storing cluster logs in Elasticsearch and cluster metrics in InfluxDB. This tutorial will show you how to use kubed janitors to delete old data from Elasticsearch and InfluxDB.
































## What is ClusterAlert
A `ClusterAlert` is a Kubernetes `Third Party Object` (TPR). It provides declarative configuration of [Icinga services](https://www.icinga.com/docs/icinga2/latest/doc/09-object-types/#service) for cluster level alerts in a Kubernetes native way. You only need to describe the desired check command and notifier in a ClusterAlert object, and the Searchlight operator will create Icinga2 hosts, services and notifications to the desired state for you.


### Check Command
Check commands are used by Icinga to periodically test some condition. If the test return positive appropriate notifications are sent. The following check commands are supported for pods:

- [ca_cert](/docs/cluster-alerts/ca_cert.md) - To check expiration of CA certificate used by Kubernetes api server.
- [component_status](/docs/cluster-alerts/component_status.md) - To check Kubernetes component status.
- [event](/docs/cluster-alerts/event.md) - To check Kubernetes Warning events.
- [json_path](/docs/cluster-alerts/json_path.md) - To check any JSON HTTP response using [jq](https://stedolan.github.io/jq/).
- [node_exists](/docs/cluster-alerts/node_exists.md) - To check existence of Kubernetes nodes.
- [pod_exists](/docs/cluster-alerts/pod_exists.md) - To check existence of Kubernetes pods.
- [NodeAlerts] - This article introduces the concept of `NodeAlert` to periodically check various commands on nodes in a Kubernetes cluster. Also, visit the links below to learn about various check commands are supported for nodess:
  - [influx_query](/docs/node-alerts/influx_query.md) - To check InfluxDB query result.
  - [node_status](/docs/node-alerts/node_status.md) - To check Kubernetes Node status.
  - [node_volume](/docs/node-alerts/node_volume.md) - To check Node Disk stat.
- [PodAlerts] - This article introduces the concept of `PodAlert` to periodically check various commands on pods in a Kubernetes cluster. Also, visit the links below to learn about various check commands are supported for pods:
 - [influx_query](/docs/pod-alerts/influx_query.md) - To check InfluxDB query result.
 - [pod_exec](/docs/pod-alerts/pod_exec.md) - To check Kubernetes exec command. Returns OK if exit code is zero, otherwise, returns CRITICAL
 - [pod_status](/docs/pod-alerts/pod_status.md) - To check Kubernetes pod status.
 - [pod_volume](/docs/pod-alerts/pod_volume.md) - To check Pod volume stat.
 - [Supported Notifiers](/docs/tutorials/notifiers.md): This article documents how to configure Kubed to send notifications via Email, SMS or Chat
