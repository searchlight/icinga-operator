package controller

import (
	"fmt"
	"reflect"
	"time"

	"github.com/appscode/errors"
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	acs "github.com/appscode/searchlight/client/clientset"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/client/icinga"
	"github.com/appscode/searchlight/pkg/controller/event"
	"github.com/appscode/searchlight/pkg/controller/host"
	_ "github.com/appscode/searchlight/pkg/controller/host/localhost"
	_ "github.com/appscode/searchlight/pkg/controller/host/node"
	_ "github.com/appscode/searchlight/pkg/controller/host/pod"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/events"
	"github.com/appscode/searchlight/pkg/stash"
	kapi "k8s.io/kubernetes/pkg/api"
	kerr "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/api/unversioned"
	clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	"k8s.io/kubernetes/pkg/labels"
)

type IcingaController struct {
	ctx *types.Context
}

func New(kubeClient clientset.Interface,
	icingaClient *icinga.IcingaClient,
	extClient acs.ExtensionInterface,
	storage *stash.Storage) *IcingaController {
	data, err := getIcingaDataMap()
	if err != nil {
		log.Errorln("Icinga data not found")
	}
	ctx := &types.Context{
		KubeClient:   kubeClient,
		ExtClient:    extClient,
		IcingaData:   data,
		IcingaClient: icingaClient,
		Storage:      storage,
	}
	return &IcingaController{ctx: ctx}
}

func (b *IcingaController) Handle(e *events.Event) error {
	var err error
	switch e.ResourceType {
	case events.Alert:
		err = b.handleAlert(e)
	case events.Pod:
		err = b.handlePod(e)
	case events.Node:
		err = b.handleNode(e)
	case events.Service:
		err = b.handleService(e)
	case events.AlertEvent:
		err = b.handleAlertEvent(e)
	}

	if err != nil {
		log.Errorln(err)
	}

	return nil
}

func (b *IcingaController) handleAlert(e *events.Event) error {
	alert := e.RuntimeObj

	if e.EventType.IsAdded() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		var err error
		_alert := alert[0].(*aci.Alert)
		if _alert.Status.CreationTime == nil {
			// Set Status
			t := unversioned.Now()
			_alert.Status.CreationTime = &t
			_alert.Status.Phase = aci.AlertPhaseCreating
			_alert, err = b.ctx.ExtClient.Alert(_alert.Namespace).Update(_alert)
			if err != nil {
				return errors.New().WithCause(err).Err()
			}
		}

		b.ctx.Resource = _alert

		if err := b.IsObjectExists(); err != nil {
			// Update Status
			t := unversioned.Now()
			_alert.Status.UpdateTime = &t
			_alert.Status.Phase = aci.AlertPhaseFailed
			_alert.Status.Reason = err.Error()
			if _, err := b.ctx.ExtClient.Alert(_alert.Namespace).Update(_alert); err != nil {
				return errors.New().WithCause(err).Err()
			}
			if kerr.IsNotFound(err) {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonNotFound, err.Error())
				return nil
			} else {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToProceed, err.Error())
				return errors.New().WithCause(err).Err()
			}
		}

		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonCreating)

		if err := b.Create(); err != nil {
			// Update Status
			t := unversioned.Now()
			_alert.Status.UpdateTime = &t
			_alert.Status.Phase = aci.AlertPhaseFailed
			_alert.Status.Reason = err.Error()
			if _, err := b.ctx.ExtClient.Alert(_alert.Namespace).Update(_alert); err != nil {
				return errors.New().WithCause(err).Err()
			}

			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToCreate, err.Error())
			return errors.New().WithCause(err).Err()
		}

		t := unversioned.Now()
		_alert.Status.UpdateTime = &t
		_alert.Status.Phase = aci.AlertPhaseCreated
		_alert.Status.Reason = ""
		if _, err = b.ctx.ExtClient.Alert(_alert.Namespace).Update(_alert); err != nil {
			return errors.New().WithCause(err).Err()
		}
		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulCreate)
	} else if e.EventType.IsUpdated() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		oldConfig := alert[0].(*aci.Alert)
		newConfig := alert[1].(*aci.Alert)

		if reflect.DeepEqual(oldConfig.Spec, newConfig.Spec) {
			return nil
		}

		if err := host.CheckAlertConfig(oldConfig, newConfig); err != nil {
			return errors.New().WithCause(err).Err()
		}

		b.ctx.Resource = newConfig

		if err := b.IsObjectExists(); err != nil {
			if kerr.IsNotFound(err) {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonNotFound, err.Error())
				return nil
			} else {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToProceed, err.Error())
				return errors.New().WithCause(err).Err()
			}
		}

		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonUpdating)

		if err := b.Update(); err != nil {
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToUpdate, err.Error())
			return errors.New().WithCause(err).Err()
		}

		// Set Status
		_alert := b.ctx.Resource
		t := unversioned.Now()
		_alert.Status.UpdateTime = &t
		if _, err := b.ctx.ExtClient.Alert(_alert.Namespace).Update(_alert); err != nil {
			return errors.New().WithCause(err).Err()
		}
		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulUpdate)
	} else if e.EventType.IsDeleted() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		b.ctx.Resource = alert[0].(*aci.Alert)
		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonDeleting)

		b.parseAlertOptions()
		if err := b.Delete(); err != nil {
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToDelete, err.Error())
			return errors.New().WithCause(err).Err()
		}
		event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulDelete)
	}
	return nil
}

func (b *IcingaController) handlePod(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}
	ancestors := b.getParentsForPod(e.RuntimeObj[0])
	if host.IsIcingaApp(ancestors, e.MetaData.Namespace) {
		if e.EventType.IsAdded() {
			go b.handleIcingaPod()
		}
	} else {
		return b.handleRegularPod(e, ancestors)
	}

	return nil
}

func (b *IcingaController) handleIcingaPod() {
	log.Debugln("Icinga pod is created...")
	then := time.Now()
	for {
		log.Debugln("Waiting for Icinga to UP")
		if b.checkIcingaAvailability() {
			break
		}
		now := time.Now()
		if now.Sub(then) > time.Minute*10 {
			log.Debugln("Icinga is down for more than 10 minutes..")
			return
		}
		time.Sleep(time.Second * 30)
	}

	icingaUp := false
	alertList, err := b.ctx.ExtClient.Alert(kapi.NamespaceAll).List(kapi.ListOptions{LabelSelector: labels.Everything()})
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, alert := range alertList.Items {
		if !icingaUp && !b.checkIcingaAvailability() {
			log.Debugln("Icinga is down...")
			return
		}
		icingaUp = true

		fakeEvent := &events.Event{
			ResourceType: events.Alert,
			EventType:    events.Added,
			RuntimeObj:   make([]interface{}, 0),
		}
		fakeEvent.RuntimeObj = append(fakeEvent.RuntimeObj, &alert)

		if err := b.handleAlert(fakeEvent); err != nil {
			log.Debugln(err)
		}
	}

	return
}

func (b *IcingaController) handleRegularPod(e *events.Event, ancestors []*types.Ancestors) error {
	namespace := e.MetaData.Namespace
	icingaUp := false
	ancestorItself := &types.Ancestors{
		Type:  events.Pod.String(),
		Names: []string{e.MetaData.Name},
	}
	ancestors = append(ancestors, ancestorItself)

	for _, ancestor := range ancestors {
		objectType := ancestor.Type
		for _, objectName := range ancestor.Names {
			lb, err := host.GetLabelSelector(objectType, objectName)
			if err != nil {
				return errors.New().WithCause(err).Err()
			}

			alertList, err := b.ctx.ExtClient.Alert(namespace).List(kapi.ListOptions{
				LabelSelector: lb,
			})
			if err != nil {
				return errors.New().WithCause(err).Err()
			}

			for _, alert := range alertList.Items {
				if !icingaUp && !b.checkIcingaAvailability() {
					return errors.New("Icinga is down").Err()
				}
				icingaUp = true

				if command, found := b.ctx.IcingaData[alert.Spec.CheckCommand]; found {
					if hostType, found := command.HostType[b.ctx.ObjectType]; found {
						if hostType != host.HostTypePod {
							continue
						}
					}
				}

				if e.EventType.IsAdded() {
					// Waiting for POD IP to use as Icinga Host IP
					then := time.Now()
					for {
						hasPodIP, err := b.checkPodIPAvailability(e.MetaData.Name, namespace)
						if err != nil {
							return errors.New().WithCause(err).Err()
						}
						if hasPodIP {
							break
						}
						log.Debugln("Waiting for pod IP")
						now := time.Now()
						if now.Sub(then) > time.Minute*2 {
							return errors.New("Pod IP is not available for 2 minutes").Err()
						}
						time.Sleep(time.Second * 30)
					}

					b.ctx.Resource = &alert

					additionalMessage := fmt.Sprintf(`pod "%v.%v"`, e.MetaData.Name, e.MetaData.Namespace)
					event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSync, additionalMessage)
					b.parseAlertOptions()

					if err := b.Create(e.MetaData.Name); err != nil {
						event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
						return errors.New().WithCause(err).Err()
					}
					event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulSync, additionalMessage)
				} else if e.EventType.IsDeleted() {
					b.ctx.Resource = &alert

					additionalMessage := fmt.Sprintf(`pod "%v.%v"`, e.MetaData.Name, e.MetaData.Namespace)
					event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSync, additionalMessage)
					b.parseAlertOptions()

					if err := b.Delete(e.MetaData.Name); err != nil {
						event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
						return errors.New().WithCause(err).Err()
					}
					event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulSync, additionalMessage)
				}

				t := unversioned.Now()
				alert.Status.UpdateTime = &t
				b.ctx.ExtClient.Alert(alert.Namespace).Update(&alert)
			}
		}
	}
	return nil
}

func (b *IcingaController) handleNode(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}

	lb, err := host.GetLabelSelector(events.Cluster.String(), "")
	if err != nil {
		return errors.New().WithCause(err).Err()
	}
	lb1, err := host.GetLabelSelector(events.Node.String(), e.MetaData.Name)
	if err != nil {
		return errors.New().WithCause(err).Err()
	}

	requirements, _ := lb1.Requirements()
	lb.Add(requirements...)

	icingaUp := false

	alertList, err := b.ctx.ExtClient.Alert(kapi.NamespaceAll).List(kapi.ListOptions{
		LabelSelector: lb,
	})
	if err != nil {
		return errors.New().WithCause(err).Err()
	}

	for _, alert := range alertList.Items {
		if !icingaUp && !b.checkIcingaAvailability() {
			return errors.New("Icinga is down").Err()
		}
		icingaUp = true

		if command, found := b.ctx.IcingaData[alert.Spec.CheckCommand]; found {
			if hostType, found := command.HostType[b.ctx.ObjectType]; found {
				if hostType != host.HostTypeNode {
					continue
				}
			}
		}

		if e.EventType.IsAdded() {
			b.ctx.Resource = &alert

			additionalMessage := fmt.Sprintf(`node "%v"`, e.MetaData.Name)
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSync, additionalMessage)
			b.parseAlertOptions()

			if err := b.Create(e.MetaData.Name); err != nil {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.New().WithCause(err).Err()
			}
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulSync, additionalMessage)

		} else if e.EventType.IsDeleted() {
			b.ctx.Resource = &alert

			additionalMessage := fmt.Sprintf(`node "%v"`, e.MetaData.Name)
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSync, additionalMessage)
			b.parseAlertOptions()

			if err := b.Delete(e.MetaData.Name); err != nil {
				event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.New().WithCause(err).Err()
			}
			event.CreateAlertEvent(b.ctx.KubeClient, b.ctx.Resource, types.EventReasonSuccessfulSync, additionalMessage)
		}

		t := unversioned.Now()
		alert.Status.UpdateTime = &t
		b.ctx.ExtClient.Alert(alert.Namespace).Update(&alert)
	}

	return nil
}

func (b *IcingaController) handleService(e *events.Event) error {
	if e.EventType.IsAdded() {
		if checkIcingaService(e.MetaData.Name, e.MetaData.Namespace) {
			service, err := b.ctx.KubeClient.Core().Services(e.MetaData.Namespace).Get(e.MetaData.Name)
			if err != nil {
				return errors.New().WithCause(err).Err()
			}
			endpoint := fmt.Sprintf("https://%v:5665/v1", service.Spec.ClusterIP)
			b.ctx.IcingaClient = b.ctx.IcingaClient.SetEndpoint(endpoint)
		}
	}
	return nil
}

func (b *IcingaController) handleAlertEvent(e *events.Event) error {
	var alertEvents []interface{}
	if e.ResourceType == events.AlertEvent {
		alertEvents = e.RuntimeObj
	}

	if e.EventType.IsAdded() {
		if len(alertEvents) == 0 {
			return errors.New("Missing event data").Err()
		}
		alertEvent := alertEvents[0].(*kapi.Event)

		if _, found := alertEvent.Annotations[types.AcknowledgeTimestamp]; found {
			return errors.New("Event is already handled").Err()
		}

		eventRefObjKind := alertEvent.InvolvedObject.Kind

		if eventRefObjKind != events.ObjectKindAlert.String() {
			return errors.New("For acknowledgement, Reference object should be Alert").Err()
		}

		eventRefObjNamespace := alertEvent.InvolvedObject.Namespace
		eventRefObjName := alertEvent.InvolvedObject.Name

		alert, err := b.ctx.ExtClient.Alert(eventRefObjNamespace).Get(eventRefObjName)
		if err != nil {
			return errors.New().WithCause(err).Err()
		}

		b.ctx.Resource = alert
		return b.Acknowledge(alertEvent)
	}
	return nil
}

func getIcingaDataMap() (map[string]*types.IcingaData, error) {
	icingaData, err := data.LoadIcingaData()
	if err != nil {
		return nil, errors.New().WithCause(err).Err()
	}

	icingaDataMap := make(map[string]*types.IcingaData)
	for _, command := range icingaData.Command {
		varsMap := make(map[string]data.CommandVar)
		for _, v := range command.Vars {
			varsMap[v.Name] = v
		}

		icingaDataMap[command.Name] = &types.IcingaData{
			HostType: command.ObjectToHost,
			VarInfo:  varsMap,
		}
	}
	return icingaDataMap, nil
}
