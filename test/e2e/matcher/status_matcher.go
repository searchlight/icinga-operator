package matcher

import (
	"fmt"

	"github.com/onsi/gomega/types"
	apiv1 "k8s.io/client-go/pkg/api/v1"
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
	podList := actual.(apiv1.PodList)

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
	return fmt.Sprintf("Expected\n\t%#v\nto contain sidecar container \n\t%#v", actual, matcher.expected)
}

func (matcher *statusMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain the sidecar container\n\t%#v", actual, matcher.expected)
}
