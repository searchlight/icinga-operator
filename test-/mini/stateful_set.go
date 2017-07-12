package mini

import (
	"errors"
	"time"

	"github.com/appscode/go/types"
	"github.com/appscode/searchlight/pkg/controller"
	"github.com/appscode/searchlight/pkg/icinga"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
)

func CreateStatefulSet(w *controller.Controller, namespace string) (*apps.StatefulSet, error) {
	// Create Service
	service, err := CreateService(w, namespace, nil)
	if err != nil {
		return nil, err
	}

	statefulSet := &apps.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
		},
		Spec: apps.StatefulSetSpec{
			ServiceName: service.Name,
		},
	}

	if err := CreateKubernetesObject(w.KubeClient, statefulSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		ss, err := w.Storage.StatefulSetStore.StatefulSets(statefulSet.Namespace).Get(statefulSet.Name)
		if err != nil {
			return nil, err
		}
		if ss.Status.Replicas == *statefulSet.Spec.Replicas {
			return ss, nil
		}
		if check > 6 {
			return nil, errors.New("Fail to create StatefulSet")
		}
		check++
	}
}

func DeleteStatefulSet(w *controller.Controller, statefulSet *apps.StatefulSet) error {
	statefulSet, err := w.KubeClient.AppsV1beta1().StatefulSets(statefulSet.Namespace).Get(statefulSet.Name, metav1.GetOptions{})
	if err != nil {
		return err
	}
	// Update StatefulSet
	statefulSet.Spec.Replicas = types.Int32P(0)
	if _, err := w.KubeClient.AppsV1beta1().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		return err
	}

	labelSelector, err := icinga.GetLabels(w.KubeClient, statefulSet.Namespace, icinga.TypeStatefulSet, statefulSet.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 30)
		podList, err := w.Storage.PodStore.List(labelSelector)
		if err != nil {
			return err
		}
		if len(podList) == 0 {
			break
		}

		if check > 6 {
			return errors.New("Fail to delete StatefulSet Pods")
		}
		check++
	}

	// Delete StatefulSet
	if err := w.KubeClient.AppsV1beta1().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil); err != nil {
		return err
	}

	return DeleteService(w, &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSet.Spec.ServiceName,
			Namespace: statefulSet.Namespace,
		},
	})
}
