package util

import (
	"errors"
	"fmt"
	"strings"
	"time"

	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/controller/host"
)

func getIcingaHostType(commandName, objectType string) (string, error) {
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

func icingaHostSearchQuery(objectList []*host.KubeObjectInfo) string {
	matchHost := ""
	for id, object := range objectList {
		if id > 0 {
			matchHost = matchHost + "||"
		}
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",host.name)`, object.Name)
	}
	return fmt.Sprintf(`{"filter": "(%s)"}`, matchHost)
}

func countIcingaService(watcher *app.Watcher, objectList []*host.KubeObjectInfo, serviceName string, expectZero bool) error {
	in := host.IcingaServiceSearchQuery(serviceName, objectList)
	var respService host.ResponseObject

	try := 0
	for {
		var err error
		if _, err = watcher.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
			err = errors.New("can't check icinga service")
		} else {
			if expectZero {
				if len(respService.Results) != 0 {
					err = errors.New("Service Found")
				}
			} else {
				if len(respService.Results) != len(objectList) {
					err = errors.New("Total Service Mismatch")
				}
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
		time.Sleep(time.Second * 30)
		try++
	}

	return nil
}

func countIcingaHost(watcher *app.Watcher, objectList []*host.KubeObjectInfo, expectZero bool) error {
	in := icingaHostSearchQuery(objectList)
	var respHost host.ResponseObject

	try := 0
	for {
		var err error
		if _, err = watcher.IcingaClient.Objects().Hosts("").Get([]string{}, in).Do().Into(&respHost); err != nil {
			err = errors.New("can't check icinga service")
		} else {
			if expectZero {
				if len(respHost.Results) != 0 {
					err = errors.New("Host Found")
				}
			} else {
				if len(respHost.Results) != len(objectList) {
					err = errors.New("Total Host Mismatch")
				}
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
		time.Sleep(time.Second * 30)
		try++
	}

	return nil
}

func GetObjectList(watcher *app.Watcher, alert *aci.Alert) ([]*host.KubeObjectInfo, error) {
	objectType, objectName := host.GetObjectInfo(alert.Labels)
	checkCommand := alert.Spec.CheckCommand

	// create all alerts for pod_status
	hostType, err := getIcingaHostType(checkCommand, objectType)
	if err != nil {
		return nil, err
	}
	objectList, err := host.GetObjectList(watcher.Client, checkCommand, hostType, alert.Namespace, objectType, objectName, "")
	if err != nil {
		return nil, err
	}

	return objectList, nil
}

func CheckIcingaObjectsForAlert(watcher *app.Watcher, alert *aci.Alert, expectZeroHost, expectZeroService bool) (err error) {
	objectList, err := GetObjectList(watcher, alert)
	if err != nil {
		return err
	}

	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Host")
	if err = countIcingaHost(watcher, objectList, expectZeroHost); err != nil {
		return
	}

	// Count Icinga Service for 1st Alert. Should be found
	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)
	fmt.Println("----> Counting Icinga Service")
	if err = countIcingaService(watcher, objectList, serviceName, expectZeroService); err != nil {
		return
	}
	return
}

func CheckIcingaObjects(watcher *app.Watcher, alert *aci.Alert, objectList []*host.KubeObjectInfo, expectZeroHost, expectZeroService bool) (err error) {
	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Host")
	if err = countIcingaHost(watcher, objectList, expectZeroHost); err != nil {
		return
	}

	// Count Icinga Service for 1st Alert. Should be found
	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)
	fmt.Println("----> Counting Icinga Service")
	if err = countIcingaService(watcher, objectList, serviceName, expectZeroService); err != nil {
		return
	}
	return
}

func CheckIcingaObjectsForPod(watcher *app.Watcher, podName, namespace string, expectedService int32) error {

	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Service")

	objectList := []*host.KubeObjectInfo{
		&host.KubeObjectInfo{
			Name: fmt.Sprintf("%v@%v", podName, namespace),
		},
	}

	in := icingaHostSearchQuery(objectList)
	var respService host.ResponseObject

	try := 0
	for {
		var err error
		if _, err = watcher.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
			return errors.New("can't check icinga service")
		}

		validService := int32(0)
		for _, service := range respService.Results {
			if service.Attrs.Name != "ping4" {
				validService++
			}
		}

		if expectedService != validService {
			err = errors.New("Service Mismatch")
			fmt.Println(err.Error())
		} else {
			break
		}

		if try > 5 {
			return err
		}

		fmt.Println("--> Waiting for 30 second more in count process")
		time.Sleep(time.Second * 30)
		try++
	}

	return nil
}
