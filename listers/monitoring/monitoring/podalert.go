/*
Copyright 2017 The Searchlight Authors.

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

// This file was automatically generated by lister-gen

package monitoring

import (
	monitoring "github.com/appscode/searchlight/apis/monitoring"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// PodAlertLister helps list PodAlerts.
type PodAlertLister interface {
	// List lists all PodAlerts in the indexer.
	List(selector labels.Selector) (ret []*monitoring.PodAlert, err error)
	// PodAlerts returns an object that can list and get PodAlerts.
	PodAlerts(namespace string) PodAlertNamespaceLister
	PodAlertListerExpansion
}

// podAlertLister implements the PodAlertLister interface.
type podAlertLister struct {
	indexer cache.Indexer
}

// NewPodAlertLister returns a new PodAlertLister.
func NewPodAlertLister(indexer cache.Indexer) PodAlertLister {
	return &podAlertLister{indexer: indexer}
}

// List lists all PodAlerts in the indexer.
func (s *podAlertLister) List(selector labels.Selector) (ret []*monitoring.PodAlert, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*monitoring.PodAlert))
	})
	return ret, err
}

// PodAlerts returns an object that can list and get PodAlerts.
func (s *podAlertLister) PodAlerts(namespace string) PodAlertNamespaceLister {
	return podAlertNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// PodAlertNamespaceLister helps list and get PodAlerts.
type PodAlertNamespaceLister interface {
	// List lists all PodAlerts in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*monitoring.PodAlert, err error)
	// Get retrieves the PodAlert from the indexer for a given namespace and name.
	Get(name string) (*monitoring.PodAlert, error)
	PodAlertNamespaceListerExpansion
}

// podAlertNamespaceLister implements the PodAlertNamespaceLister
// interface.
type podAlertNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all PodAlerts in the indexer for a given namespace.
func (s podAlertNamespaceLister) List(selector labels.Selector) (ret []*monitoring.PodAlert, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*monitoring.PodAlert))
	})
	return ret, err
}

// Get retrieves the PodAlert from the indexer for a given namespace and name.
func (s podAlertNamespaceLister) Get(name string) (*monitoring.PodAlert, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(monitoring.Resource("podalert"), name)
	}
	return obj.(*monitoring.PodAlert), nil
}
