package pod

import (
	"fmt"
	"testing"

	aci "github.com/appscode/searchlight/api"
	"github.com/appscode/searchlight/data"
	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/stretchr/testify/assert"
)

func TestSetParameterizedVariables(t *testing.T) {
	alertSpec := aci.AlertSpec{
		Vars: map[string]interface{}{
			"A": `Fake Query for pod_name   =  '?'`,
			"B": `Fake Query for pod_name='?'`,
		},
	}

	commandVars := map[string]data.CommandVar{
		"A": {
			Parameterized: true,
		},
		"B": {
			Parameterized: true,
		},
	}

	fakePodName := "test-pod"
	mp, err := setParameterizedVariables(alertSpec, fakePodName, commandVars, make(map[string]interface{}))
	assert.Nil(t, err)

	for key := range alertSpec.Vars {
		mpVal, found := mp[host.IVar(key)]
		if assert.True(t, found) {
			assert.EqualValues(t, mpVal, fmt.Sprintf(`Fake Query for pod_name='%v'`, fakePodName))
		}
	}

	alertSpec.Vars = map[string]interface{}{
		"A": `Fake Query for pod_name   =  '?'`,
		"C": `Fake Query for pod_name='?'`,
	}

	mp, err = setParameterizedVariables(alertSpec, fakePodName, commandVars, make(map[string]interface{}))
	assert.NotNil(t, err)

	alertSpec.Vars = map[string]interface{}{
		// Invalid Query. Should be (= '?')
		"A": `Fake Query for pod_name  = ?`,
	}

	mp, err = setParameterizedVariables(alertSpec, fakePodName, commandVars, make(map[string]interface{}))
	assert.Nil(t, err)

	for key := range alertSpec.Vars {
		mpVal, found := mp[host.IVar(key)]
		if assert.True(t, found) {
			assert.NotEqual(t, mpVal, fmt.Sprintf(`Fake Query for pod_name='%v'`, fakePodName))
		}
	}
}
