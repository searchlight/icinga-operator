package util

/*import (
	"errors"
	"fmt"
	"strings"
	"time"

	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
)*/

/*

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

func icingaHostSearchQuery(objectList []*icinga.IcingaHost) string {
	matchHost := ""
	for id, object := range objectList {
		if id > 0 {
			matchHost = matchHost + "||"
		}
		matchHost = matchHost + fmt.Sprintf(`match(\"%s\",icinga.name)`, object.Hostname)
	}
	return fmt.Sprintf(`{"filter": "(%s)"}`, matchHost)

}

func countIcingaService(w *controller.Controller, objectList []*icinga.IcingaHost, serviceName string, expectZero bool) error {
	in := icinga.IcingaServiceSearchQuery(serviceName, objectList)
	var respService icinga.ResponseObject

	try := 0
	for {
		var err error
		if _, err = w.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
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

func countIcingaHost(w *controller.Controller, objectList []*icinga.IcingaHost, expectZero bool) error {
	in := icingaHostSearchQuery(objectList)
	var respHost icinga.ResponseObject

	try := 0
	for {
		var err error
		if _, err = w.IcingaClient.Objects().Hosts("").Get([]string{}, in).Do().Into(&respHost); err != nil {
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

func GetIcingaHostList(w *controller.Controller, alert *aci.PodAlert) ([]*icinga.IcingaHost, error) {
	objectType, objectName := icinga.GetObjectInfo(alert.Labels)
	check := alert.Spec.Check

	// create all alerts for pod_status
	hostType, err := getIcingaHostType(check, objectType)
	if err != nil {
		return nil, err
	}
	objectList, err := icinga.GetObjectList(w.KubeClient, check, hostType, alert.Namespace, objectType, objectName, "")
	if err != nil {
		return nil, err
	}

	return objectList, nil
}*/
/*
func CheckIcingaObjectsForAlert(w *controller.Controller, alert *aci.PodAlert, expectZeroHost, expectZeroService bool) (err error) {
	objectList, err := GetIcingaHostList(w, alert)
	if err != nil {
		return err
	}

	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Host")
	if err = countIcingaHost(w, objectList, expectZeroHost); err != nil {
		return
	}

	// Count Icinga Service for 1st Alert. Should be found
	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)
	fmt.Println("----> Counting Icinga Service")
	if err = countIcingaService(w, objectList, serviceName, expectZeroService); err != nil {
		return
	}
	return
}

func CheckIcingaObjects(w *controller.Controller, alert *aci.PodAlert, objectList []*icinga.IcingaHost, expectZeroHost, expectZeroService bool) (err error) {
	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Host")
	if err = countIcingaHost(w, objectList, expectZeroHost); err != nil {
		return
	}

	// Count Icinga Service for 1st Alert. Should be found
	serviceName := strings.Replace(alert.Name, "_", "-", -1)
	serviceName = strings.Replace(serviceName, ".", "-", -1)
	fmt.Println("----> Counting Icinga Service")
	if err = countIcingaService(w, objectList, serviceName, expectZeroService); err != nil {
		return
	}
	return
}

func CheckIcingaObjectsForPod(w *controller.Controller, podName, namespace string, expectedService int32) error {
	// Count Icinga Host in Icinga2. Should be found
	fmt.Println("----> Counting Icinga Service")

	objectList := []*icinga.IcingaHost{
		{
			Hostname: fmt.Sprintf("%v@%v", podName, namespace),
		},
	}

	in := icingaHostSearchQuery(objectList)
	var respService icinga.ResponseObject

	try := 0
	for {
		var err error
		if _, err = w.IcingaClient.Objects().Service("").Get([]string{}, in).Do().Into(&respService); err != nil {
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
}*/
