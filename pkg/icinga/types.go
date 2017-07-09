package icinga

import (
	"fmt"
	"os"
)

const (
	internalIP = "InternalIP"

	TypePods   = "pods"
	TypeNodes  = "nodes"
	ObjectType = "alert.appscode.com/objectType"
	ObjectName = "alert.appscode.com/objectName"
)

type KHost struct {
	Name string
	IP   string
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
