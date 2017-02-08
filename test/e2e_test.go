package e2e

import (
	"fmt"
	"testing"

	"github.com/appscode/searchlight/pkg/controller/host"
	"github.com/appscode/searchlight/test/mini"
	"github.com/appscode/searchlight/test/util"
	"github.com/stretchr/testify/assert"
)

func TestMultipleAlerts(t *testing.T) {

	// Run KubeD
	// runKubeD(setIcingaClient bool)
	// Pass true to set IcingaClient in watcher
	watcher, err := runKubeD(true)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	// Create ReplicaSet
	fmt.Println("--> Creating ReplicaSet")
	replicaSet, err := mini.CreateReplicaSet(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println("--> Creating 1st Alert on ReplicaSet")
	labelMap := map[string]string{
		"objectType": host.TypeReplicasets,
		"objectName": replicaSet.Name,
	}
	firstAlert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, host.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := util.CheckIcingaObjects(watcher, firstAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	fmt.Println("--> Creating 2nd Alert on ReplicaSet")
	secondAlert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, host.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 2nd Alert")
	if err := util.CheckIcingaObjects(watcher, secondAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete 1st Alert
	fmt.Println("--> Deleting 1st Alert")
	if err := mini.DeleteAlert(watcher, firstAlert); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := util.CheckIcingaObjects(watcher, firstAlert, false, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete 2nd Alert
	fmt.Println("--> Deleting 2nd Alert")
	if err := mini.DeleteAlert(watcher, secondAlert); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 2nd Alert")
	if err := util.CheckIcingaObjects(watcher, secondAlert, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	if err := mini.DeleteReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}
}
