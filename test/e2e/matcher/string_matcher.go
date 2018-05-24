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
