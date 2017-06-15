package controller

import (
	"os"

	"github.com/appscode/errors"
	"github.com/appscode/log"
	"github.com/appscode/searchlight/pkg/controller/types"
	"github.com/appscode/searchlight/pkg/events"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
)

func (b *IcingaController) IsObjectExists() error {
	log.Infoln("Checking Kubernetes Object existence", b.ctx.Resource.ObjectMeta)
	b.parseAlertOptions()

	var err error
	switch b.ctx.ObjectType {
	case events.Service.String():
		_, err = b.ctx.KubeClient.CoreV1().Services(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.RC.String():
		_, err = b.ctx.KubeClient.CoreV1().ReplicationControllers(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.DaemonSet.String():
		_, err = b.ctx.KubeClient.ExtensionsV1beta1().DaemonSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.Deployments.String():
		_, err = b.ctx.KubeClient.ExtensionsV1beta1().Deployments(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.StatefulSet.String():
		_, err = b.ctx.KubeClient.AppsV1beta1().StatefulSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.ReplicaSet.String():
		_, err = b.ctx.KubeClient.ExtensionsV1beta1().ReplicaSets(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.Pod.String():
		_, err = b.ctx.KubeClient.CoreV1().Pods(b.ctx.Resource.Namespace).Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.Node.String():
		_, err = b.ctx.KubeClient.CoreV1().Nodes().Get(b.ctx.ObjectName, metav1.GetOptions{})
	case events.Cluster.String():
		err = nil
	default:
		err = errors.Newf(`Invalid Object Type "%s"`, b.ctx.ObjectType).Err()
	}
	return err
}

func (b *IcingaController) getParentsForPod(o interface{}) []*types.Ancestors {
	pod := o.(*apiv1.Pod)
	result := make([]*types.Ancestors, 0)

	svc, err := b.ctx.Storage.ServiceStore.GetPodServices(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range svc {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.Service.String(),
			Names: names,
		})
	}

	rc, err := b.ctx.Storage.RcStore.GetPodControllers(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range rc {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.RC.String(),
			Names: names,
		})
	}

	rs, err := b.ctx.Storage.ReplicaSetStore.GetPodReplicaSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range rs {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.ReplicaSet.String(),
			Names: names,
		})
	}

	ps, err := b.ctx.Storage.StatefulSetStore.GetPodStatefulSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range ps {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.StatefulSet.String(),
			Names: names,
		})
	}

	ds, err := b.ctx.Storage.DaemonSetStore.GetPodDaemonSets(pod)
	if err == nil {
		names := make([]string, 0)
		for _, s := range ds {
			names = append(names, s.Name)
		}
		result = append(result, &types.Ancestors{
			Type:  events.DaemonSet.String(),
			Names: names,
		})
	}
	return result
}

func (b *IcingaController) checkIcingaAvailability() bool {
	log.Debugln("Checking Icinga client")
	if b.ctx.IcingaClient == nil {
		return false
	}
	resp := b.ctx.IcingaClient.Check().Get([]string{}).Do()
	if resp.Status != 200 {
		return false
	}
	return true
}

func (b *IcingaController) checkPodIPAvailability(podName, namespace string) (bool, error) {
	log.Debugln("Checking pod IP")
	pod, err := b.ctx.KubeClient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return false, errors.New().WithCause(err).Err()
	}
	if pod.Status.PodIP == "" {
		return false, nil
	}
	return true, nil
}

func checkIcingaService(serviceName, namespace string) bool {
	icingaService := os.Getenv("ICINGA_SERVICE_NAME")
	if serviceName != icingaService {
		return false
	}
	icingaServiceNamespace := os.Getenv("ICINGA_SERVICE_NAMESPACE")
	if namespace != icingaServiceNamespace {
		return false
	}
	return true
}
