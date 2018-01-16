package notifier

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/appscode/envconfig"
	"github.com/appscode/go-notify"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/log"
	logs "github.com/appscode/go/log/golog"
	api "github.com/appscode/searchlight/apis/monitoring/v1alpha1"
	cs "github.com/appscode/searchlight/client/typed/monitoring/v1alpha1"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	HostName  string
	AlertName string
	Type      string
	State     string
	Output    string
	// The time object is used in icinga to send request. This
	// indicates detection time from icinga.
	Time    int64
	Author  string
	Comment string
}

type Secret struct {
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}

func getLoader(client kubernetes.Interface, alert api.Alert) (envconfig.LoaderFunc, error) {
	cfg, err := client.CoreV1().Secrets(alert.GetNamespace()).Get(alert.GetNotifierSecretName(), metav1.GetOptions{})
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

func getAlert(kh *icinga.IcingaHost, extClient cs.MonitoringV1alpha1Interface, alertName string) (api.Alert, error) {
	switch kh.Type {
	case icinga.TypePod:
		return extClient.PodAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	case icinga.TypeNode:
		return extClient.NodeAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	case icinga.TypeCluster:
		return extClient.ClusterAlerts(kh.AlertNamespace).Get(alertName, metav1.GetOptions{})
	}
	return nil, fmt.Errorf("unknown host type %s", kh.Type)
}

func sendNotification(req *Request) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		log.Fatalln(err)
	}

	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		log.Fatalln(err)
	}

	alert, err := getAlert(host, cs.NewForConfigOrDie(config), req.AlertName)
	if err != nil {
		log.Fatalln(err)
	}

	loader, err := getLoader(kubernetes.NewForConfigOrDie(config), alert)
	if err != nil {
		log.Fatalln(err)
	}

	receivers := alert.GetReceivers()

	for _, receiver := range receivers {
		if receiver.State != req.State || len(receiver.To) == 0 {
			continue
		}
		notifyVia, err := unified.LoadVia(receiver.Notifier, loader)
		if err != nil {
			log.Errorln(err)
			continue
		}

		switch n := notifyVia.(type) {
		case notify.ByEmail:
			var body string
			body, err = RenderMail(alert, req)
			if err != nil {
				log.Errorf("Failed to render email. Reason: %s", err)
				break
			}
			err = n.To(receiver.To[0], receiver.To[1:]...).
				WithSubject(RenderSubject(alert, req)).
				WithBody(body).
				WithNoTracking().
				SendHtml()
		case notify.BySMS:
			err = n.To(receiver.To[0], receiver.To[1:]...).
				WithBody(RenderSMS(alert, req)).
				Send()
		case notify.ByChat:
			err = n.To(receiver.To[0], receiver.To[1:]...).
				WithBody(RenderSMS(alert, req)).
				Send()
		case notify.ByPush:
			err = n.To(receiver.To[0:]...).
				WithBody(RenderSMS(alert, req)).
				Send()
		}

		if err != nil {
			log.Errorln(err)
		} else {
			log.Infof("Notification sent using %s", receiver.Notifier)
		}
	}
}

func NewCmd() *cobra.Command {
	var req Request
	var eventTime string

	c := &cobra.Command{
		Use:   "notifier",
		Short: "AppsCode Icinga2 Notifier",
		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "alert", "host", "type", "state", "time")
			t, err := time.Parse("2006-01-02 15:04:05 +0000", eventTime)
			if err != nil {
				log.Errorln(err)
				os.Exit(1)

			}
			req.Time = t.Unix()
			sendNotification(&req)
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")

	c.Flags().StringVarP(&req.HostName, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.AlertName, "alert", "A", "", "Kubernetes alert object name")
	c.Flags().StringVar(&req.Type, "type", "", "Notification type (PROBLEM | ACKNOWLEDGEMENT | RECOVERY)")
	c.Flags().StringVar(&req.State, "state", "", "Service state (OK | WARNING | CRITICAL)")
	c.Flags().StringVar(&req.Output, "output", "", "Service output")
	c.Flags().StringVar(&eventTime, "time", "", "Event time")
	c.Flags().StringVarP(&req.Author, "author", "a", "", "Event author name")
	c.Flags().StringVarP(&req.Comment, "comment", "c", "", "Event comment")

	c.Flags().AddGoFlagSet(flag.CommandLine)
	logs.InitLogs()

	return c
}
