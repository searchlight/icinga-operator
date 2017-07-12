# Test Icinga2 Custom Plugins

### Configure Kubernetes Client

It reads local `~/.kube/config` data and uses `current-context` for cluster and auth information.

### Run Test

Run following command to test Plugins

* __component_status__

        go test -v github.com/appscode/searchlight/test -run ^TestComponentStatus$

* __json_path__

        go test -v github.com/appscode/searchlight/test -run ^TestJsonPath$

* __kube_event__

        go test -v github.com/appscode/searchlight/test -run ^TestKubeEvent$

* __node_count__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeCount$

* __node_status__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeStatus$

> Following two will create temporary kubernetes objects to test

* __kube_exec__

        go test -v github.com/appscode/searchlight/test -run ^TestKubeExec$

* __pod_exists & pod_status__

        go test -v github.com/appscode/searchlight/test -run ^TestPodExistsPodStatus$

> To run all Test

    go test -v github.com/appscode/searchlight/test

Run following command for E2E test

* __MultipleAlerts__

        go test -v github.com/appscode/searchlight/test -run ^TestMultipleAlerts$

* __MultipleAlertsOnMultipleObjects__

        go test -v github.com/appscode/searchlight/test -run ^TestMultipleAlertsOnMultipleObjects$
