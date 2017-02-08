package plugin

import (
	"errors"
	"reflect"
	"strings"

	"github.com/appscode/searchlight/pkg/controller/host"
)

func GetKubeObjectInfo(hostname string) (objectType string, objectName string, namespace string, err error) {
	parts := strings.Split(hostname, "@")
	if len(parts) != 2 {
		err = errors.New("Invalid icinga host.name")
		return
	}
	name := parts[0]
	namespace = parts[1]

	objectType = ""
	objectName = ""
	if name != host.CheckCommandPodExists && name != host.CheckCommandPodStatus {
		parts = strings.Split(name, "|")
		if len(parts) == 1 {
			objectType = host.TypePods
			objectName = parts[0]
		} else if len(parts) == 2 {
			objectType = parts[0]
			objectName = parts[1]
		} else {
			err = errors.New("Invalid icinga host.name")
			return
		}
	}
	return
}

func FillStruct(data map[string]interface{}, result interface{}) {
	t := reflect.ValueOf(result).Elem()
	for k, v := range data {
		val := t.FieldByName(k)
		val.Set(reflect.ValueOf(v))
	}
}
