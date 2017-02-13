package app

import (
	aci "github.com/appscode/k8s-addons/api"
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/k8s-addons/pkg/stash"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
)

type Watcher struct {
	acw.Watcher

	IcingaClient *icinga.IcingaClient
}

func (watch *Watcher) Run() {
	watch.setup()
	watch.Service()
	watch.StatefulSet()
	watch.DaemonSet()
	watch.ReplicaSet()
	watch.RC()
	watch.Pod()
	watch.Alert()
	watch.AlertEvent()
	watch.Node()
	watch.Deployment()
}

func (w *Watcher) setup() {
	log.Infoln("Ensuring ThirdPartyResource")
	if err := w.ensureThirdPartyResource(); err != nil {
		log.Fatalln(err)
	}
	w.Watcher.Dispatch = w.Dispatch
	w.Storage = &stash.Storage{}
}

func (w *Watcher) ensureThirdPartyResource() error {
	resourceName := "alert" + "." + aci.V1beta1SchemeGroupVersion.Group

	_, err := w.Client.Extensions().ThirdPartyResources().Get(resourceName)
	if !errors.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: unversioned.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: kapi.ObjectMeta{
			Name: resourceName,
		},
		Versions: []extensions.APIVersion{
			{
				Name: aci.V1beta1SchemeGroupVersion.Version,
			},
		},
	}

	_, err = w.Client.Extensions().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (w *Watcher) Dispatch(e *events.Event) error {
	if e.Ignorable() {
		return nil
	}
	log.Debugln("Dispatching event with resource", e.ResourceType, "event", e.EventType)
	if e.ResourceType == events.Alert || e.ResourceType == events.Node ||
		e.ResourceType == events.Pod || e.ResourceType == events.Service || e.ResourceType == events.AlertEvent {
		return controller.New(w.Client, w.IcingaClient, w.AppsCodeExtensionClient, w.Storage).Handle(e)
	}
	return nil
}
