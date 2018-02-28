package icinga

import (
	"fmt"
	"os"
	"strings"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/clientset/versioned"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	internalIP = "InternalIP"

	TypePod     = "pod"
	TypeNode    = "node"
	TypeCluster = "cluster"
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
	return "", errors.Errorf("unknown host type %s", kh.Type)
}

func (kh IcingaHost) GetAlert(extClient cs.Interface, alertName string) (api.Alert, error) {
	switch kh.Type {
	case TypePod:
		return extClient.MonitoringV1alpha1().PodAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	case TypeNode:
		return extClient.MonitoringV1alpha1().NodeAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	case TypeCluster:
		return extClient.MonitoringV1alpha1().ClusterAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	}
	return nil, errors.Errorf("unknown host type %s", kh.Type)
}

func ParseHost(name string) (*IcingaHost, error) {
	parts := strings.SplitN(name, "@", 3)
	if !(len(parts) == 2 || len(parts) == 3) {
		return nil, errors.Errorf("host %s has a bad format", name)
	}
	t := parts[1]
	switch t {
	case TypePod, TypeNode:
		if len(parts) != 3 {
			return nil, errors.Errorf("host %s has a bad format", name)
		}
		return &IcingaHost{
			AlertNamespace: parts[0],
			Type:           t,
			ObjectName:     parts[2],
		}, nil
	case TypeCluster:
		if len(parts) != 2 {
			return nil, errors.Errorf("host %s has a bad format", name)
		}
		return &IcingaHost{
			AlertNamespace: parts[0],
			Type:           t,
		}, nil
	}
	return nil, errors.Errorf("unknown host type %s", t)
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
			LastState       float64                `json:"last_state"`
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
