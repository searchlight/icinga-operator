package controller

import (
	"fmt"
	"reflect"
	"time"

	"github.com/appscode/errors"
	"github.com/appscode/kubed/pkg/events"
	"github.com/appscode/log"
	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/pkg/analytics"
	"github.com/appscode/searchlight/pkg/controller/eventer"
	"github.com/appscode/searchlight/pkg/controller/types"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (c *Controller) Handle(e *events.Event) error {
	var err error
	switch e.ResourceType {
	case events.Alert:
		err = c.handleAlert(e)
		sendEventForAlert(e.EventType, err)
	case events.Pod:
		err = c.handlePod(e)
	case events.Node:
		err = c.handleNode(e)
	case events.Service:
		err = c.handleService(e)
	case events.AlertEvent:
		err = c.handleAlertEvent(e)
	}

	if err != nil {
		log.Errorln(err)
	}

	return nil
}

func (c *Controller) handleAlert(e *events.Event) error {
	alert := e.RuntimeObj

	if e.EventType.IsAdded() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		var err error
		_alert := alert[0].(*aci.PodAlert)
		if _alert.Status.CreationTime == nil {
			// Set Status
			t := metav1.Now()
			_alert.Status.CreationTime = &t
			_alert.Status.Phase = aci.AlertPhaseCreating
			_alert, err = c.opt.ExtClient.PodAlerts(_alert.Namespace).Update(_alert)
			if err != nil {
				return errors.FromErr(err).Err()
			}
		}

		c.opt.Resource = _alert

		if err := c.IsObjectExists(); err != nil {
			// Update Status
			t := metav1.Now()
			_alert.Status.UpdateTime = &t
			_alert.Status.Phase = aci.AlertPhaseFailed
			_alert.Status.Reason = err.Error()
			if _, err := c.opt.ExtClient.PodAlerts(_alert.Namespace).Update(_alert); err != nil {
				return errors.FromErr(err).Err()
			}
			if kerr.IsNotFound(err) {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonNotFound, err.Error())
				return nil
			} else {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToProceed, err.Error())
				return errors.FromErr(err).Err()
			}
		}

		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonCreating)

		if err := c.Create(); err != nil {
			// Update Status
			t := metav1.Now()
			_alert.Status.UpdateTime = &t
			_alert.Status.Phase = aci.AlertPhaseFailed
			_alert.Status.Reason = err.Error()
			if _, err := c.opt.ExtClient.PodAlerts(_alert.Namespace).Update(_alert); err != nil {
				return errors.FromErr(err).Err()
			}

			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToCreate, err.Error())
			return errors.FromErr(err).Err()
		}

		t := metav1.Now()
		_alert.Status.UpdateTime = &t
		_alert.Status.Phase = aci.AlertPhaseCreated
		_alert.Status.Reason = ""
		if _, err = c.opt.ExtClient.PodAlerts(_alert.Namespace).Update(_alert); err != nil {
			return errors.FromErr(err).Err()
		}
		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulCreate)
	} else if e.EventType.IsUpdated() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		oldConfig := alert[0].(*aci.PodAlert)
		newConfig := alert[1].(*aci.PodAlert)

		if reflect.DeepEqual(oldConfig.Spec, newConfig.Spec) {
			return nil
		}

		if err := CheckAlertConfig(oldConfig, newConfig); err != nil {
			return errors.FromErr(err).Err()
		}

		c.opt.Resource = newConfig

		if err := c.IsObjectExists(); err != nil {
			if kerr.IsNotFound(err) {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonNotFound, err.Error())
				return nil
			} else {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToProceed, err.Error())
				return errors.FromErr(err).Err()
			}
		}

		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonUpdating)

		if err := c.Update(); err != nil {
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToUpdate, err.Error())
			return errors.FromErr(err).Err()
		}

		// Set Status
		_alert := c.opt.Resource
		t := metav1.Now()
		_alert.Status.UpdateTime = &t
		if _, err := c.opt.ExtClient.PodAlerts(_alert.Namespace).Update(_alert); err != nil {
			return errors.FromErr(err).Err()
		}
		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulUpdate)
	} else if e.EventType.IsDeleted() {
		if len(alert) == 0 {
			return errors.New("Missing alert data").Err()
		}

		c.opt.Resource = alert[0].(*aci.PodAlert)
		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonDeleting)

		c.parseAlertOptions()
		if err := c.Delete(); err != nil {
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToDelete, err.Error())
			return errors.FromErr(err).Err()
		}
		eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulDelete)
	}
	return nil
}

func (c *Controller) handlePod(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}
	ancestors := c.getParentsForPod(e.RuntimeObj[0])
	if IsIcingaApp(ancestors, e.MetaData.Namespace) {
		if e.EventType.IsAdded() {
			go c.handleIcingaPod()
		}
	} else {
		return c.handleRegularPod(e, ancestors)
	}

	return nil
}

func (c *Controller) handleIcingaPod() {
	log.Debugln("Icinga pod is created...")
	then := time.Now()
	for {
		log.Debugln("Waiting for Icinga to UP")
		if c.checkIcingaAvailability() {
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
	alertList, err := c.opt.ExtClient.PodAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: labels.Everything().String(),
	})
	if err != nil {
		log.Errorln(err)
		return
	}

	for _, alert := range alertList.Items {
		if !icingaUp && !c.checkIcingaAvailability() {
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

		if err := c.handleAlert(fakeEvent); err != nil {
			log.Debugln(err)
		}
	}

	return
}

func (c *Controller) handleRegularPod(e *events.Event, ancestors []*types.Ancestors) error {
	namespace := e.MetaData.Namespace
	icingaUp := false
	ancestorItself := &types.Ancestors{
		Type:  events.Pod.String(),
		Names: []string{e.MetaData.Name},
	}

	syncAlert := func(alert aci.PodAlert) error {
		if e.EventType.IsAdded() {
			// Waiting for POD IP to use as Icinga Host IP
			then := time.Now()
			for {
				hasPodIP, err := c.checkPodIPAvailability(e.MetaData.Name, namespace)
				if err != nil {
					return errors.FromErr(err).Err()
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

			c.opt.Resource = &alert

			additionalMessage := fmt.Sprintf(`pod "%v.%v"`, e.MetaData.Name, e.MetaData.Namespace)
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSync, additionalMessage)
			c.parseAlertOptions()

			if err := c.Create(e.MetaData.Name); err != nil {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.FromErr(err).Err()
			}
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulSync, additionalMessage)
		} else if e.EventType.IsDeleted() {
			c.opt.Resource = &alert

			additionalMessage := fmt.Sprintf(`pod "%v.%v"`, e.MetaData.Name, e.MetaData.Namespace)
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSync, additionalMessage)
			c.parseAlertOptions()

			if err := c.Delete(e.MetaData.Name); err != nil {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.FromErr(err).Err()
			}
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulSync, additionalMessage)
		}
		return nil
	}

	ancestors = append(ancestors, ancestorItself)
	for _, ancestor := range ancestors {
		objectType := ancestor.Type
		for _, objectName := range ancestor.Names {
			lb, err := GetLabelSelector(objectType, objectName)
			if err != nil {
				return errors.FromErr(err).Err()
			}

			alertList, err := c.opt.ExtClient.PodAlerts(namespace).List(metav1.ListOptions{
				LabelSelector: lb.String(),
			})
			if err != nil {
				return errors.FromErr(err).Err()
			}

			for _, alert := range alertList.Items {
				if !icingaUp && !c.checkIcingaAvailability() {
					return errors.New("Icinga is down").Err()
				}
				icingaUp = true

				if command, found := c.opt.IcingaData[alert.Spec.Check]; found {
					if hostType, found := command.HostType[c.opt.ObjectType]; found {
						if hostType != HostTypePod {
							continue
						}
					}
				}

				err = syncAlert(alert)
				sendEventForSync(e.EventType, err)

				if err != nil {
					return err
				}

				t := metav1.Now()
				alert.Status.UpdateTime = &t
				c.opt.ExtClient.PodAlerts(alert.Namespace).Update(&alert)
			}
		}
	}
	return nil
}

func (c *Controller) handleNode(e *events.Event) error {
	if !(e.EventType.IsAdded() || e.EventType.IsDeleted()) {
		return nil
	}

	lb, err := GetLabelSelector(events.Cluster.String(), "")
	if err != nil {
		return errors.FromErr(err).Err()
	}
	lb1, err := GetLabelSelector(events.Node.String(), e.MetaData.Name)
	if err != nil {
		return errors.FromErr(err).Err()
	}

	requirements, _ := lb1.Requirements()
	lb.Add(requirements...)

	icingaUp := false

	syncAlert := func(alert aci.PodAlert) error {
		if e.EventType.IsAdded() {
			c.opt.Resource = &alert

			additionalMessage := fmt.Sprintf(`node "%v"`, e.MetaData.Name)
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSync, additionalMessage)
			c.parseAlertOptions()

			if err := c.Create(e.MetaData.Name); err != nil {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.FromErr(err).Err()
			}
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulSync, additionalMessage)

		} else if e.EventType.IsDeleted() {
			c.opt.Resource = &alert

			additionalMessage := fmt.Sprintf(`node "%v"`, e.MetaData.Name)
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSync, additionalMessage)
			c.parseAlertOptions()

			if err := c.Delete(e.MetaData.Name); err != nil {
				eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonFailedToSync, additionalMessage, err.Error())
				return errors.FromErr(err).Err()
			}
			eventer.CreateAlertEvent(c.opt.KubeClient, c.opt.Resource, types.EventReasonSuccessfulSync, additionalMessage)
		}
		return nil
	}

	alertList, err := c.opt.ExtClient.PodAlerts(apiv1.NamespaceAll).List(metav1.ListOptions{
		LabelSelector: lb.String(),
	})
	if err != nil {
		return errors.FromErr(err).Err()
	}

	for _, alert := range alertList.Items {
		if !icingaUp && !c.checkIcingaAvailability() {
			return errors.New("Icinga is down").Err()
		}
		icingaUp = true

		if command, found := c.opt.IcingaData[alert.Spec.Check]; found {
			if hostType, found := command.HostType[c.opt.ObjectType]; found {
				if hostType != HostTypeNode {
					continue
				}
			}
		}

		err = syncAlert(alert)
		sendEventForSync(e.EventType, err)

		if err != nil {
			return err
		}

		t := metav1.Now()
		alert.Status.UpdateTime = &t
		c.opt.ExtClient.PodAlerts(alert.Namespace).Update(&alert)
	}

	return nil
}

func (c *Controller) handleService(e *events.Event) error {
	if e.EventType.IsAdded() {
		if checkIcingaService(e.MetaData.Name, e.MetaData.Namespace) {
			service, err := c.opt.KubeClient.CoreV1().Services(e.MetaData.Namespace).Get(e.MetaData.Name, metav1.GetOptions{})
			if err != nil {
				return errors.FromErr(err).Err()
			}
			endpoint := fmt.Sprintf("https://%v:5665/v1", service.Spec.ClusterIP)
			c.opt.IcingaClient = c.opt.IcingaClient.SetEndpoint(endpoint)
		}
	}
	return nil
}

func (c *Controller) handleAlertEvent(e *events.Event) error {
	var alertEvents []interface{}
	if e.ResourceType == events.AlertEvent {
		alertEvents = e.RuntimeObj
	}

	if e.EventType.IsAdded() {
		if len(alertEvents) == 0 {
			return errors.New("Missing event data").Err()
		}
		alertEvent := alertEvents[0].(*apiv1.Event)

		if _, found := alertEvent.Annotations[types.AcknowledgeTimestamp]; found {
			return errors.New("Event is already handled").Err()
		}

		eventRefObjKind := alertEvent.InvolvedObject.Kind

		if eventRefObjKind != events.ObjectKindAlert.String() {
			return errors.New("For acknowledgement, Reference object should be Alert").Err()
		}

		eventRefObjNamespace := alertEvent.InvolvedObject.Namespace
		eventRefObjName := alertEvent.InvolvedObject.Name

		alert, err := c.opt.ExtClient.PodAlerts(eventRefObjNamespace).Get(eventRefObjName)
		if err != nil {
			return errors.FromErr(err).Err()
		}

		c.opt.Resource = alert
		return c.Acknowledge(alertEvent)
	}
	return nil
}

func sendEventForAlert(eventType events.EventType, err error) {
	label := "success"
	if err != nil {
		label = "failure"
	}

	switch eventType {
	case events.Added:
		analytics.SendEvent("Alert", "created", label)
	case events.Updated:
		analytics.SendEvent("Alert", "updated", label)
	case events.Deleted:
		analytics.SendEvent("Alert", "deleted", label)
	}
}

func sendEventForSync(eventType events.EventType, err error) {
	label := "success"
	if err != nil {
		label = "failure"
	}

	switch eventType {
	case events.Added:
		analytics.SendEvent("Alert", "added", label)
	case events.Deleted:
		analytics.SendEvent("Alert", "removed", label)
	}
}
