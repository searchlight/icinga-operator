package plugin

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/appscode/searchlight/pkg/controller/host"
	"reflect"
)

func GetKubeObjectInfo(hostname string) (objectType string, objectName string, namespace string) {
	parts := strings.Split(hostname, "@")
	if len(parts) != 2 {
		fmt.Println(errors.New("Invalid icinga host.name"))
		os.Exit(1)
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
			fmt.Println(errors.New("Invalid icinga host.name"))
			os.Exit(1)
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