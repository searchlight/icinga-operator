package notifier

import (
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/appscode/envconfig"
	notify "github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/log"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/clientset/versioned/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/plugins"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"kmodules.xyz/client-go/logs"
	"kmodules.xyz/client-go/tools/clientcmd"
)

type notifier struct {
	client    corev1.SecretInterface
	extClient cs.MonitoringV1alpha1Interface
	options   options
}

func newPlugin(client corev1.SecretInterface, extClient cs.MonitoringV1alpha1Interface, opts options) *notifier {
	return &notifier{client, extClient, opts}
}

func newPluginFromConfig(opts options) (*notifier, error) {
	config, err := clientcmd.BuildConfigFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	extClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return newPlugin(client.CoreV1().Secrets(opts.host.AlertNamespace), extClient, opts), nil
}

type options struct {
	kubeconfigPath string
	contextName    string

	alertName        string
	notificationType string
	serviceState     string
	serviceOutput    string
	// The time object is used in icinga to send request. This
	// indicates detection time from icinga.
	time    time.Time
	author  string
	comment string
	// IcingaHost
	hostname string
	host     *icinga.IcingaHost
}

const (
	stateOK       = "OK"
	stateWarning  = "Warning"
	stateCritical = "Critical"
	stateUnknown  = "Unknown"
)

func (o *options) complete(cmd *cobra.Command) (err error) {
	o.host, err = icinga.ParseHost(o.hostname)
	if err != nil {
		return errors.New("invalid icinga host.name")
	}

	eventTime, err := cmd.Flags().GetString(flagEventTime)
	if err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02 15:04:05 +0000", eventTime)
	if err != nil {
		return err

	}
	o.time = t

	// sanitized state to preferred form
	switch strings.ToUpper(o.serviceState) {
	case "OK":
		o.serviceState = stateOK
	case "WARNING":
		o.serviceState = stateWarning
	case "CRITICAL":
		o.serviceState = stateCritical
	default:
		o.serviceState = stateUnknown
	}

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return
	}
	return nil
}

func (o *options) validate() error {
	return nil
}

type Secret struct {
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}

func (n *notifier) getLoader(alert api.Alert) (envconfig.LoaderFunc, error) {
	cfg, err := n.client.Get(alert.GetNotifierSecretName(), metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return func(key string) (value string, found bool) {
		var bytes []byte
		bytes, found = cfg.Data[key]
		value = string(bytes)
		return
	}, nil
}

func (n *notifier) getAlert() (api.Alert, error) {
	opts := n.options
	switch opts.host.Type {
	case icinga.TypePod:
		return n.extClient.PodAlerts(opts.host.AlertNamespace).Get(opts.alertName, metav1.GetOptions{})
	case icinga.TypeNode:
		return n.extClient.NodeAlerts(opts.host.AlertNamespace).Get(opts.alertName, metav1.GetOptions{})
	case icinga.TypeCluster:
		return n.extClient.ClusterAlerts(opts.host.AlertNamespace).Get(opts.alertName, metav1.GetOptions{})
	}
	return nil, fmt.Errorf("unknown host type %s", opts.host.Type)
}

func (n *notifier) sendToReceiver(alert api.Alert, receiver api.Receiver, loader envconfig.LoaderFunc) error {
	notifyVia, err := unified.LoadVia(receiver.Notifier, loader)
	if err != nil {
		return err
	}

	switch nv := notifyVia.(type) {
	case notify.ByEmail:
		var body string
		body, err = n.RenderMail(alert)
		if err != nil {
			return fmt.Errorf("failed to render email. Reason: %s", err)
		}
		return nv.To(receiver.To[0], receiver.To[1:]...).
			WithSubject(n.RenderSubject(receiver)).
			WithBody(body).
			WithNoTracking().
			SendHtml()
	case notify.BySMS:
		return nv.To(receiver.To[0], receiver.To[1:]...).
			WithBody(n.RenderSMS(receiver)).
			Send()
	case notify.ByChat:
		return nv.To(receiver.To[0], receiver.To[1:]...).
			WithBody(n.RenderSMS(receiver)).
			Send()
	case notify.ByPush:
		return nv.To(receiver.To[0:]...).
			WithBody(n.RenderSMS(receiver)).
			Send()
	default:
		return fmt.Errorf(`invalid notifier "%s"`, receiver.Notifier)
	}

}

func (n *notifier) sendNotification() {

	alert, err := n.getAlert()
	if err != nil {
		log.Fatalln(err)
	}

	loader, err := n.getLoader(alert)
	if err != nil {
		log.Fatalln(err)
	}

	serviceState := n.options.serviceState
	if api.AlertType(n.options.notificationType) == api.NotificationRecovery {
		if incident, _ := n.getIncident(); incident != nil {
			if lastNonOKState := n.getLastNonOKState(incident); lastNonOKState != "" {
				serviceState = lastNonOKState
			}
		}
	}

	receivers := alert.GetReceivers()

	for _, receiver := range receivers {
		if len(receiver.To) == 0 || !strings.EqualFold(receiver.State, serviceState) {
			continue
		}

		if err = n.sendToReceiver(alert, receiver, loader); err != nil {
			log.Errorln(err)
		} else {
			log.Infof("Notification sent using %s", receiver.Notifier)
		}
	}

	if err := n.reconcileIncident(); err != nil {
		log.Errorln(err)
	}
}

const (
	flagEventTime = "time"
	flagAlert     = "alert"
	flagType      = "type"
	flagState     = "state"
)

func NewCmd() *cobra.Command {
	var opts options

	c := &cobra.Command{
		Use:   "notifier",
		Short: "AppsCode Icinga2 Notifier",
		Run: func(cmd *cobra.Command, args []string) {

			flags.EnsureRequiredFlags(cmd, flagAlert, plugins.FlagHost, flagType, flagState, flagEventTime)

			if err := opts.complete(cmd); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			if err := opts.validate(); err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin, err := newPluginFromConfig(opts)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			plugin.sendNotification()
		},
	}

	c.Flags().StringVarP(&opts.hostname, plugins.FlagHost, "H", "", "Icinga host name")
	c.Flags().StringVarP(&opts.alertName, flagAlert, "A", "", "Kubernetes alert object name")
	c.Flags().StringVar(&opts.notificationType, flagType, "", "Notification type (PROBLEM | ACKNOWLEDGEMENT | RECOVERY)")
	c.Flags().StringVar(&opts.serviceState, flagState, "", "Service state (OK | Warning | Critical)")
	c.Flags().StringVar(&opts.serviceOutput, "output", "", "Service output")
	c.Flags().String(flagEventTime, "", "Event time")
	c.Flags().StringVarP(&opts.author, "author", "a", "", "Event author name")
	c.Flags().StringVarP(&opts.comment, "comment", "c", "", "Event comment")

	c.Flags().AddGoFlagSet(flag.CommandLine)
	logs.InitLogs()

	return c
}
