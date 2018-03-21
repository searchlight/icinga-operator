package notifier

import (
	"fmt"
	"testing"
	"time"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
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
