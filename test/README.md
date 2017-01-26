# Test Icinga2 Custom Plugins

### Configure Kubernetes Client

It reads local `~/.kube/config` data and uses `current-context` for cluster and auth information.

Note:

* Set ENV `APPSCODE_ENV` to `dev` or make it empty


### Build `hyperalert`

    ./hack/make.py build_cmd hyperalert

 This will create a binary in `dist/hyperalert/` directory


### Run Test

Run following command to test

* __component_status__

        go test -v github.com/appscode/searchlight/test -run ^TestComponentStatus$

* __json_path__

        go test -v github.com/appscode/searchlight/test -run ^TestJsonPath$

* __node_count__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeCount$

* __node_status__

        go test -v github.com/appscode/searchlight/test -run ^TestNodeStatus$

> Following two will create temporary kubernetes objects to test

* __pod_exists__

        go test -v github.com/appscode/searchlight/test -run ^TestPodExists$

* __pod_status__

        go test -v github.com/appscode/searchlight/test -run ^TestPodStatus$


### Configure Icinga2 Client

* Set ENV `E2E_ICINGA_SECRET` to kubernetes `secret` name
* Set ENV `E2E_ICINGA_URL` to Icinga2 API url [default: Service LoadBalancer.Ingress]

Following information will be collected from secret:

1. ICINGA_K8S_SERVICE
2. ICINGA_API_USER
3. ICINGA_API_PASSWORD

`ICINGA_K8S_SERVICE` will be used to get `LoadBalancer.Ingress` if `E2E_ICINGA_URL` ENV is not set

#### __General Test__

    go test -v github.com/appscode/searchlight/test -run ^TestGeneralAlert$
