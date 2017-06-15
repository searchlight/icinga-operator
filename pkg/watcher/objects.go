package watcher

import (
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/events"
	"k8s.io/apimachinery/pkg/util/wait"
	apps_listers "k8s.io/client-go/listers/apps/v1beta1"
	core_listers "k8s.io/client-go/listers/core/v1"
	extensions_listers "k8s.io/client-go/listers/extensions/v1beta1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

func (w *Watcher) Namespace() {
	log.Debugln("watching", events.Namespace.String())
	lw := &cache.ListWatch{
		ListFunc:  NamespaceListFunc(w.KubeClient),
		WatchFunc: NamespaceWatchFunc(w.KubeClient),
	}
	_, controller := w.Cache(events.Namespace, &apiv1.Namespace{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Pod() {
	log.Debugln("watching", events.Pod.String())
	lw := &cache.ListWatch{
		ListFunc:  PodListFunc(w.KubeClient),
		WatchFunc: PodWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Pod, &apiv1.Pod{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.PodStore = core_listers.NewPodLister(indexer)
}

func (w *Watcher) Service() {
	log.Debugln("watching", events.Service.String())
	lw := &cache.ListWatch{
		ListFunc:  ServiceListFunc(w.KubeClient),
		WatchFunc: ServiceWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Service, &apiv1.Service{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.ServiceStore = core_listers.NewServiceLister(indexer)
}

func (w *Watcher) RC() {
	log.Debugln("watching", events.RC.String())
	lw := &cache.ListWatch{
		ListFunc:  ReplicationControllerListFunc(w.KubeClient),
		WatchFunc: ReplicationControllerWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.RC, &apiv1.ReplicationController{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.RcStore = core_listers.NewReplicationControllerLister(indexer)
}

func (w *Watcher) ReplicaSet() {
	log.Debugln("watching", events.ReplicaSet.String())
	lw := &cache.ListWatch{
		ListFunc:  ReplicaSetListFunc(w.KubeClient),
		WatchFunc: ReplicaSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.ReplicaSet, &extensions.ReplicaSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.ReplicaSetStore = extensions_listers.NewReplicaSetLister(indexer)
}

func (w *Watcher) StatefulSet() {
	log.Debugln("watching", events.StatefulSet.String())
	lw := &cache.ListWatch{
		ListFunc:  StatefulSetListFunc(w.KubeClient),
		WatchFunc: StatefulSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.StatefulSet, &apps.StatefulSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.StatefulSetStore = apps_listers.NewStatefulSetLister(indexer)
}

func (w *Watcher) DaemonSet() {
	log.Debugln("watching", events.DaemonSet.String())
	lw := &cache.ListWatch{
		ListFunc:  DaemonSetListFunc(w.KubeClient),
		WatchFunc: DaemonSetWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.DaemonSet, &extensions.DaemonSet{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.DaemonSetStore = extensions_listers.NewDaemonSetLister(indexer)
}

func (w *Watcher) Endpoint() {
	log.Debugln("watching", events.Endpoint.String())
	lw := &cache.ListWatch{
		ListFunc:  EndpointListFunc(w.KubeClient),
		WatchFunc: EndpointWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Endpoint, &apiv1.Endpoints{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.EndpointStore = core_listers.NewEndpointsLister(indexer)
}

func (w *Watcher) Node() {
	log.Debugln("watching", events.Node.String())
	lw := &cache.ListWatch{
		ListFunc:  NodeListFunc(w.KubeClient),
		WatchFunc: NodeWatchFunc(w.KubeClient),
	}
	_, controller := w.CacheStore(events.Node, &apiv1.Node{}, lw)
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
	_, controller := w.Cache(events.AlertEvent, &apiv1.Event{}, lw)
	go controller.Run(wait.NeverStop)
}

func (w *Watcher) Deployment() {
	log.Debugln("watching", events.Deployments.String())
	lw := &cache.ListWatch{
		ListFunc:  DeploymentListFunc(w.KubeClient),
		WatchFunc: DeploymentWatchFunc(w.KubeClient),
	}
	indexer, controller := w.CacheIndexer(events.Deployments, &extensions.Deployment{}, lw, nil)
	go controller.Run(wait.NeverStop)
	w.Storage.DeploymentStore = extensions_listers.NewDeploymentLister(indexer)
}
