---
title: Incident Concepts
description: Incident Concepts
menu:
  product_searchlight_8.0.0:
    identifier: incident-concepts
    parent: incident
    name: Incident Concepts
    weight: 15
menu_name: product_searchlight_8.0.0
---

# Incident

## What is Incident
A `Incident` is a Kubernetes `Custom Resource Definition` (CRD).
It provides information on notifications sent by Searchlight for alerts.

## Incident Spec
As with all other Kubernetes objects, a Incident has `apiVersion`, `kind`, and `metadata` fields. It also has a `.spec` section. 

Below is an example Incident object maintained by Searchlight.

```yaml
apiVersion: monitoring.appscode.com/v1alpha1
kind: Incident
metadata:
  name: cluster.pod-exists-demo-0.20180428-1109
  namespace: demo
  labels:
    monitoring.appscode.com/alert: pod-exists-demo-0
    monitoring.appscode.com/alert-type: cluster
    monitoring.appscode.com/recovered: true
status:
  lastNotificationType: Recovery
  notifications:
  - type: Problem
    checkOutput: Found 10 pod(s) instead of 11
    author: searchlight-user
    firstTimestamp: 20180428-1109
    lastTimestamp: 20180428-1119
    state: Critical
  - type: Acknowledgement
    checkOutput: Found 10 pod(s) instead of 11
    author: admin
    comment: working on fix
    firstTimestamp: 20180428-1142
    lastTimestamp: 20180428-1142
    state: Critical
  - type: Recovery
    checkOutput: Found all pods
    author: searchlight-user
    firstTimestamp: 20180428-1237
    lastTimestamp: 20180428-1237
    state: Ok
```

Here,

- `metadata.name` represents Incident name with format `{host-type}.{alert-name}.{time}`
- `metadata.namespace` represents Namespace where ClusterAlert is created.
- `metadata.labels` provides additional information on Alert and Notification
- `status` provides information on notifications
- `status.lastNotificationType` represents last type of notification that was sent
- `status.notifications` provides list of notifications that were sent

#### Notification List

Searchlight plugin `hyperalert` adds notification information in `status.notifications` when a notification is occurred in Icinga.

Lets see some examples to understand this scenario.

We have an Incident with following `status` part

```yaml
status:
  lastNotificationType: Problem
  notifications:
  - type: Problem
    checkOutput: Found 10 pod(s) instead of 11
    author: searchlight-user
    firstTimestamp: 20180428-1109
    lastTimestamp: 20180428-1109
    state: Critical
```

When there is no existing Incident available, Searchlight creates one. From above, we can see that, a Icinga notification of type **Problem** has invoked for **Critical** state.
Service `output` is updated with last Service check output.

If ClusterAlert is set to send notifications on every *5m*, more notifications will be invoked with *5m* interval. Now every time a notification of type **Problem** is invoked, Searchlight updates
this existing notification information. Following fields are updated:

- checkOutput
- lastTimestamp
- state

Now, suppose, user acknowledges this problem. So another notification of type **Acknowledgement** is invoked. When notification type is **Acknowledgement**, a new notification is added in list.

Following is the latest `status` of this Incident

```yaml
status:
  lastNotificationType: Acknowledgement
  notifications:
  - type: Problem
    checkOutput: Found 10 pod(s) instead of 11
    author: searchlight-user
    firstTimestamp: 20180428-1109
    lastTimestamp: 20180428-1119
    state: Critical
  - type: Acknowledgement
    checkOutput: Found 10 pod(s) instead of 11
    author: admin
    comment: working on fix
    firstTimestamp: 20180428-1142
    lastTimestamp: 20180428-1142
    state: Critical
```

User `admin` acknowledges this issues with comment `working on fix`. `status.lastNotificationType` is updated with last notification type.

To know how to acknowledge, see [here](/docs/concepts/incident/acknowledgement.md).

When this problem is fixed, another notification of type **Recovery** is invoked. With this notification, this Incident is locked. For furthers incidents, another new Incident object will be created.

This is the final `status` of this Incident

```yaml
status:
  lastNotificationType: Recovery
  notifications:
  - type: Problem
    checkOutput: Found 10 pod(s) instead of 11
    author: searchlight-user
    firstTimestamp: 20180428-1109
    lastTimestamp: 20180428-1119
    state: Critical
  - type: Acknowledgement
    checkOutput: Found 10 pod(s) instead of 11
    author: admin
    comment: working on fix
    firstTimestamp: 20180428-1142
    lastTimestamp: 20180428-1142
    state: Critical
  - type: Recovery
    checkOutput: Found all pods
    author: searchlight-user
    firstTimestamp: 20180428-1237
    lastTimestamp: 20180428-1237
    state: Ok
```

And also, label `monitoring.appscode.com/recovered: true` is added in label. This represents that, This Incident is recovered.

