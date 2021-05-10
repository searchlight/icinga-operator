/*
Copyright AppsCode Inc. and Contributors

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

package acknowledgement

import (
	"context"
	"encoding/json"
	"fmt"

	"go.searchlight.dev/icinga-operator/apis/incidents"
	"go.searchlight.dev/icinga-operator/apis/incidents/v1alpha1"
	monitoring "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/client/clientset/versioned"
	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	restconfig "k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

type REST struct {
	client versioned.Interface
	ic     *icinga.Client // TODO: init
}

var _ rest.Creater = &REST{}
var _ rest.Scoper = &REST{}
var _ rest.GracefulDeleter = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}
var _ rest.CategoriesProvider = &REST{}

func NewREST(config *restconfig.Config, ic *icinga.Client) *REST {
	return &REST{
		client: versioned.NewForConfigOrDie(config),
		ic:     ic,
	}
}

func (r *REST) NamespaceScoped() bool {
	return true
}

func (r *REST) New() runtime.Object {
	return &incidents.Acknowledgement{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ResourceKindAcknowledgement)
}

func (r *REST) Categories() []string {
	return []string{"monitoring", "appscode", "all"}
}

func (r *REST) Create(ctx context.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ *metav1.CreateOptions) (runtime.Object, error) {
	req := obj.(*incidents.Acknowledgement)

	if errs := validate(req); len(errs) > 0 {
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: incidents.GroupName, Kind: v1alpha1.ResourceKindAcknowledgement}, req.Name, errs)
	}

	host, service, err := r.getIcingaObjects(req.Namespace, req.Name)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, service, host)
	mp["comment"] = req.Request.Comment
	mp["notify"] = !req.Request.SkipNotify
	if user, ok := apirequest.UserFrom(ctx); ok {
		mp["author"] = user.GetName()
	}

	ack, err := json.Marshal(mp)
	if err != nil {
		return nil, err
	}
	response := r.ic.Actions("acknowledge-problem").Update([]string{}, string(ack)).Do()
	if response.Err != nil {
		return nil, response.Err
	}
	var icingaresp icinga.APIResponse
	status, err := response.Into(&icingaresp)
	if err != nil {
		return nil, err
	}
	if status != 200 {
		return nil, errors.New(string(icingaresp.ResponseBody))
	}
	req.Response = incidents.AcknowledgementResponse{
		Timestamp: metav1.Now(),
	}
	return req, nil
}

func validate(o *incidents.Acknowledgement) field.ErrorList {
	klog.Infof("Validating fields for Acknowledgement %s\n", o.Name)
	errs := field.ErrorList{}

	if o.Request.Comment == "" {
		errs = append(errs,
			field.Invalid(field.NewPath("request", "comment"), o.Request.Comment, "comment must not be empty"))
	}

	// perform validation here and add to errlist using field.Invalid
	return errs
}

func (r *REST) Delete(ctx context.Context, name string, deleteValidation rest.ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	namespace, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, false, apierrors.NewBadRequest("namespace missing")
	}

	host, service, err := r.getIcingaObjects(namespace, name)
	if err != nil {
		return nil, false, err
	}

	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, service, host)
	ack, err := json.Marshal(mp)
	if err != nil {
		return nil, false, err
	}

	response := r.ic.Actions("remove-acknowledgement").Update([]string{}, string(ack)).Do()
	if response.Err != nil {
		return nil, false, response.Err
	}
	var icingaResp icinga.APIResponse
	status, err := response.Into(&icingaResp)
	if err != nil {
		return nil, false, err
	}
	if status != 200 {
		return nil, false, errors.New(string(icingaResp.ResponseBody))
	}

	resp := &incidents.Acknowledgement{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Response: incidents.AcknowledgementResponse{
			Timestamp: metav1.Now(),
		},
	}

	return resp, true, nil
}

func (r *REST) getIcingaObjects(namespace, name string) (host string, service string, err error) {
	incident, err := r.client.MonitoringV1alpha1().Incidents(namespace).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return "", "", errors.Errorf("incident %s/%s not found", namespace, name)
		}
		return "", "", errors.Wrapf(err, "failed to determine incident %s/%s", namespace, name)
	}

	icingaHost := &icinga.IcingaHost{AlertNamespace: namespace}

	service, ok := incident.Labels[monitoring.LabelKeyAlert]
	if !ok {
		return "", "", errors.Errorf("incident %s/%s is missing label %s", namespace, name, monitoring.LabelKeyAlert)
	}
	icingaHost.Type, ok = incident.Labels[monitoring.LabelKeyAlertType]
	if !ok {
		return "", "", errors.Errorf("incident %s/%s is missing label %s", namespace, name, monitoring.LabelKeyAlertType)
	} else if !icinga.IsValidHostType(icingaHost.Type) {
		return "", "", errors.Errorf("incident %s/%s has invalid value %s for label %s", namespace, name, icingaHost.Type, monitoring.LabelKeyAlertType)
	}
	if icingaHost.Type != icinga.TypeCluster {
		icingaHost.ObjectName, ok = incident.Labels[monitoring.LabelKeyObjectName]
		if !ok {
			return "", "", errors.Errorf("incident %s/%s is missing label %s", namespace, name, monitoring.LabelKeyObjectName)
		}
	}
	host, err = icingaHost.Name()
	return
}
