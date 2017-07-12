package plugin

import "github.com/appscode/searchlight/pkg/icinga"

type TestData struct {
	Data                map[string]interface{}
	ExpectedIcingaState icinga.State
}
