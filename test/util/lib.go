package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client"
	"github.com/appscode/searchlight/pkg/controller/host"
)

func GetIcingaHostType(commandName, objectType string) (string, error) {
	icingaData, err := data.LoadIcingaData()
	if err != nil {
		return "", err
	}

	for _, command := range icingaData.Command {
		if command.Name == commandName {
			if t, found := command.ObjectToHost[objectType]; found {
				return t, nil
			}
		}
	}
	return "", errors.New("Icinga host_type not found")
}

func IcingaServiceSearchQuery(icingaServiceName string, objectList []*host.KubeObjectInfo) string {
	matchHost := ""
	for id, object := range objectList {
		if id > 0 {
			matchHost = matchHost + "||"
		}
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",host.name)`, object.Name)
	}
	return fmt.Sprintf(`{"filter": "(%s)&&match(\"%s\",service.name)"}`, matchHost, icingaServiceName)
}
