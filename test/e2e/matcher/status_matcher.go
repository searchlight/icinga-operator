/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package matcher

import (
	"github.com/onsi/gomega/types"
	core "k8s.io/api/core/v1"
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
	podList := actual.(*core.PodList)
	if int(matcher.expected) != len(podList.Items) {
		return false, nil
	}
	for _, pod := range podList.Items {
		if pod.Status.Phase != core.PodRunning {
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
	podList := actual.(*core.PodList)
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
