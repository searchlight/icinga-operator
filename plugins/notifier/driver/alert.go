package driver

import (
	"os"

	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/k8s"
)

func GetAlertInfo(namespace, alertName string) (*aci.Alert, error) {
	kubeClient, err := k8s.NewClient()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	return kubeClient.AppscodeExtensionClient.Alert(namespace).Get(alertName)
}
