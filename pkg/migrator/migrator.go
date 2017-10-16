package migrator

import (
	"errors"
	"fmt"
	"time"

	"github.com/appscode/go/log"
	aci "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/hashicorp/go-version"
	extensions "k8s.io/api/extensions/v1beta1"
	extensionsobj "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
)

type migrationState struct {
	tprRegDeleted bool
	crdCreated    bool
}

type migrator struct {
	kubeClient       clientset.Interface
	apiExtKubeClient apiextensionsclient.Interface
	extClient        cs.MonitoringV1alpha1Interface

	migrationState *migrationState
}

func NewMigrator(kubeClient clientset.Interface, apiExtKubeClient apiextensionsclient.Interface, extClient cs.MonitoringV1alpha1Interface) *migrator {
	return &migrator{
		migrationState:   &migrationState{},
		kubeClient:       kubeClient,
		apiExtKubeClient: apiExtKubeClient,
		extClient:        extClient,
	}
}

func (m *migrator) isMigrationNeeded() (bool, error) {
	v, err := m.kubeClient.Discovery().ServerVersion()
	if err != nil {
		return false, err
	}

	ver, err := version.NewVersion(v.String())
	if err != nil {
		return false, err
	}

	mv := ver.Segments()[1]

	if mv == 7 {
		_, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(
			aci.ResourceNameClusterAlert+"."+api.SchemeGroupVersion.Group,
			metav1.GetOptions{},
		)
		if err != nil {
			if !kerr.IsNotFound(err) {
				return false, err
			}
		} else {
			return true, nil
		}

		_, err = m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(
			aci.ResourceNameNodeAlert+"."+api.SchemeGroupVersion.Group,
			metav1.GetOptions{},
		)
		if err != nil {
			if !kerr.IsNotFound(err) {
				return false, err
			}
		} else {
			return true, nil
		}

		_, err = m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(
			aci.ResourceNamePodAlert+"."+api.SchemeGroupVersion.Group,
			metav1.GetOptions{},
		)
		if err != nil {
			if !kerr.IsNotFound(err) {
				return false, err
			}
		} else {
			return true, nil
		}
	}

	return false, nil
}

func (m *migrator) RunMigration() error {
	needed, err := m.isMigrationNeeded()
	if err != nil {
		return err
	}

	if needed {
		if err := m.migrateTPR2CRD(); err != nil {
			return m.rollback()
		}
	}

	return nil
}

func (m *migrator) migrateTPR2CRD() error {
	log.Debugln("Performing TPR to CRD migration.")

	log.Debugln("Deleting TPRs.")
	if err := m.deleteTPRs(); err != nil {
		return errors.New("Failed to Delete TPRs")
	}

	m.migrationState.tprRegDeleted = true

	log.Debugln("Creating CRDs.")
	if err := m.createCRDs(); err != nil {
		return errors.New("Failed to create CRDs")
	}

	m.migrationState.crdCreated = true

	log.Debugln("Waiting for CRDs to be ready.")
	if err := m.waitForCRDsReady(); err != nil {
		return errors.New("Failed to be ready CRDs")
	}

	return nil
}

func (m *migrator) deleteTPRs() error {
	tprClient := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources()

	deleteTPR := func(resourceName string) error {
		name := resourceName + "." + api.SchemeGroupVersion.Group
		if err := tprClient.Delete(name, &metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to remove %s TPR", name)
		}
		return nil
	}

	if err := deleteTPR(aci.ResourceNameClusterAlert); err != nil {
		return err
	}
	if err := deleteTPR(aci.ResourceNameNodeAlert); err != nil {
		return err
	}
	if err := deleteTPR(aci.ResourceNamePodAlert); err != nil {
		return err
	}
	return nil
}

func (m *migrator) createCRDs() error {
	if err := m.createCRD(aci.ResourceKindClusterAlert, aci.ResourceTypeClusterAlert); err != nil {
		return err
	}
	if err := m.createCRD(aci.ResourceKindNodeAlert, aci.ResourceTypeNodeAlert); err != nil {
		return err
	}
	if err := m.createCRD(aci.ResourceKindPodAlert, aci.ResourceTypePodAlert); err != nil {
		return err
	}
	return nil
}

func (m *migrator) createCRD(resourceKind, resourceType string) error {
	crd := &extensionsobj.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceType + "." + api.SchemeGroupVersion.Group,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Spec: extensionsobj.CustomResourceDefinitionSpec{
			Group:   api.SchemeGroupVersion.Group,
			Version: api.SchemeGroupVersion.Version,
			Scope:   extensionsobj.NamespaceScoped,
			Names: extensionsobj.CustomResourceDefinitionNames{
				Plural: resourceType,
				Kind:   resourceKind,
			},
		},
	}

	crdClient := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()
	_, err := crdClient.Create(crd)
	if err != nil && !kerr.IsAlreadyExists(err) {
		return fmt.Errorf(`Failed to create CRD "%v"`, crd.Spec.Names.Kind)
	}

	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crdEst, err := crdClient.Get(crd.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, cond := range crdEst.Status.Conditions {
			switch cond.Type {
			case extensionsobj.Established:
				if cond.Status == extensionsobj.ConditionTrue {
					return true, err
				}
			case extensionsobj.NamesAccepted:
				if cond.Status == extensionsobj.ConditionFalse {
					fmt.Printf("Name conflict. Reason: %v\n", cond.Reason)
				}
			}
		}
		return false, err
	})

	return nil
}

func (m *migrator) waitForCRDsReady() error {
	labelMap := map[string]string{
		"app": "searchlight",
	}

	return wait.Poll(3*time.Second, 10*time.Minute, func() (bool, error) {
		crdList, err := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions().List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelMap).String(),
		})
		if err != nil {
			return false, err
		}

		if len(crdList.Items) == 3 {
			return true, nil
		}

		return false, nil
	})
}

func (m *migrator) rollback() error {
	log.Debugln("Rolling back migration.")

	ms := m.migrationState

	if ms.crdCreated {
		log.Debugln("Deleting CRDs.")
		err := m.deleteCRDs()
		if err != nil {
			return errors.New("Failed to delete CRDs")
		}
	}

	if ms.tprRegDeleted {
		log.Debugln("Creating TPRs.")
		err := m.createTPRs()
		if err != nil {
			return errors.New("Failed to recreate TPRs")
		}

		err = m.waitForTPRsReady()
		if err != nil {
			return errors.New("Failed to be ready TPRs")
		}
	}

	return nil
}

func (m *migrator) deleteCRDs() error {
	crdClient := m.apiExtKubeClient.ApiextensionsV1beta1().CustomResourceDefinitions()

	deleteCRD := func(resourceType string) error {
		name := resourceType + "." + api.SchemeGroupVersion.Group
		err := crdClient.Delete(name, &metav1.DeleteOptions{})
		if err != nil {
			return fmt.Errorf(`Failed to delete CRD "%s""`, name)
		}
		return nil
	}

	if err := deleteCRD(aci.ResourceTypeClusterAlert); err != nil {
		return err
	}
	if err := deleteCRD(aci.ResourceTypeNodeAlert); err != nil {
		return err
	}
	if err := deleteCRD(aci.ResourceTypePodAlert); err != nil {
		return err
	}
	return nil
}

func (m *migrator) createTPRs() error {
	if err := m.createTPR(aci.ResourceNameClusterAlert); err != nil {
		return err
	}
	if err := m.createTPR(aci.ResourceNameNodeAlert); err != nil {
		return err
	}
	if err := m.createTPR(aci.ResourceNamePodAlert); err != nil {
		return err
	}
	return nil
}

func (m *migrator) createTPR(resourceName string) error {
	name := resourceName + "." + api.SchemeGroupVersion.Group
	_, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Get(name, metav1.GetOptions{})
	if !kerr.IsNotFound(err) {
		return err
	}

	thirdPartyResource := &extensions.ThirdPartyResource{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "ThirdPartyResource",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"app": "searchlight",
			},
		},
		Description: "Searchlight by AppsCode - Alerts for Kubernetes",
		Versions: []extensions.APIVersion{
			{
				Name: api.SchemeGroupVersion.Version,
			},
		},
	}

	_, err = m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().Create(thirdPartyResource)
	return err
}

func (m *migrator) waitForTPRsReady() error {
	labelMap := map[string]string{
		"app": "searchlight",
	}

	return wait.Poll(3*time.Second, 10*time.Minute, func() (bool, error) {
		crdList, err := m.kubeClient.ExtensionsV1beta1().ThirdPartyResources().List(metav1.ListOptions{
			LabelSelector: labels.SelectorFromSet(labelMap).String(),
		})
		if err != nil {
			return false, err
		}

		if len(crdList.Items) == 3 {
			return true, nil
		}

		return false, nil
	})
}
