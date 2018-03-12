package notifier

import (
	"fmt"
	"testing"
	"time"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
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
	req := Request{
		HostName:  "demo@cluster",
		AlertName: alert.Name,
		Type:      "WHAT_IS_THE_CORRECT_VAL?",
		State:     "WARNING",
		Output:    "Check command output",
		Time:      time.Now(),
		Author:    "<searchight-user>",
		Comment:   "This is a test",
	}
	config, err := RenderMail(&alert, &req)
	fmt.Println(err)
	assert.Nil(t, err)
	fmt.Println(config)
}
