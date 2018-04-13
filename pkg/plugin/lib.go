package plugin

import (
	"encoding/json"
	"fmt"
	"io"

	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/ghodss/yaml"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	stateOK       = "OK"
	stateWarning  = "Warning"
	stateCritical = "Critical"
	stateUnknown  = "Unknown"
)

// CustomResourceDefinitionTypeMeta set the default kind/apiversion of CRD
var PluginTypeMeta = metav1.TypeMeta{
	Kind:       "SearchlightPlugin",
	APIVersion: "monitoring.appscode.com/v1alpha1",
}

func MarshallPlugin(w io.Writer, plugin *api.SearchlightPlugin, outputFormat string) {
	jsonBytes, err := json.MarshalIndent(plugin, "", "    ")
	if err != nil {
		fmt.Println("error:", err)
	}

	if outputFormat == "json" {
		w.Write(jsonBytes)
	} else {
		yamlBytes, err := yaml.JSONToYAML(jsonBytes)
		if err != nil {
			fmt.Println("error:", err)
		}
		w.Write([]byte("---\n"))
		w.Write(yamlBytes)
	}
}
