package driver

import (
	"os"

	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/util"
)

func GetAlertInfo(namespace, alertName string) (*aci.PodAlert, error) {
	kubeClient, err := util.NewClient()
	if err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
	return kubeClient.ExtClient.PodAlerts(namespace).Get(alertName)
}
