/*
Copyright 2018 The Searchlight Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package fake

import (
	internalversion "github.com/appscode/searchlight/client/clientset/internalversion/typed/monitoring/internalversion"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeMonitoring struct {
	*testing.Fake
}

func (c *FakeMonitoring) ClusterAlerts(namespace string) internalversion.ClusterAlertInterface {
	return &FakeClusterAlerts{c, namespace}
}

func (c *FakeMonitoring) Incidents(namespace string) internalversion.IncidentInterface {
	return &FakeIncidents{c, namespace}
}

func (c *FakeMonitoring) NodeAlerts(namespace string) internalversion.NodeAlertInterface {
	return &FakeNodeAlerts{c, namespace}
}

func (c *FakeMonitoring) PodAlerts(namespace string) internalversion.PodAlertInterface {
	return &FakePodAlerts{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeMonitoring) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
