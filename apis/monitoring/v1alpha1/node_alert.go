package v1alpha1

import (
	"fmt"
	"time"

	core "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourceKindNodeAlert = "NodeAlert"
	ResourceTypeNodeAlert = "nodealerts"
)

var _ Alert = &NodeAlert{}

func (a NodeAlert) GetName() string {
	return a.Name
}

func (a NodeAlert) GetNamespace() string {
	return a.Namespace
}

func (a NodeAlert) Command() string {
	return string(a.Spec.Check)
}

func (a NodeAlert) GetCheckInterval() time.Duration {
	return a.Spec.CheckInterval.Duration
}

func (a NodeAlert) GetAlertInterval() time.Duration {
	return a.Spec.AlertInterval.Duration
}

func (a NodeAlert) IsValid(kc kubernetes.Interface) error {
	if a.Spec.NodeName != nil && len(a.Spec.Selector) > 0 {
		return fmt.Errorf("can't specify both node name and selector")
	}

	cmd, ok := NodeCommands[a.Spec.Check]
	if !ok {
		return fmt.Errorf("%s is not a valid node check command", a.Spec.Check)
	}
	for k := range a.Spec.Vars {
		if _, ok := cmd.Vars[k]; !ok {
			return fmt.Errorf("var %s is unsupported for check command %s", k, a.Spec.Check)
		}
	}
	for _, rcv := range a.Spec.Receivers {
		found := false
		for _, state := range cmd.States {
			if state == rcv.State {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("state %s is unsupported for check command %s", rcv.State, a.Spec.Check)
		}
	}

	return checkNotifiers(kc, a)
}

func (a NodeAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a NodeAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a NodeAlert) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindNodeAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
