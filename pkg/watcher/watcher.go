package watcher

import (
	"reflect"
	"sync"
	"time"

	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/events"
	"github.com/appscode/searchlight/pkg/stash"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/fields"
	"k8s.io/kubernetes/pkg/runtime"
)

type Watcher struct {
	KubeClient   clientset.Interface
	ExtClient    acs.ExtensionInterface
	IcingaClient *icinga.IcingaClient

	// sync time to sync the list.
	SyncPeriod time.Duration

	// lister store
	Storage *stash.Storage
	sync.Mutex
}

func (w *Watcher) Run() {
	w.setup()
	w.Service()
	w.StatefulSet()
	w.DaemonSet()
	w.ReplicaSet()
	w.RC()
	w.Pod()
	w.Alert()
	w.AlertEvent()
	w.Node()
	w.Deployment()
}

func (w *Watcher) setup() {
	log.Infoln("Ensuring ThirdPartyResource")
	if err := w.ensureThirdPartyResource(); err != nil {
		log.Fatalln(err)
	}
	w.Storage = &stash.Storage{}
}

func (w *Watcher) ensureThirdPartyResource() error {
	resourceName := "alert" + "." + aci.V1alpha1SchemeGroupVersion.Group

	_, err := w.KubeClient.Extensions().ThirdPartyResources().Get(resourceName)
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
				Name: aci.V1alpha1SchemeGroupVersion.Version,
			},
		},
	}

	_, err = w.KubeClient.Extensions().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (w *Watcher) Dispatch(e *events.Event) error {
	if e.Ignorable() {
		return nil
	}
	log.Debugln("Dispatching event with resource", e.ResourceType, "event", e.EventType)
	if e.ResourceType == events.Alert || e.ResourceType == events.Node ||
		e.ResourceType == events.Pod || e.ResourceType == events.Service || e.ResourceType == events.AlertEvent {
		return controller.New(w.KubeClient, w.IcingaClient, w.ExtClient, w.Storage).Handle(e)
	}
	return nil
}

func (w *Watcher) Cache(resource events.ObjectType, object runtime.Object, lw *cache.ListWatch) (cache.Store, *cache.Controller) {
	var listWatch *cache.ListWatch
	if lw != nil {
		listWatch = lw
	} else {
		listWatch = cache.NewListWatchFromClient(w.KubeClient.Core().RESTClient(), resource.String(), kapi.NamespaceAll, fields.Everything())
	}

	return cache.NewInformer(
		listWatch,
		object,
		w.SyncPeriod,
		eventHandlerFuncs(w),
	)
}

func (w *Watcher) CacheStore(resource events.ObjectType, object runtime.Object, lw *cache.ListWatch) (cache.Store, *cache.Controller) {
	if lw == nil {
		lw = cache.NewListWatchFromClient(w.KubeClient.Core().RESTClient(), resource.String(), kapi.NamespaceAll, fields.Everything())
	}

	return stash.NewInformerPopulated(
		lw,
		object,
		w.SyncPeriod,
		eventHandlerFuncs(w),
	)
}

func (w *Watcher) CacheIndexer(resource events.ObjectType, object runtime.Object, lw *cache.ListWatch, indexers cache.Indexers) (cache.Indexer, *cache.Controller) {
	if lw == nil {
		lw = cache.NewListWatchFromClient(w.KubeClient.Core().RESTClient(), resource.String(), kapi.NamespaceAll, fields.Everything())
	}
	if indexers == nil {
		indexers = cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}
	}

	return stash.NewIndexerInformerPopulated(
		lw,
		object,
		w.SyncPeriod,
		eventHandlerFuncs(w),
		indexers,
	)
}

func eventHandlerFuncs(k *Watcher) cache.ResourceEventHandlerFuncs {
	return cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			e := events.New(events.Added, obj)
			k.Dispatch(e)
		},
		DeleteFunc: func(obj interface{}) {
			e := events.New(events.Deleted, obj)
			k.Dispatch(e)
		},
		UpdateFunc: func(old, new interface{}) {
			if !reflect.DeepEqual(old, new) {
				e := events.New(events.Updated, old, new)
				k.Dispatch(e)
			}
		},
	}
}
