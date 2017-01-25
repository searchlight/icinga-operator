# Alert Third Party Resource

## Create the Alert Third Party Resource

Save the following contents to `alert-third-party-resource.yaml`:

```yaml
apiVersion: extensions/v1beta1
kind: ThirdPartyResource
description: "Alert support for Kubernetes by appscode.com"
metadata:
  name: alert.appscode.com
versions:
  - name: v1beta1
```

Submit the Third Party Resource configuration to the Kubernetes API server:

```sh
kubectl create -f alert-third-party-resource.yaml
```

At this point we can now create [Kubernetes Alert Objects](objects.md).
