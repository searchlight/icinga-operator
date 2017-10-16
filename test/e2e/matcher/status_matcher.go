package matcher

import (
	"github.com/onsi/gomega/types"
	apiv1 "k8s.io/api/core/v1"
)

func HaveRunningPods(expected int32) types.GomegaMatcher {
	return &statusMatcher{
		expected: expected,
	}
}

type statusMatcher struct {
	expected int32
}

func (matcher *statusMatcher) Match(actual interface{}) (success bool, err error) {
	podList := actual.(*apiv1.PodList)
	if int(matcher.expected) != len(podList.Items) {
		return false, nil
	}
	for _, pod := range podList.Items {
		if pod.Status.Phase != apiv1.PodRunning {
			return false, nil
		}
	}
	return true, nil
}

func (matcher *statusMatcher) FailureMessage(actual interface{}) (message string) {
	return "Expected to be Running all Pods"
}

func (matcher *statusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "Expected to be not Running all Pods"
}

func HavePods(expected int32) types.GomegaMatcher {
	return &countMatcher{
		expected: expected,
	}
}

type countMatcher struct {
	expected int32
}

func (matcher *countMatcher) Match(actual interface{}) (success bool, err error) {
	podList := actual.(*apiv1.PodList)
	if int(matcher.expected) != len(podList.Items) {
		return false, nil
	}
	return true, nil
}

func (matcher *countMatcher) FailureMessage(actual interface{}) (message string) {
	return "Expected to have all Pods"
}

func (matcher *countMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return "Expected to not have all Pods"
}
