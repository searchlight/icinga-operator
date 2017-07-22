package framework

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

func (f *Framework) CountNode() (int32, error) {
	nodeList, err := f.kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, err
	}
	return int32(len(nodeList.Items)), nil
}
