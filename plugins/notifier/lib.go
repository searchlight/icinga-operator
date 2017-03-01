package notifier

import (
	"flag"
	"fmt"
	"os"
	"time"

	api "github.com/appscode/api/kubernetes/v1beta1"
	"github.com/appscode/client"
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
)

const (
	appscodeConfigPath = "/var/run/config/appscode/"
	appscodeSecretPath = "/var/run/secrets/appscode/"

	notifyVia = "NOTIFY_VIA"
)

type Secret struct {
	Namespace string `json:"namespace"`
	Token     string `json:"token"`
}

func notifyViaAppsCode(req *api.IncidentNotifyRequest) error {
	cluster_uid, err := io.ReadFile(appscodeConfigPath + "cluster-uid")
	if err != nil {
		return err
	}
	req.KubernetesCluster = cluster_uid

	grpc_endpoint, err := io.ReadFile(appscodeConfigPath + "appscode-api-grpc-endpoint")
	if err != nil {
		return err
	}

	apiOptions := client.NewOption(grpc_endpoint)

	var secretData Secret
	if err := io.ReadFileAs(appscodeSecretPath+"api-token", &secretData); err != nil {
		return err
	}

	apiOptions = apiOptions.BearerAuth(secretData.Namespace, secretData.Token)
	conn, err := client.New(apiOptions)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.Kubernetes().V1beta1().Incident().Notify(conn.Context(), req); err != nil {
		return err
	}

	return nil
}

func sendNotification(req *api.IncidentNotifyRequest) {
	if err := notifyViaAppsCode(req); err != nil {
		log.Debug(err)
	} else {
		log.Debug("Notification sent via AppsCode")
		os.Exit(0)
	}

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
	var req api.IncidentNotifyRequest
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
