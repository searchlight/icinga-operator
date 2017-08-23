# Project Status

## Versioning Policy
There are 2 parts to versioning policy:
 - Operator version: Searchlight __does not follow semver__, rather the _major_ version of operator points to the
Kubernetes [client-go](https://github.com/kubernetes/client-go#branches-and-tags) version.
You can verify this from the `glide.yaml` file. This means there might be breaking changes
between point releases of the operator. This generally manifests as changed annotation keys or their meaning.
Please always check the release notes for upgrade instructions.
 - TPR version: monitoring.appscode.com/v1alpha1 is considered in alpha. This means breaking changes to the YAML format
might happen among different releases of the operator.

### Release 3.x.x
This is going to be the supported release for Kubernetes 1.5 & 1.6 .

### Release 4.x.x
This relased will be based on client-go 4.0.0. This is going to include a number of breaking changes (example, use CustomResoureDefinition instead of TPRs) and be supported for Kubernetes 1.7+. Please see the issues in release milestone [here](https://github.com/appscode/searchlight/milestone/3).

### External Dependencies
Searchlight 3.0.1 depends on the following version of Icinga2 and friends:

| Name                   | Version    |
|------------------------|------------|
| Icinga2                | 2.6.3-1    |
| Icingaweb2             | 2.4.1      |
| Monitoring Plugins     | 2.2-r1     |
| Postgres               | 9.5-alpine |
