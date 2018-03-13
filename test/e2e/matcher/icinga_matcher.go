package matcher

import (
	"fmt"

	"github.com/onsi/gomega/types"
)

type IcingaServiceState struct {
	OK       int32
	Warning  int32
	Critical int32
	Unknown  int32
}

func HaveIcingaObject(expected IcingaServiceState) types.GomegaMatcher {
	return &icingaObjectMatcher{
		expected: expected,
	}
}

type icingaObjectMatcher struct {
	expected IcingaServiceState
}

func (matcher *icingaObjectMatcher) Match(actual interface{}) (success bool, err error) {
	switch obj := actual.(type) {
	case IcingaServiceState:
		if obj.OK != matcher.expected.OK {
			return false, nil
		}
		if obj.Warning != matcher.expected.Warning {
			return false, nil
		}
		if obj.Critical != matcher.expected.Critical {
			return false, nil
		}
		if obj.Unknown != matcher.expected.Unknown {
			return false, nil
		}
	}
	return true, nil
}

func (matcher *icingaObjectMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Found icinga object: %v", actual)
}

func (matcher *icingaObjectMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Found icinga object: %v", actual)
}
