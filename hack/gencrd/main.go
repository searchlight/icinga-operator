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
	"k8s.io/apimachinery/pkg/apimachinery/announced"
	"k8s.io/apimachinery/pkg/apimachinery/registered"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
)

func generateCRDDefinitions() {
	filename := gort.GOPath() + "/src/github.com/appscode/searchlight/apis/monitoring/v1alpha1/crds.yaml"

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	crds := []*crd_api.CustomResourceDefinition{
		slitev1alpha1.ClusterAlert{}.CustomResourceDefinition(),
		slitev1alpha1.NodeAlert{}.CustomResourceDefinition(),
		slitev1alpha1.PodAlert{}.CustomResourceDefinition(),
		slitev1alpha1.Incident{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		crdutils.MarshallCrd(f, crd, "yaml")
	}
}

func generateSwaggerJson() {
	var (
		groupFactoryRegistry = make(announced.APIGroupFactoryRegistry)
		registry             = registered.NewOrDie("")
		Scheme               = runtime.NewScheme()
		Codecs               = serializer.NewCodecFactory(Scheme)
	)

	stashinstall.Install(groupFactoryRegistry, registry, Scheme)
	repoinstall.Install(groupFactoryRegistry, registry, Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Registry: registry,
		Scheme:   Scheme,
		Codecs:   Codecs,
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
		Resources: []schema.GroupVersionResource{
			slitev1alpha1.SchemeGroupVersion.WithResource(slitev1alpha1.ResourcePluralClusterAlert),
			slitev1alpha1.SchemeGroupVersion.WithResource(slitev1alpha1.ResourcePluralNodeAlert),
			slitev1alpha1.SchemeGroupVersion.WithResource(slitev1alpha1.ResourcePluralPodAlert),
			slitev1alpha1.SchemeGroupVersion.WithResource(slitev1alpha1.ResourcePluralIncident),
		},
		ImmutableResources: []schema.GroupVersionResource{
			incidentv1alpha1.SchemeGroupVersion.WithResource(incidentv1alpha1.ResourcePluralAcknowledgement),
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/github.com/appscode/searchlight/apis/swagger.json"
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	generateCRDDefinitions()
	generateSwaggerJson()
}
