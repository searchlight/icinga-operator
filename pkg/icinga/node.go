package icinga

import (
	"bytes"
	"text/template"

	"github.com/appscode/go/errors"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

type NodeHost struct {
	commonHost

	KubeClient kubernetes.Interface
	ExtClient  cs.MonitoringV1alpha1Interface
	//*types.Context
}

func NewNodeHost(kubeClient kubernetes.Interface, extClient cs.MonitoringV1alpha1Interface, IcingaClient *Client) *NodeHost {
	return &NodeHost{
		KubeClient: kubeClient,
		ExtClient:  extClient,
		commonHost: commonHost{
			IcingaClient: IcingaClient,
		},
	}
}

func (h *NodeHost) getHost(alert api.NodeAlert, node apiv1.Node) IcingaHost {
	nodeIP := "127.0.0.1"
	for _, ip := range node.Status.Addresses {
		if ip.Type == internalIP {
			nodeIP = ip.Address
			break
		}
	}
	return IcingaHost{
		ObjectName:     node.Name,
		Type:           TypeNode,
		AlertNamespace: alert.Namespace,
		IP:             nodeIP,
	}
}

func (h *NodeHost) expandVars(alertSpec api.NodeAlertSpec, kh IcingaHost, attrs map[string]interface{}) error {
	commandVars := api.NodeCommands[alertSpec.Check].Vars
	for key, val := range alertSpec.Vars {
		if v, found := commandVars[key]; found {
			if v.Parameterized {
				type Data struct {
					NodeName string
					NodeIP   string
				}
				tmpl, err := template.New("").Parse(val)
				if err != nil {
					return err
				}
				var buf bytes.Buffer
				err = tmpl.Execute(&buf, Data{NodeName: kh.ObjectName, NodeIP: kh.IP})
				if err != nil {
					return err
				}
				attrs[IVar(key)] = buf.String()
			} else {
				attrs[IVar(key)] = val
			}
		} else {
			return errors.Newf("variable %v not found", key).Err()
		}
	}
	return nil
}

// set Alert in Icinga LocalHost
func (h *NodeHost) Create(alert api.NodeAlert, node apiv1.Node) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert, node)

	if has, err := h.CheckIcingaService(alert.Name, kh); err != nil || has {
		return err
	}

	if err := h.CreateIcingaHost(kh); err != nil {
		return errors.FromErr(err).Err()
	}

	attrs := make(map[string]interface{})
	attrs["check_command"] = alertSpec.Check
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	if err := h.expandVars(alertSpec, kh, attrs); err != nil {
		return err
	}
	if err := h.CreateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}

	return h.CreateIcingaNotification(alert, kh)
}

func (h *NodeHost) Update(alert api.NodeAlert, node apiv1.Node) error {
	alertSpec := alert.Spec
	kh := h.getHost(alert, node)

	attrs := make(map[string]interface{})
	if alertSpec.CheckInterval.Seconds() > 0 {
		attrs["check_interval"] = alertSpec.CheckInterval.Seconds()
	}
	if err := h.expandVars(alertSpec, kh, attrs); err != nil {
		return err
	}
	if err := h.UpdateIcingaService(alert.Name, kh, attrs); err != nil {
		return errors.FromErr(err).Err()
	}

	return h.UpdateIcingaNotification(alert, kh)
}

func (h *NodeHost) Delete(alert api.NodeAlert, node apiv1.Node) error {
	kh := h.getHost(alert, node)

	if err := h.DeleteIcingaService(alert.Name, kh); err != nil {
		return errors.FromErr(err).Err()
	}
	return h.DeleteIcingaHost(kh)
}
