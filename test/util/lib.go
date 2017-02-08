package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/cmd/searchlight/app"
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

func IcingaHostSearchQuery(objectList []*host.KubeObjectInfo) string {
	matchHost := ""
	for id, object := range objectList {
		if id > 0 {
			matchHost = matchHost + "||"
		}
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",host.name)`, object.Name)
	}
	return fmt.Sprintf(`{"filter": "(%s)"}`, matchHost)
}

func countIcingaService(watcher *app.Watcher, alert *aci.Alert, expectZero bool) error {

	checkCommand := alert.Spec.CheckCommand
	objectType := alert.Labels["alert.appscode.com/objectType"]
	objectName := alert.Labels["alert.appscode.com/objectName"]
	namespace := alert.Namespace
	// create all alerts for pod_status
	hostType, err := GetIcingaHostType(checkCommand, objectType)
	if err != nil {
		return err
	}
	objectList, err := host.GetObjectList(watcher.Client, checkCommand, hostType, namespace, objectType, objectName, "")
	if err != nil {
		return err
	}

	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)

	in := host.IcingaServiceSearchQuery(serviceName, objectList)
	var respService host.ResponseObject

	try := 0
	for {
		time.Sleep(time.Second * 30)

		if _, err = watcher.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
			return errors.New("can't check icinga service")
		}

		if expectZero {
			if len(respService.Results) != 0 {
				err = errors.New(fmt.Sprintf("Service Found for %s:%s", objectType, objectName))
			}
		} else {
			if len(respService.Results) != len(objectList) {
				err = errors.New(fmt.Sprintf("Total Service Mismatch for %s:%s", objectType, objectName))
			}
		}
		if err != nil {
			fmt.Println(err.Error())
		} else {
			break
		}
		if try > 5 {
			return err
		}

		fmt.Println("--> Waiting for 30 second more in count process")

	}

	return nil
}

func countIcingaHost(watcher *app.Watcher, fakeAlert *aci.Alert, expectZero bool) error {
	objectType, objectName := host.GetObjectInfo(fakeAlert.Labels)
	checkCommand := fakeAlert.Spec.CheckCommand

	// create all alerts for pod_status
	hostType, err := GetIcingaHostType(checkCommand, objectType)
	if err != nil {
		return err
	}
	objectList, err := host.GetObjectList(watcher.Client, checkCommand, hostType, fakeAlert.Namespace, objectType, objectName, "")
	if err != nil {
		return err
	}

	in := IcingaHostSearchQuery(objectList)
	var respHost host.ResponseObject

	try := 0
	for {
		time.Sleep(time.Second * 30)

		if _, err = watcher.IcingaClient.Objects().Hosts("").Get([]string{}, in).Do().Into(&respHost); err != nil {
			return errors.New("can't check icinga service")
		}

		if expectZero {
			if len(respHost.Results) != 0 {
				err = errors.New(fmt.Sprintf("Host Found for %s:%s", objectType, objectName))
			}
		} else {
			if len(respHost.Results) != len(objectList) {
				err = errors.New(fmt.Sprintf("Total Host Mismatch for %s:%s", objectType, objectName))
			}
		}
		if err != nil {
			fmt.Println(err.Error())
		} else {
			break
		}
		if try > 5 {
			return err
		}

		fmt.Println("--> Waiting for 30 second more in count process")

	}

	return nil
}

func CheckIcingaObjects(watcher *app.Watcher, fakeAlert *aci.Alert, expectZeroHost, expectZeroService bool) (err error) {
	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Host")
	if err = countIcingaHost(watcher, fakeAlert, expectZeroHost); err != nil {
		return
	}

	// Count Icinga Service for 1st Alert. Should be found
	fmt.Println("----> Counting Icinga Service")
	if err = countIcingaService(watcher, fakeAlert, expectZeroService); err != nil {
		return
	}
	return
}