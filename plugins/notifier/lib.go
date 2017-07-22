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
	"github.com/appscode/log"
	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

type Request struct {
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

func getLoader(client clientset.Interface) (envconfig.LoaderFunc, error) {
	secretName := os.Getenv(icinga.ICINGA_NOTIFIER_SECRET_NAME)
	secretNamespace := util.OperatorNamespace()

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

func sendNotification(req *Request) {
	client, err := util.NewClient()
	if err != nil {
		log.Fatalln(err)
	}

	host, err := icinga.ParseHost(req.HostName)
	if err != nil {
		log.Fatalln(err)
	}

	alert, err := host.GetAlert(client.ExtClient, req.AlertName)
	if err != nil {
		log.Fatalln(err)
	}

	loader, err := getLoader(client.Client)
	if err != nil {
		log.Fatalln(err)
	}

	receivers := alert.GetReceivers()

	for _, receiver := range receivers {
		if receiver.State != req.State || len(receiver.To) == 0 {
			continue
		}
		notifyVia, err := unified.LoadVia(receiver.Method, loader)
		if err != nil {
			log.Errorln(err)
			continue
		}

		switch n := notifyVia.(type) {
		case notify.ByEmail:
			subject := "Notification"
			if sub, found := subjectMap[req.Type]; found {
				subject = sub
			}
			var mailBody string
			mailBody, err = RenderMail(alert, req)
			if err != nil {
				break
			}
			err = n.To(receiver.To[0], receiver.To[1:]...).WithSubject(subject).WithBody(mailBody).Send()
		case notify.BySMS:
			var smsBody string
			smsBody, err = RenderSMS(alert, req)
			if err != nil {
				break
			}
			err = n.To(receiver.To[0], receiver.To[1:]...).WithBody(smsBody).Send()
		case notify.ByChat:
			var smsBody string
			smsBody, err = RenderSMS(alert, req)
			if err != nil {
				break
			}
			err = n.To(receiver.To[0], receiver.To[1:]...).WithBody(smsBody).Send()
		}

		if err != nil {
			log.Errorln(err)
		} else {
			log.Debug(fmt.Sprintf("Notification sent using %s", receiver.Method))
		}
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

	c.Flags().StringVarP(&req.HostName, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.AlertName, "alert", "A", "", "Kubernetes alert object name")
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
