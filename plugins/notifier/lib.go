package notifier

import (
	"flag"
	"fmt"
	"os"
	"time"

	api "github.com/appscode/api/kubernetes/v1beta1"
	apiv1 "k8s.io/client-go/pkg/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/io"
	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/plugins/notifier/driver/extpoints"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/hipchat"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/mailgun"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/plivo"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/slack"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/smtp"
	_ "github.com/appscode/searchlight/plugins/notifier/driver/twilio"
	"github.com/spf13/cobra"
	"github.com/appscode/searchlight/pkg/icinga"
	"io/ioutil"
	"strings"
	"github.com/appscode/envconfig"
	"github.com/appscode/searchlight/pkg/util"
	clientset "k8s.io/client-go/kubernetes"
	"github.com/appscode/go-notify/unified"
	"github.com/appscode/go-notify"
)

const (
	appscodeConfigPath = "/var/run/config/appscode/"
	appscodeSecretPath = "/var/run/secrets/appscode/"

	notifyVia = "NOTIFY_VIA"
)

type Request struct {
	AlertPhid string `protobuf:"bytes,1,opt,name=alert_phid,json=alertPhid" json:"alert_phid,omitempty"`
	HostName  string `protobuf:"bytes,2,opt,name=host_name,json=hostName" json:"host_name,omitempty"`
	Type      string `protobuf:"bytes,3,opt,name=type" json:"type,omitempty"`
	State     string `protobuf:"bytes,4,opt,name=state" json:"state,omitempty"`
	Output    string `protobuf:"bytes,5,opt,name=output" json:"output,omitempty"`
	// The time object is used in icinga to send request. This
	// indicates detection time from icinga.
	Time                int64  `protobuf:"varint,6,opt,name=time" json:"time,omitempty"`
	Author              string `protobuf:"bytes,7,opt,name=author" json:"author,omitempty"`
	Comment             string `protobuf:"bytes,8,opt,name=comment" json:"comment,omitempty"`
	KubernetesAlertName string `protobuf:"bytes,9,opt,name=kubernetes_alert_name,json=kubernetesAlertName" json:"kubernetes_alert_name,omitempty"`
	KubernetesCluster   string `protobuf:"bytes,10,opt,name=kubernetes_cluster,json=kubernetesCluster" json:"kubernetes_cluster,omitempty"`
}

type Secret struct {
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}

func namespace() string {
	if ns := os.Getenv("OPERATOR_NAMESPACE"); ns != "" {
		return ns
	}
	if data, err := ioutil.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		if ns := strings.TrimSpace(string(data)); len(ns) > 0 {
			return ns
		}
	}
	return apiv1.NamespaceDefault
}


func getLoader(client clientset.Interface) (envconfig.LoaderFunc, error) {

	secretName := os.Getenv(icinga.ICINGA_NOTIFIER_SECRET_NAME)
	secretNamespace := namespace()

	cfg, err := client.CoreV1().
		Secrets(secretNamespace).
		Get(secretName, metav1.GetOptions{})
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

func sendNotification(req *api.IncidentNotifyRequest) {

	client, err := util.NewClient()
	if err != nil {

	}

	loader, err := getLoader(client.Client)
	if err != nil {
		return
	}
	notifier, err := unified.Load(loader)
	if err != nil {
		return
	}

	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		return err
	}

	client.ExtClient.PodAlerts(namespace).Get(alertName)

	alert, err := driver.GetAlertInfo(host.AlertNamespace, req.KubernetesAlertName)
	if err != nil {
		return err
	}
	

	switch n := notifier.(type) {
	case notify.ByEmail:
		receivers := getArray(loader, "CLUSTER_ADMIN_EMAIL")
		if len(receivers) == 0 {
			return n.UID(), errors.New("Missing / invalid cluster admin email(s)")
		}
		n = n.To(receivers[0], receivers[1:]...)
		return n.UID(), n.WithSubject("Cluster CA Certificate").WithBody(msg).Send()
	case notify.BySMS:
		receivers := getArray(loader, "CLUSTER_ADMIN_PHONE")
		if len(receivers) == 0 {
			return n.UID(), errors.New("Missing / invalid cluster admin phone number(s)")
		}
		n = n.To(receivers[0], receivers[1:]...)
		return n.UID(), n.WithBody(msg).Send()
	case notify.ByChat:
		return n.UID(), n.WithBody(msg).Send()
	}
	return "", errors.New("Unknown notifier")


	notifyVia := os.Getenv(notifyVia)
	if notifyVia == "" {
		log.Errorln("No fallback notifier set")
		os.Exit(1)
	}

	cluster_uid, err := io.ReadFile(appscodeConfigPath + "cluster-name")
	if err != nil {
		cluster_uid = ""
	}

	req.KubernetesCluster = cluster_uid
	driver := extpoints.Drivers.Lookup(notifyVia)
	if driver == nil {
		log.Errorln("Invalid failback notifier")
		os.Exit(1)
	}

	if err := driver.Notify(req); err != nil {
		log.Errorln(err)
	} else {
		log.Debug(fmt.Sprintf("Notification sent via %s", notifyVia))
	}
}

func NewCmd() *cobra.Command {
	var req Request
	var eventTime string

	c := &cobra.Command{
		Use:     "notifier",
		Short:   "AppsCode Icinga2 Notifier",
		Example: "",
		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "alert", "host", "type", "state", "output", "time")
			t, err := time.Parse("2006-01-02 15:04:05 +0000", eventTime)
			if err != nil {
				log.Errorln(err)
				os.Exit(1)

			}
			req.Time = t.Unix()
			sendNotification(&req)
		},
	}

	c.Flags().StringVarP(&req.KubernetesAlertName, "alert", "A", "", "Kubernetes alert object name")
	c.Flags().StringVarP(&req.HostName, "host", "H", "", "Icinga host name")
	c.Flags().StringVar(&req.Type, "type", "", "Notification type")
	c.Flags().StringVar(&req.State, "state", "", "Service state")
	c.Flags().StringVar(&req.Output, "output", "", "Service output")
	c.Flags().StringVar(&eventTime, "time", "", "Event time")
	c.Flags().StringVarP(&req.Author, "author", "a", "", "Event author name")
	c.Flags().StringVarP(&req.Comment, "comment", "c", "", "Event comment")

	c.Flags().AddGoFlagSet(flag.CommandLine)
	logs.InitLogs()

	return c
}
