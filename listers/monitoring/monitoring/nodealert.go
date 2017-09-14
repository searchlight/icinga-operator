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

// NodeAlertLister helps list NodeAlerts.
type NodeAlertLister interface {
	// List lists all NodeAlerts in the indexer.
	List(selector labels.Selector) (ret []*monitoring.NodeAlert, err error)
	// NodeAlerts returns an object that can list and get NodeAlerts.
	NodeAlerts(namespace string) NodeAlertNamespaceLister
	NodeAlertListerExpansion
}

// nodeAlertLister implements the NodeAlertLister interface.
type nodeAlertLister struct {
	indexer cache.Indexer
}

// NewNodeAlertLister returns a new NodeAlertLister.
func NewNodeAlertLister(indexer cache.Indexer) NodeAlertLister {
	return &nodeAlertLister{indexer: indexer}
}

// List lists all NodeAlerts in the indexer.
func (s *nodeAlertLister) List(selector labels.Selector) (ret []*monitoring.NodeAlert, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*monitoring.NodeAlert))
	})
	return ret, err
}

// NodeAlerts returns an object that can list and get NodeAlerts.
func (s *nodeAlertLister) NodeAlerts(namespace string) NodeAlertNamespaceLister {
	return nodeAlertNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// NodeAlertNamespaceLister helps list and get NodeAlerts.
type NodeAlertNamespaceLister interface {
	// List lists all NodeAlerts in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*monitoring.NodeAlert, err error)
	// Get retrieves the NodeAlert from the indexer for a given namespace and name.
	Get(name string) (*monitoring.NodeAlert, error)
	NodeAlertNamespaceListerExpansion
}

// nodeAlertNamespaceLister implements the NodeAlertNamespaceLister
// interface.
type nodeAlertNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all NodeAlerts in the indexer for a given namespace.
func (s nodeAlertNamespaceLister) List(selector labels.Selector) (ret []*monitoring.NodeAlert, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*monitoring.NodeAlert))
	})
	return ret, err
}

// Get retrieves the NodeAlert from the indexer for a given namespace and name.
func (s nodeAlertNamespaceLister) Get(name string) (*monitoring.NodeAlert, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(monitoring.Resource("nodealert"), name)
	}
	return obj.(*monitoring.NodeAlert), nil
}
