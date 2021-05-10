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

package notifier

import (
	"fmt"
	"testing"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRenderMail(t *testing.T) {
	alert := api.ClusterAlert{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ca-cert-demo",
			Namespace: metav1.NamespaceDefault,
		},
		Spec: api.ClusterAlertSpec{
			Check:              api.CheckPodExists,
			CheckInterval:      metav1.Duration{Duration: 1 * time.Minute},
			AlertInterval:      metav1.Duration{Duration: 5 * time.Minute},
			NotifierSecretName: "notifier-conf",
			Vars: map[string]string{
				"name": "busybox",
			},
		},
	}
	hostname := "demo@cluster"
	host, err := icinga.ParseHost(hostname)
	assert.Nil(t, err)

	opts := options{
		hostname:         hostname,
		alertName:        alert.Name,
		notificationType: "WHAT_IS_THE_CORRECT_VAL?",
		serviceState:     "Warning",
		serviceOutput:    "Check command output",
		time:             time.Now(),
		author:           "<searchight-user>",
		comment:          "This is a test",
		host:             host,
	}

	config, err := newPlugin(nil, nil, opts).RenderMail(&alert)
	fmt.Println(err)
	assert.Nil(t, err)
	fmt.Println(config)
}
