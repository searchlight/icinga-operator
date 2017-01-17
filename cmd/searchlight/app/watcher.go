package app

import (
	"github.com/appscode/k8s-addons/pkg/events"
	"github.com/appscode/k8s-addons/pkg/stash"
	acw "github.com/appscode/k8s-addons/pkg/watcher"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller"
)

type Watcher struct {
	acw.Watcher

	IcingaClient *icinga.IcingaClient
}

func (watch *Watcher) Run() {
	watch.Watcher.Dispatch = watch.Dispatch
	watch.Storage = &stash.Storage{}
	watch.Service()
	watch.StatefulSet()
	watch.DaemonSet()
	watch.ReplicaSet()
	watch.RC()
	watch.Pod()
	watch.Alert()
	watch.AlertEvent()
	watch.Node()
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
