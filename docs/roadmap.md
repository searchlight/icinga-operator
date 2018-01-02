---
title: Roadmap | Searchlight
description: Roadmap of searchlight
menu:
  product_searchlight_4.0.1:
    identifier: roadmap-searchlight
    name: Roadmap
    parent: getting-started
    weight: 30
product_name: searchlight
menu_name: product_searchlight_4.0.1
section_menu_id: getting-started
url: /products/searchlight/4.0.1/getting-started/roadmap/
---

# Project Status

## Versioning Policy
There are 2 parts to versioning policy:
 - Operator version: Searchlight __does not follow semver__, rather the _major_ version of operator points to the
Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version.
You can verify this from the `glide.yaml` file. This means there might be breaking changes
between point releases of the operator. This generally manifests as changed annotation keys or their meaning.
Please always check the release notes for upgrade instructions.
 - CRD version: monitoring.appscode.com/v1alpha1 is considered in alpha. This means breaking changes to the YAML format
might happen among different releases of the operator.

### External Dependencies
Searchlight 4.0.1 depends on the following version of Icinga2 and friends:

| Name                   | Version    |
|------------------------|------------|
| Icinga2                | 2.6.3-1    |
| Icingaweb2             | 2.4.1      |
| Monitoring Plugins     | 2.2-r1     |
| Postgres               | 9.5-alpine |