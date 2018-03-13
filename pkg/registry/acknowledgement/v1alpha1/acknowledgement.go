package v1alpha1

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	monitoring "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	"github.com/appscode/searchlight/client/clientset/versioned"
	"github.com/appscode/searchlight/pkg/icinga"
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
	return &api.Acknowledgement{}
}

func (r *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return api.SchemeGroupVersion.WithKind(api.ResourceKindAcknowledgement)
}

func (r *REST) Create(ctx apirequest.Context, obj runtime.Object, _ rest.ValidateObjectFunc, _ bool) (runtime.Object, error) {
	req := obj.(*api.Acknowledgement)

	if errs := validate(req); len(errs) > 0 {
		return nil, apierrors.NewInvalid(schema.GroupKind{Group: api.GroupName, Kind: api.ResourceKindAcknowledgement}, req.Name, errs)
	}

	incident, err := r.client.MonitoringV1alpha1().Incidents(req.Namespace).Get(req.Name, metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) {
			return nil, errors.Errorf("incident %s/%s not found", req.Namespace, req.Name)
		}
		return nil, errors.Wrapf(err, "failed to determine incident %s/%s", req.Namespace, req.Name)
	}

	host := &icinga.IcingaHost{AlertNamespace: req.Namespace}

	alertName, ok := incident.Labels[monitoring.LabelKeyAlert]
	if !ok {
		return nil, errors.Errorf("incident %s/%s is missing label %s", req.Namespace, req.Name, monitoring.LabelKeyAlert)
	}
	host.Type, ok = incident.Labels[monitoring.LabelKeyAlertType]
	if !ok {
		return nil, errors.Errorf("incident %s/%s is missing label %s", req.Namespace, req.Name, monitoring.LabelKeyAlertType)
	} else if !icinga.IsValidHostType(host.Type) {
		return nil, errors.Errorf("incident %s/%s has invalid value %s for label %s", req.Namespace, req.Name, host.Type, monitoring.LabelKeyAlertType)
	}
	if host.Type != icinga.TypeCluster {
		host.ObjectName, ok = incident.Labels[monitoring.LabelKeyObjectName]
		if !ok {
			return nil, errors.Errorf("incident %s/%s is missing label %s", req.Namespace, req.Name, monitoring.LabelKeyObjectName)
		}
	}
	hostName, err := host.Name()
	if err != nil {
		return nil, err
	}

	mp := make(map[string]interface{})
	mp["type"] = "Service"
	mp["filter"] = fmt.Sprintf(`service.name == "%s" && host.name == "%s"`, alertName, hostName)
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
	req.Response = api.AcknowledgementResponse{
		Timestamp: metav1.Now(),
	}
	return req, nil
}

func validate(o *api.Acknowledgement) field.ErrorList {
	log.Infof("Validating fields for Acknowledgement %s\n", o.Name)
	errs := field.ErrorList{}

	if o.Request.Comment == "" {
		errs = append(errs,
			field.Invalid(field.NewPath("request", "comment"), o.Request.Comment, "comment must not be empty"))
	}

	// perform validation here and add to errlist using field.Invalid
	return errs
}
