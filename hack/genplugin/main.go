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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/plugin"

	"gomodules.xyz/runtime"
	"k8s.io/klog/v2"
)

func main() {
	pluginFolder := runtime.GOPath() + "/src/go.searchlight.dev/icinga-operator/hack/deploy"
	checkCommandFolder := runtime.GOPath() + "/src/go.searchlight.dev/icinga-operator/docs/examples/plugins/check-command"

	plugins := []*api.SearchlightPlugin{
		plugin.GetComponentStatusPlugin(),
		plugin.GetJsonPathPlugin(),
		plugin.GetNodeExistsPlugin(),
		plugin.GetPodExistsPlugin(),
		plugin.GetEventPlugin(),
		plugin.GetCACertPlugin(),
		plugin.GetCertPlugin(),
		plugin.GetNodeStatusPlugin(),
		plugin.GetNodeVolumePlugin(),
		plugin.GetPodStatusPlugin(),
		plugin.GetPodVolumePlugin(),
		plugin.GetPodExecPlugin(),
	}

	f, err := os.OpenFile(filepath.Join(pluginFolder, "plugins.yaml"), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		klog.Fatal(err)
	}
	defer f.Close()

	for _, p := range plugins {
		ioutil.WriteFile(filepath.Join(checkCommandFolder, fmt.Sprintf("%s.conf", p.Name)), []byte(plugin.GenerateCheckCommand(p)), 0666)
		plugin.MarshallPlugin(f, p, "yaml")
	}
}
