/*
Copyright AppsCode Inc. and Contributors

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

package check_pod_status

import (
	"testing"
	"time"

	"go.searchlight.dev/icinga-operator/client/clientset/versioned/scheme"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"
	clientSetScheme "k8s.io/client-go/kubernetes/scheme"
)

var (
	cs *fake.Clientset
)

const (
	TIMEOUT = 2 * time.Minute
)

func TestPlugin_Check(t *testing.T) {
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)
	RunSpecsWithDefaultAndCustomReporters(t, "check_pod_status Suite", []Reporter{})
}

var _ = BeforeSuite(func() {
	scheme.AddToScheme(clientSetScheme.Scheme)
	cs = fake.NewSimpleClientset()
})

var _ = AfterSuite(func() {})
