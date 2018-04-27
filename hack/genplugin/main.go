package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/appscode/go/runtime"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/plugin"
)

func main() {
	pluginFolder := runtime.GOPath() + "/src/github.com/appscode/searchlight/hack/deploy"
	checkCommandFolder := runtime.GOPath() + "/src/github.com/appscode/searchlight/docs/examples/plugins/check-command"

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
		log.Fatal(err)
	}
	defer f.Close()

	for _, p := range plugins {
		ioutil.WriteFile(filepath.Join(checkCommandFolder, fmt.Sprintf("%s.conf", p.Name)), []byte(plugin.GenerateCheckCommand(p)), 0666)
		plugin.MarshallPlugin(f, p, "yaml")
	}
}
