package icinga

import (
	"fmt"
	"os"
	"strings"

	tapi "github.com/appscode/searchlight/api"
	tcs "github.com/appscode/searchlight/client/clientset"
)

const (
	internalIP = "InternalIP"

	TypePod     = "pod"
	TypeNode    = "node"
	TypeCluster = "cluster"

	ObjectType = "alert.appscode.com/objectType"
	ObjectName = "alert.appscode.com/objectName"
)

type IcingaHost struct {
	Type           string
	AlertNamespace string
	ObjectName     string
	IP             string
}

func (kh IcingaHost) Name() (string, error) {
	switch kh.Type {
	case TypePod:
		return kh.AlertNamespace + "@" + kh.Type + "@" + kh.ObjectName, nil
	case TypeNode:
		return kh.AlertNamespace + "@" + kh.Type + "@" + kh.ObjectName, nil
	case TypeCluster:
		return kh.AlertNamespace + "@" + kh.Type, nil
	}
	return "", fmt.Errorf("Unknown host type %s", kh.Type)
}

func (kh IcingaHost) GetAlert(extClient tcs.ExtensionInterface, alertName string) (tapi.Alert, error) {
	switch kh.Type {
	case TypePod:
		return extClient.PodAlerts(kh.AlertNamespace).Get(alertName)
	case TypeNode:
		return extClient.NodeAlerts(kh.AlertNamespace).Get(alertName)
	case TypeCluster:
		return extClient.ClusterAlerts(kh.AlertNamespace).Get(alertName)
	}
	return nil, fmt.Errorf("Unknown host type %s", kh.Type)
}

func ParseHost(name string) (*IcingaHost, error) {
	parts := strings.SplitN(name, "@", 3)
	if len(parts) == 2 || len(parts) == 3 {
		return nil, fmt.Errorf("Host %s has a bad format", name)
	}
	t := parts[1]
	switch t {
	case TypePod, TypeNode:
		if len(parts) != 3 {
			return nil, fmt.Errorf("Host %s has a bad format", name)
		}
		return &IcingaHost{
			AlertNamespace: parts[0],
			Type:           t,
			ObjectName:     parts[2],
		}, nil
	case TypeCluster:
		if len(parts) != 2 {
			return nil, fmt.Errorf("Host %s has a bad format", name)
		}
		return &IcingaHost{
			AlertNamespace: parts[0],
			Type:           t,
		}, nil
	}
	return nil, fmt.Errorf("Unknown host type %s", t)
}

type IcingaObject struct {
	Templates []string               `json:"templates,omitempty"`
	Attrs     map[string]interface{} `json:"attrs"`
}

type ResponseObject struct {
	Results []struct {
		Attrs struct {
			Name            string                 `json:"name"`
			CheckInterval   float64                `json:"check_interval"`
			Vars            map[string]interface{} `json:"vars"`
			Acknowledgement float64                `json:"acknowledgement"`
		} `json:"attrs"`
		Name string `json:"name"`
	} `json:"results"`
}

func IVar(value string) string {
	return "vars." + value
}

type State int32

const (
	OK       State = iota // 0
	WARNING               // 1
	CRITICAL              // 2
	UNKNOWN               // 3
)

func Output(s State, message interface{}) {
	fmt.Fprintln(os.Stdout, s, ":", message)
	os.Exit(int(s))
}
