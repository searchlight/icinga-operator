package check_any_cert

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type request struct {
	masterURL      string
	kubeconfigPath string

	Namespace  string
	Selector   string
	SecretName string
	SecretKey  []string
	Count      int

	warning  time.Duration
	critical time.Duration
}

func getCertSecrets(req *request) (*core.SecretList, error) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return nil, err
	}

	kubeClient := kubernetes.NewForConfigOrDie(config)

	secretList := &core.SecretList{}

	if req.SecretName != "" {
		secret, err := kubeClient.CoreV1().Secrets(req.Namespace).Get(req.SecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		secretList.Items = append(secretList.Items, *secret)
	} else {
		secretList, err = kubeClient.CoreV1().Secrets(req.Namespace).List(metav1.ListOptions{
			LabelSelector: req.Selector,
		})
		if err != nil {
			return nil, err
		}
	}
	return secretList, nil
}

func checkNotAfter(cert *x509.Certificate, req *request) (icinga.State, time.Duration, bool) {
	remaining := cert.NotAfter.Sub(time.Now())

	if remaining.Seconds() < req.critical.Seconds() {
		return icinga.Critical, remaining, false
	}

	if remaining.Seconds() < req.warning.Seconds() {
		return icinga.Warning, remaining, false
	}

	return icinga.OK, remaining, true
}

func checkCert(data []byte, secret *core.Secret, key string, req *request) (icinga.State, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return icinga.Unknown, fmt.Errorf(
			`failed to parse certificate for key "%s" in Secret "%s/%s"`,
			key, secret.Namespace, secret.Name,
		)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return icinga.Unknown, fmt.Errorf(
			`failed to parse certificate for key "%s" in Secret "%s/%s"`,
			key, secret.Namespace, secret.Name,
		)
	}

	if state, remaining, ok := checkNotAfter(cert, req); !ok {
		return state, fmt.Errorf(
			`certificate for key "%s" in Secret "%s/%s" will be expired within %v hours`,
			key, secret.Namespace, secret.Name, remaining.Hours(),
		)
	}

	return icinga.OK, nil
}

func checkCertPerSecretKey(secret *core.Secret, req *request) (icinga.State, error) {
	for _, key := range req.SecretKey {
		data, ok := secret.Data[key]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "%s" not found in Secret "%s/%s"`, key, secret.Namespace, secret.Name)
		}

		if state, err := checkCert(data, secret, key, req); err != nil {
			return state, err
		}
	}

	if len(req.SecretKey) == 0 && secret.Type == core.SecretTypeTLS {
		data, ok := secret.Data["tls.crt"]
		if !ok {
			return icinga.Warning, fmt.Errorf(`key "tls.crt" not found in Secret "%s/%s"`, secret.Namespace, secret.Name)
		}

		if state, err := checkCert(data, secret, "tls.crt", req); err != nil {
			return state, err
		}
	}

	return icinga.OK, nil
}

func checkAnyCert(req *request) (icinga.State, interface{}) {
	secretList, err := getCertSecrets(req)
	if err != nil {
		return icinga.Unknown, err
	}

	for _, secret := range secretList.Items {
		if state, err := checkCertPerSecretKey(&secret, req); err != nil {
			return state, err
		}
	}

	return icinga.OK, fmt.Sprintf("Certificate expirity check is succeeded")
}

func NewCmd() *cobra.Command {
	var req request
	var icingaHost string
	var commaSeparatedKeys string

	cmd := &cobra.Command{
		Use:     "check_any_cert",
		Short:   "Check Certificate expire date",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host.name")
				os.Exit(3)
			}
			if host.Type != icinga.TypeCluster {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host type")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			req.SecretKey = make([]string, 0)
			keys := strings.Split(commaSeparatedKeys, ",")
			for _, key := range keys {
				req.SecretKey = append(req.SecretKey, strings.Trim(key, " "))
			}
			icinga.Output(checkAnyCert(&req))
		},
	}

	cmd.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='")
	cmd.Flags().StringVarP(&req.SecretName, "secretName", "s", "", "Name of secret from where certificates are checked")
	cmd.Flags().StringVarP(&commaSeparatedKeys, "secretKey", "k", "", "Name of secret key where certificates are kept")
	cmd.Flags().DurationVarP(&req.warning, "warning", "w", time.Hour*360, `Remaining duration for warning state. [Default: 360h]`)
	cmd.Flags().DurationVarP(&req.critical, "critical", "c", time.Hour*120, `Remaining duration for critical state. [Default: 120h]`)
	return cmd
}
