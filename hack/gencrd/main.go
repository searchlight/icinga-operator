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
	"io/ioutil"
	"os"
	"path/filepath"

	incidentinstall "go.searchlight.dev/icinga-operator/apis/incidents/install"
	incidentv1alpha1 "go.searchlight.dev/icinga-operator/apis/incidents/v1alpha1"
	moninstall "go.searchlight.dev/icinga-operator/apis/monitoring/install"
	monv1alpha1 "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"

	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	gort "gomodules.xyz/runtime"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	"kmodules.xyz/client-go/openapi"
)

func main() {
	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	moninstall.Install(Scheme)
	incidentinstall.Install(Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "icinga-operator",
			Version: "v0",
			Contact: &spec.ContactInfo{
				Name:  "AppsCode Inc.",
				URL:   "https://appscode.com",
				Email: "hello@appscode.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			monv1alpha1.GetOpenAPIDefinitions,
			incidentv1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []openapi.TypeInfo{
			{monv1alpha1.SchemeGroupVersion, monv1alpha1.ResourcePluralClusterAlert, monv1alpha1.ResourceKindClusterAlert, true},
			{monv1alpha1.SchemeGroupVersion, monv1alpha1.ResourcePluralNodeAlert, monv1alpha1.ResourceKindNodeAlert, true},
			{monv1alpha1.SchemeGroupVersion, monv1alpha1.ResourcePluralPodAlert, monv1alpha1.ResourceKindPodAlert, true},
			{monv1alpha1.SchemeGroupVersion, monv1alpha1.ResourcePluralIncident, monv1alpha1.ResourceKindIncident, true},
			{monv1alpha1.SchemeGroupVersion, monv1alpha1.ResourcePluralSearchlightPlugin, monv1alpha1.ResourceKindSearchlightPlugin, true},
		},
		CDResources: []openapi.TypeInfo{
			{incidentv1alpha1.SchemeGroupVersion, incidentv1alpha1.ResourcePluralAcknowledgement, incidentv1alpha1.ResourceKindAcknowledgement, true},
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/go.searchlight.dev/icinga-operator/openapi/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}
