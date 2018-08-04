package main

import (
	"io/ioutil"
	"os"

	"github.com/appscode/go/log"
	gort "github.com/appscode/go/runtime"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/openapi"
	repoinstall "github.com/appscode/searchlight/apis/incidents/install"
	incidentv1alpha1 "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	stashinstall "github.com/appscode/searchlight/apis/monitoring/install"
	slitev1alpha1 "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/go-openapi/spec"
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	"path/filepath"
)

func generateCRDDefinitions() {
	slitev1alpha1.EnableStatusSubresource = true

	filename := gort.GOPath() + "/src/github.com/appscode/searchlight/apis/monitoring/v1alpha1/crds.yaml"
	os.Remove(filename)

	err := os.MkdirAll(filepath.Join(gort.GOPath(), "/src/github.com/appscode/searchlight/api/crds"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	crds := []*crd_api.CustomResourceDefinition{
		slitev1alpha1.ClusterAlert{}.CustomResourceDefinition(),
		slitev1alpha1.NodeAlert{}.CustomResourceDefinition(),
		slitev1alpha1.PodAlert{}.CustomResourceDefinition(),
		slitev1alpha1.Incident{}.CustomResourceDefinition(),
		slitev1alpha1.SearchlightPlugin{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		filename := filepath.Join(gort.GOPath(), "/src/github.com/appscode/searchlight/api/crds", crd.Spec.Names.Singular+".yaml")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		crdutils.MarshallCrd(f, crd, "yaml")
		f.Close()
	}
}

func generateSwaggerJson() {
	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	stashinstall.Install(Scheme)
	repoinstall.Install(Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "stash-server",
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
			slitev1alpha1.GetOpenAPIDefinitions,
			incidentv1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []openapi.TypeInfo{
			{slitev1alpha1.SchemeGroupVersion, slitev1alpha1.ResourcePluralClusterAlert, slitev1alpha1.ResourceKindClusterAlert, true},
			{slitev1alpha1.SchemeGroupVersion, slitev1alpha1.ResourcePluralNodeAlert, slitev1alpha1.ResourceKindNodeAlert, true},
			{slitev1alpha1.SchemeGroupVersion, slitev1alpha1.ResourcePluralPodAlert, slitev1alpha1.ResourceKindPodAlert, true},
			{slitev1alpha1.SchemeGroupVersion, slitev1alpha1.ResourcePluralIncident, slitev1alpha1.ResourceKindIncident, true},
			{slitev1alpha1.SchemeGroupVersion, slitev1alpha1.ResourcePluralSearchlightPlugin, slitev1alpha1.ResourceKindSearchlightPlugin, true},
		},
		CDResources: []openapi.TypeInfo{
			{incidentv1alpha1.SchemeGroupVersion, incidentv1alpha1.ResourcePluralAcknowledgement, incidentv1alpha1.ResourceKindAcknowledgement, true},
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/github.com/appscode/searchlight/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	generateCRDDefinitions()
	generateSwaggerJson()
}
