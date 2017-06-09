package watcher

import (
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/events"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
	ext "k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/util/wait"
)

func (w *Watcher) Namespace() {
	log.Debugln("watching", events.Namespace.String())
	lw := &cache.ListWatch{
		ListFunc:  NamespaceListFunc(w.KubeClient),
		WatchFunc: NamespaceWatchFunc(w.KubeClient),
	}
	_, controller := w.Cache(events.Namespace, &kapi.Namespace{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Pod() {
	log.Debugln("watching", events.Pod.String())
	lw := &cache.ListWatch{
		ListFunc:  PodListFunc(w.KubeClient),
		WatchFunc: PodWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Pod, &kapi.Pod{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.PodStore = cache.StoreToPodLister{indexer}
}

func (w *Watcher) Service() {
	log.Debugln("watching", events.Service.String())
	lw := &cache.ListWatch{
		ListFunc:  ServiceListFunc(w.KubeClient),
		WatchFunc: ServiceWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Service, &kapi.Service{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.ServiceStore = cache.StoreToServiceLister{indexer}
}

func (w *Watcher) RC() {
	log.Debugln("watching", events.RC.String())
	lw := &cache.ListWatch{
		ListFunc:  ReplicationControllerListFunc(w.KubeClient),
		WatchFunc: ReplicationControllerWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.RC, &kapi.ReplicationController{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.RcStore = cache.StoreToReplicationControllerLister{indexer}
}

func (w *Watcher) ReplicaSet() {
	log.Debugln("watching", events.ReplicaSet.String())
	lw := &cache.ListWatch{
		ListFunc:  ReplicaSetListFunc(w.KubeClient),
		WatchFunc: ReplicaSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.ReplicaSet, &ext.ReplicaSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.ReplicaSetStore = cache.StoreToReplicaSetLister{indexer}
}

func (w *Watcher) StatefulSet() {
	log.Debugln("watching", events.StatefulSet.String())
	lw := &cache.ListWatch{
		ListFunc:  StatefulSetListFunc(w.KubeClient),
		WatchFunc: StatefulSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.StatefulSet, &apps.StatefulSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.StatefulSetStore = cache.StoreToStatefulSetLister{indexer}
}

func (w *Watcher) DaemonSet() {
	log.Debugln("watching", events.DaemonSet.String())
	lw := &cache.ListWatch{
		ListFunc:  DaemonSetListFunc(w.KubeClient),
		WatchFunc: DaemonSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.DaemonSet, &ext.DaemonSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.DaemonSetStore = cache.StoreToDaemonSetLister{indexer}
}

func (w *Watcher) Endpoint() {
	log.Debugln("watching", events.Endpoint.String())
	lw := &cache.ListWatch{
		ListFunc:  EndpointListFunc(w.KubeClient),
		WatchFunc: EndpointWatchFunc(w.KubeClient),
	}
	store, controller := w.CacheStore(events.Endpoint, &kapi.Endpoints{}, lw)
	go controller.Run(wait.NeverStop)
	w.Storage.EndpointStore = cache.StoreToEndpointsLister{store}
}

func (w *Watcher) Node() {
	log.Debugln("watching", events.Node.String())
	lw := &cache.ListWatch{
		ListFunc:  NodeListFunc(w.KubeClient),
		WatchFunc: NodeWatchFunc(w.KubeClient),
	}
	_, controller := w.CacheStore(events.Node, &kapi.Node{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Alert() {
	log.Debugln("watching", events.Alert.String())
	lw := &cache.ListWatch{
		ListFunc:  AlertListFunc(w.ExtClient),
		WatchFunc: AlertWatchFunc(w.ExtClient),
	}
	_, controller := w.Cache(events.Alert, &aci.Alert{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) AlertEvent() {
	log.Debugln("watching", events.AlertEvent.String())
	lw := &cache.ListWatch{
		ListFunc:  AlertEventListFunc(w.KubeClient),
		WatchFunc: AlertEventWatchFunc(w.KubeClient),
	}
	_, controller := w.Cache(events.AlertEvent, &kapi.Event{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Deployment() {
	log.Debugln("watching", events.Deployments.String())
	lw := &cache.ListWatch{
		ListFunc:  DeploymentListFunc(w.KubeClient),
		WatchFunc: DeploymentWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Deployments, &ext.Deployment{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.DeploymentStore = cache.StoreToDeploymentLister{indexer}
}
