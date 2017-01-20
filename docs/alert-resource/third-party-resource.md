# Alert Third Party Resource

## Create the Alert Third Party Resource

Save the following contents to `alert.yaml`:

```
apiVersion: extensions/v1beta1
kind: ThirdPartyResource
description: "Alert support for Kubernetes by appscode.com"
metadata:
  name: alert.appscode.com
versions:
  - name: v1beta1
```

Submit the Third Party Resource configuration to the Kubernetes API server:

```
kubectl create -f alert.yaml
```

At this point you can now create [Alert Objects](docs/alert-resource/objects.md).
