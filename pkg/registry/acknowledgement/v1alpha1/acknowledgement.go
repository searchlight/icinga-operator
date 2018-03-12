package v1alpha1

import (
	"encoding/json"
	"fmt"

	incidents "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	monitoring "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/client/clientset/versioned"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	restconfig "k8s.io/client-go/rest"
)

type REST struct {
	client versioned.Interface
	ic     *icinga.Client // TODO: init
}

var _ rest.Creater = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}

func NewREST(config *restconfig.Config, ic *icinga.Client) *REST {
	return &REST{
		client: versioned.NewForConfigOrDie(config),
		ic:     ic,
	}
}

func (r *REST) New() runtime.Object {
	return &incidents.Acknowledgement{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return incidents.SchemeGroupVersion.WithKind("Acknowledgement")
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*incidents.Acknowledgement)

	in, err := r.client.MonitoringV1alpha1().Incidents(req.Namespace).Get(req.Name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil, errors.Errorf("incident %s/%s not found", req.Namespace, req.Name)
		}
		return nil, errors.Wrapf(err, "failed to determine incident %s/%s", req.Namespace, req.Name)
	}

	alertName, ok := in.Annotations[monitoring.LabelKeyAlert]
	if !ok {
		return nil, errors.Errorf("incident %s/%s is missing annotation %s", req.Namespace, req.Name, monitoring.LabelKeyAlert)
	}
	alertType, ok := in.Annotations[monitoring.LabelKeyAlertType]
	if !ok {
		return nil, errors.Errorf("incident %s/%s is missing annotation %s", req.Namespace, req.Name, monitoring.LabelKeyAlertType)
	}
	objName, ok := in.Annotations[monitoring.LabelKeyObjectName]
	if !ok {
		return nil, errors.Errorf("incident %s/%s is missing annotation %s", req.Namespace, req.Name, monitoring.LabelKeyObjectName)
	}

	host := &icinga.IcingaHost{
		ObjectName:     objName,
		AlertNamespace: req.Namespace,
	}

	switch alertType {
	case api.ResourceTypePodAlert:
		host.Type = icinga.TypePod
	case api.ResourceTypeNodeAlert:
		host.Type = icinga.TypeNode
	case api.ResourceTypeClusterAlert:
		host.Type = icinga.TypeCluster
	}

	hostName, err := host.Name()
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, alertName, hostName)
	mp["comment"] = req.Request
	mp["notify"] = !req.Request.SkipNotify
	if user, ok := apirequest.UserFrom(ctx); ok {
		mp["author"] = user.GetName()
	}

	ack, err := json.Marshal(mp)
	if err != nil {
		return nil, err
	}
	response := r.ic.Actions("acknowledge-problem").Update([]string{}, string(ack)).Do()
	if response.Status != 200 {
		return nil, response.Err
	}
	req.Response = &incidents.AcknowledgementResponse{
		Timestamp: metav1.Now(),
	}
	return req, nil
}
