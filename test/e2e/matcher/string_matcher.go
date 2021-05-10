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
	"fmt"
	"regexp"
	"strings"

	"github.com/onsi/gomega/types"
)

func ReceiveNotification(expected string) types.GomegaMatcher {
	return &notificationMatcher{
		expected: expected,
	}
}

func ReceiveNotificationWithExp(expected string) types.GomegaMatcher {
	return &notificationMatcher{
		expected: strings.Replace(expected, "[", `\[`, -1),
	}
}

type notificationMatcher struct {
	expected string
}

func (matcher *notificationMatcher) Match(actual interface{}) (success bool, err error) {
	regexpExpected, err := regexp.Compile(matcher.expected)
	if err != nil {
		return false, err
	}
	return regexpExpected.MatchString(actual.(string)), nil
}

func (matcher *notificationMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Found notification message: %v", actual)
}

func (matcher *notificationMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Found notification message: %v", actual)
}
