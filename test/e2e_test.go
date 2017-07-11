package e2e

import (
	"fmt"
	"testing"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/test/mini"
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
		"objectType": icinga.TypeReplicasets,
		"objectName": replicaSet.Name,
	}
	firstAlert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, firstAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	fmt.Println("--> Creating 2nd Alert on ReplicaSet")
	secondAlert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 2nd Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, secondAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Increment Replica
	fmt.Println("--> Incrementing Replica")
	replicaSet.Spec.Replicas++
	if replicaSet, err = mini.UpdateReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}

	// Get Last Replica
	fmt.Println("--> Getting Last Replica")
	lastPod, err := mini.GetLastReplica(watcher, replicaSet)
	if !assert.Nil(t, err) {
		return
	}

	// Checking Icinga Objects for This Pod
	fmt.Println("----> Checking Icinga Objects for last Pod")
	if err = icinga.CheckIcingaObjectsForPod(watcher, lastPod.Name, lastPod.Namespace, 2); !assert.Nil(t, err) {
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
	if err := icinga.CheckIcingaObjectsForAlert(watcher, firstAlert, false, true); !assert.Nil(t, err) {
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
	if err := icinga.CheckIcingaObjectsForAlert(watcher, secondAlert, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete ReplicaSet
	fmt.Println("--> Deleting ReplicaSet")
	if err := mini.DeleteReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}
}

func TestMultipleAlertsOnMultipleObjects(t *testing.T) {
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

	// Create Service on ReplicaSet
	fmt.Println("--> Creating Service")
	service, err := mini.CreateService(watcher, replicaSet.Namespace, replicaSet.Spec.Template.Labels)
	if !assert.Nil(t, err) {
		return
	}

	// Create 1st Alert
	fmt.Println("--> Creating 1st Alert on ReplicaSet")
	labelMap := map[string]string{
		"objectType": icinga.TypeReplicasets,
		"objectName": replicaSet.Name,
	}
	firstAlert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, firstAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Create 2nd Alert
	fmt.Println("--> Creating 2nd Alert on Service")
	labelMap = map[string]string{
		"objectType": icinga.TypeServices,
		"objectName": service.Name,
	}
	secondAlert, err := mini.CreateAlert(watcher, service.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 2nd Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, secondAlert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Get Pod
	fmt.Println("--> Getting Pod")
	pod, err := mini.GetLastReplica(watcher, replicaSet)
	if !assert.Nil(t, err) {
		return
	}

	// Checking Icinga Objects for Pod
	fmt.Println("----> Checking Icinga Objects for Pod")
	if err = icinga.CheckIcingaObjectsForPod(watcher, pod.Name, pod.Namespace, 2); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete Pod
	fmt.Println("--> Deleting Pod")
	if err := mini.DeletePod(watcher, pod); !assert.Nil(t, err) {
		return
	}

	// Checking Icinga Objects for Pod
	fmt.Println("----> Checking Icinga Objects for Pod")
	if err = icinga.CheckIcingaObjectsForPod(watcher, pod.Name, pod.Namespace, 0); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Getting ReplicaSetObjectList
	replicaSetObjectList, err := icinga.GetIcingaHostList(watcher, firstAlert)
	if !assert.Nil(t, err) {
		return
	}

	// Delete ReplicaSet
	fmt.Println("--> Deleting ReplicaSet")
	if err := mini.DeleteReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := icinga.CheckIcingaObjects(watcher, firstAlert, replicaSetObjectList, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete Service
	fmt.Println("--> Deleting Service")
	if err := mini.DeleteService(watcher, service); !assert.Nil(t, err) {
		return
	}

	// Delete 1st Alert
	fmt.Println("--> Deleting 1st Alert")
	if err := mini.DeleteAlert(watcher, firstAlert); !assert.Nil(t, err) {
		return
	}

	// Delete 2nd Alert
	fmt.Println("--> Deleting 2nd Alert")
	if err := mini.DeleteAlert(watcher, secondAlert); !assert.Nil(t, err) {
		return
	}
}

func TestAlertWhileReCreateKubeObject(t *testing.T) {
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

	fmt.Println("--> Creating Alert on ReplicaSet")
	labelMap := map[string]string{
		"objectType": icinga.TypeReplicasets,
		"objectName": replicaSet.Name,
	}
	alert, err := mini.CreateAlert(watcher, replicaSet.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Getting ReplicaSetObjectList
	replicaSetObjectList, err := icinga.GetIcingaHostList(watcher, alert)
	if !assert.Nil(t, err) {
		return
	}

	// Delete ReplicaSet
	fmt.Println("--> Deleting ReplicaSet")
	if err := mini.DeleteReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjects(watcher, alert, replicaSetObjectList, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Create ReplicaSet with same name
	fmt.Println("--> ReCreating ReplicaSet with same name")
	newReplicaSet, err := mini.ReCreateReplicaSet(watcher, replicaSet)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete 1st Alert
	fmt.Println("--> Deleting 1st Alert")
	if err := mini.DeleteAlert(watcher, alert); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete ReplicaSet
	fmt.Println("--> Deleting ReplicaSet")
	if err := mini.DeleteReplicaSet(watcher, newReplicaSet); !assert.Nil(t, err) {
		return
	}
}

func TestAlertOnPod(t *testing.T) {
	// Run KubeD
	// runKubeD(setIcingaClient bool)
	// Pass true to set IcingaClient in watcher
	watcher, err := runKubeD(true)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	// Create Pod
	fmt.Println("--> Creating Pod")
	pod, err := mini.CreatePod(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println("--> Creating Alert on Pod")
	labelMap := map[string]string{
		"objectType": icinga.TypePod,
		"objectName": pod.Name,
	}
	alert, err := mini.CreateAlert(watcher, pod.Namespace, labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Getting PodObjectList
	replicaSetObjectList, err := icinga.GetIcingaHostList(watcher, alert)
	if !assert.Nil(t, err) {
		return
	}

	// Delete Pod
	fmt.Println("--> Deleting Pod")
	if err := mini.DeletePod(watcher, pod); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 1st Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjects(watcher, alert, replicaSetObjectList, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Create Pod with same name
	fmt.Println("--> ReCreating Pod with same name")
	newPod, err := mini.ReCreatePod(watcher, pod)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete 1st Alert
	fmt.Println("--> Deleting 1st Alert")
	if err := mini.DeleteAlert(watcher, alert); !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for 2nd Alert.
	fmt.Println("----> Checking Icinga Objects for 1st Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, true, true); !assert.Nil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete Pod
	fmt.Println("--> Deleting Pod")
	if err := mini.DeletePod(watcher, newPod); !assert.Nil(t, err) {
		return
	}
}

func TestInvalidNamespace(t *testing.T) {
	// Run KubeD
	// runKubeD(setIcingaClient bool)
	// Pass true to set IcingaClient in watcher
	watcher, err := runKubeD(true)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	// Create Pod
	fmt.Println("--> Creating Pod")
	pod, err := mini.CreatePod(watcher, "kube-system")
	if !assert.Nil(t, err) {
		return
	}

	fmt.Println("--> Creating Alert on Pod")
	labelMap := map[string]string{
		"objectType": icinga.TypePod,
		"objectName": pod.Name,
	}
	alert, err := mini.CreateAlert(watcher, "default", labelMap, icinga.CheckCommandVolume)
	if !assert.Nil(t, err) {
		return
	}

	// Check Icinga Objects for Alert.
	fmt.Println("----> Checking Icinga Objects for Alert")
	if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); !assert.NotNil(t, err) {
		return
	}
	fmt.Println("---->> Check Successful")

	// Delete Pod
	fmt.Println("--> Deleting Pod")
	if err := mini.DeletePod(watcher, pod); !assert.Nil(t, err) {
		return
	}

	// Delete Alert
	fmt.Println("--> Deleting Alert")
	if err := mini.DeleteAlert(watcher, alert); !assert.Nil(t, err) {
		return
	}
}

func TestAlertOnPodAncestors(t *testing.T) {
	// Run KubeD
	// runKubeD(setIcingaClient bool)
	// Pass true to set IcingaClient in watcher
	watcher, err := runKubeD(true)
	if !assert.Nil(t, err) {
		return
	}
	fmt.Println("--> Running kubeD")

	testAncestor := func(objectType, objectName, namespace string) error {
		fmt.Println("--> Creating Alert on", objectType)
		labelMap := map[string]string{
			"objectType": objectType,
			"objectName": objectName,
		}
		alert, err := mini.CreateAlert(watcher, namespace, labelMap, icinga.CheckCommandVolume)
		if err != nil {
			return err
		}

		// Check Icinga Objects for 1st Alert.
		fmt.Println("----> Checking Icinga Objects for Alert")
		if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, false, false); err != nil {
			return err
		}
		fmt.Println("---->> Check Successful")

		// Delete 1st Alert
		fmt.Println("--> Deleting Alert")
		if err := mini.DeleteAlert(watcher, alert); err != nil {
			return err
		}

		// Check Icinga Objects for Alert.
		fmt.Println("----> Checking Icinga Objects for Alert")
		if err := icinga.CheckIcingaObjectsForAlert(watcher, alert, true, true); err != nil {
			return err
		}
		fmt.Println("---->> Check Successful")
		return nil
	}

	// Create DaemonSet
	fmt.Println("--> Creating DaemonSet")
	daemonSet, err := mini.CreateDaemonSet(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	// Test DaemonSet Ancestor
	fmt.Println("--> Testing DaemonSet Ancestor")
	if err := testAncestor(icinga.TypeDaemonsets, daemonSet.Name, daemonSet.Namespace); !assert.Nil(t, err) {
		return
	}

	// Delete DaemonSet
	fmt.Println("--> Deleting DaemonSet")
	if err := mini.DeleteDaemonSet(watcher, daemonSet); !assert.Nil(t, err) {
		return
	}

	// Create Deployment
	fmt.Println("--> Creating Deployment")
	deployment, err := mini.CreateDeployment(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	// Test Deployment Ancestor
	fmt.Println("--> Testing Deployment Ancestor")
	if err := testAncestor(icinga.TypeDeployments, deployment.Name, deployment.Namespace); !assert.Nil(t, err) {
		return
	}

	// Delete Deployment
	fmt.Println("--> Deleting Deployment")
	if err := mini.DeleteDeployment(watcher, deployment); !assert.Nil(t, err) {
		return
	}

	// Create ReplicaSet
	fmt.Println("--> Creating ReplicaSet")
	replicaSet, err := mini.CreateReplicaSet(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	// Test ReplicaSet Ancestor
	fmt.Println("--> Testing ReplicaSet Ancestor")
	if err := testAncestor(icinga.TypeReplicasets, replicaSet.Name, replicaSet.Namespace); !assert.Nil(t, err) {
		return
	}

	// Delete ReplicaSet
	fmt.Println("--> Deleting ReplicaSet")
	if err := mini.DeleteReplicaSet(watcher, replicaSet); !assert.Nil(t, err) {
		return
	}

	// Create ReplicationController
	fmt.Println("--> Creating ReplicationController")
	rc, err := mini.CreateReplicationController(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	// Test ReplicationController Ancestor
	fmt.Println("--> Testing ReplicationController Ancestor")
	if err := testAncestor(icinga.TypeReplicationcontrollers, rc.Name, rc.Namespace); !assert.Nil(t, err) {
		return
	}

	// Delete ReplicationController
	fmt.Println("--> Deleting ReplicationController")
	if err := mini.DeleteReplicationController(watcher, rc); !assert.Nil(t, err) {
		return
	}

	// Create StatefulSet
	fmt.Println("--> Creating StatefulSet")
	statefulSet, err := mini.CreateStatefulSet(watcher, "default")
	if !assert.Nil(t, err) {
		return
	}

	// Test StatefulSet Ancestor
	fmt.Println("--> Testing StatefulSet Ancestor")
	if err := testAncestor(icinga.TypeStatefulSet, statefulSet.Name, statefulSet.Namespace); !assert.Nil(t, err) {
		return
	}

	// Delete StatefulSet
	fmt.Println("--> Deleting StatefulSet")
	if err := mini.DeleteStatefulSet(watcher, statefulSet); !assert.Nil(t, err) {
		return
	}
}
