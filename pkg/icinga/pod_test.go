package icinga

import (
	"fmt"
	"testing"

	tapi "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/data"
	"github.com/stretchr/testify/assert"
)

func TestSetParameterizedPodVariables(t *testing.T) {
	alertSpec := tapi.PodAlertSpec{
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
		mpVal, found := mp[IVar(key)]
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
		mpVal, found := mp[IVar(key)]
		if assert.True(t, found) {
			assert.NotEqual(t, mpVal, fmt.Sprintf(`Fake Query for pod_name='%v'`, fakePodName))
		}
	}
}
