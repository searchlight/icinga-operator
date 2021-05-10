/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"encoding/json"
	"fmt"
	"io"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

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
