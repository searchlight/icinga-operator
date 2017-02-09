package mini

import (
	"errors"
	"time"

	"github.com/appscode/k8s-addons/pkg/testing"
	"github.com/appscode/searchlight/cmd/searchlight/app"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/util"
	kapi "k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/apps"
)

func CreateStatefulSet(watcher *app.Watcher, namespace string) (*apps.StatefulSet, error) {
	// Create Service
	service, err := CreateService(watcher, namespace, nil)
	if err != nil {
		return nil, err
	}

	statefulSet := &apps.StatefulSet{
		ObjectMeta: kapi.ObjectMeta{
			Namespace: namespace,
		},
		Spec: apps.StatefulSetSpec{
			ServiceName: service.Name,
		},
	}

	if err := testing.CreateKubernetesObject(watcher.Client, statefulSet); err != nil {
		return nil, err
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		nStatefulSet, exists, err := watcher.Storage.StatefulSetStore.Get(statefulSet)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, errors.New("StatefulSet not found")
		}

		if nStatefulSet.(*apps.StatefulSet).Status.Replicas == statefulSet.Spec.Replicas {
			return nStatefulSet.(*apps.StatefulSet), nil
		}

		if check > 6 {
			return nil, errors.New("Fail to create StatefulSet")
		}
		check++
	}
}

func DeleteStatefulSet(watcher *app.Watcher, statefulSet *apps.StatefulSet) error {
	// Update StatefulSet
	statefulSet.Spec.Replicas = 0
	if _, err := watcher.Client.Apps().StatefulSets(statefulSet.Namespace).Update(statefulSet); err != nil {
		return err
	}

	labelSelector, err := util.GetLabels(watcher.Client, statefulSet.Namespace, host.TypeStatefulSet, statefulSet.Name)
	if err != nil {
		return err
	}

	check := 0
	for {
		time.Sleep(time.Second * 10)
		podList, err := watcher.Storage.PodStore.List(labelSelector)
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
	if err := watcher.Client.Apps().StatefulSets(statefulSet.Namespace).Delete(statefulSet.Name, nil); err != nil {
		return err
	}

	return DeleteService(watcher, &kapi.Service{
		ObjectMeta: kapi.ObjectMeta{
			Name:      statefulSet.Spec.ServiceName,
			Namespace: statefulSet.Namespace,
		},
	})
}
