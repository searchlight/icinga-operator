## Searchlight Watcher

Searchlight Watcher watches following Kubernetes objects:

* Service
* StatefulSet
* DaemonSet
* ReplicaSet
* ReplicationController
* Pod
* Alert
* Node

Events on following objects are detected by Searchlight Controller:

* [Alert](../user-guide/alert-object.md) (The Third Party Resource)
* Pod
* Node
* Service

Other objects are watched only to find ancestors of pods.

Keep in Mind:

1. When Searchlight Controller starts or restarts, it starts with empty cache.
2. Watcher starts watching and caching all Kubernetes objects.
3. Controller detects all objects of Kubernetes type Alert, Pod, Node and Service.
4. And finally all alerts are reassigned.
