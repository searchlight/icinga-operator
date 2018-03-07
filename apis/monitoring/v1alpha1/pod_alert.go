package v1alpha1

import (
	"fmt"
	"time"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourceKindPodAlert = "PodAlert"
	ResourceTypePodAlert = "podalerts"
)

var _ Alert = &PodAlert{}

func (a PodAlert) GetName() string {
	return a.Name
}

func (a PodAlert) GetNamespace() string {
	return a.Namespace
}

func (a PodAlert) Command() string {
	return string(a.Spec.Check)
}

func (a PodAlert) GetCheckInterval() time.Duration {
	return a.Spec.CheckInterval.Duration
}

func (a PodAlert) GetAlertInterval() time.Duration {
	return a.Spec.AlertInterval.Duration
}

func (a PodAlert) IsValid(kc kubernetes.Interface) error {
	if a.Spec.PodName != nil && a.Spec.Selector != nil {
		return fmt.Errorf("can't specify both pod name and selector")
	}
	if a.Spec.PodName == nil && a.Spec.Selector == nil {
		return fmt.Errorf("specify either pod name or selector")
	}
	if a.Spec.Selector != nil {
		_, err := metav1.LabelSelectorAsSelector(a.Spec.Selector)
		if err != nil {
			return err
		}
	}

	cmd, ok := PodCommands[a.Spec.Check]
	if !ok {
		return fmt.Errorf("%s is not a valid pod check command", a.Spec.Check)
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

func (a PodAlert) GetNotifierSecretName() string {
	return a.Spec.NotifierSecretName
}

func (a PodAlert) GetReceivers() []Receiver {
	return a.Spec.Receivers
}

func (a PodAlert) ObjectReference() *core.ObjectReference {
	return &core.ObjectReference{
		APIVersion:      SchemeGroupVersion.String(),
		Kind:            ResourceKindPodAlert,
		Namespace:       a.Namespace,
		Name:            a.Name,
		UID:             a.UID,
		ResourceVersion: a.ResourceVersion,
	}
}
