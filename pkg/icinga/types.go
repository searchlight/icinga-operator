package icinga

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
