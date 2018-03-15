package check_any_cert

import (
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/cert"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	Namespace  string
	Selector   string
	SecretName string
	SecretKey  []string

	Warning  time.Duration
	Critical time.Duration
}

type CertContext struct {
	Client corev1.SecretInterface
}

func NewCertContext(req *Request) (*CertContext, error) {
	config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
	if err != nil {
		return nil, err
	}
	return &CertContext{
		Client: kubernetes.NewForConfigOrDie(config).CoreV1().Secrets(req.Namespace),
	}, nil
}

func (cc *CertContext) getCertSecrets(req *Request) ([]core.Secret, error) {
	if req.SecretName != "" {
		var secret *core.Secret
		secret, err := cc.Client.Get(req.SecretName, metav1.GetOptions{})
		if err != nil {
			return nil, err
		}
		return []core.Secret{*secret}, nil
	} else {
		secretList, err := cc.Client.List(metav1.ListOptions{
			LabelSelector: req.Selector,
		})
		if err != nil {
			return nil, err
		}
		return secretList.Items, nil
	}
	return nil, nil
}

func checkNotAfter(cert *x509.Certificate, req *Request) (icinga.State, time.Duration, bool) {
	remaining := cert.NotAfter.Sub(time.Now())
	if remaining.Seconds() < req.Critical.Seconds() {
		return icinga.Critical, remaining, false
	}

	if remaining.Seconds() < req.Warning.Seconds() {
		return icinga.Warning, remaining, false
	}

	return icinga.OK, remaining, true
}

func checkCert(data []byte, secret *core.Secret, key string, req *Request) (icinga.State, error) {
	certs, err := cert.ParseCertsPEM(data)
	if err != nil {
		return icinga.Unknown, fmt.Errorf(
			`failed to parse certificate for key "%s" in Secret "%s/%s"`,
			key, secret.Namespace, secret.Name,
		)
	}

	for _, cert := range certs {
		if state, remaining, ok := checkNotAfter(cert, req); !ok {
			return state, fmt.Errorf(
				`certificate for key "%s" in Secret "%s/%s" will be expired within %v hours`,
				key, secret.Namespace, secret.Name, remaining.Hours(),
			)
		}
	}
	return icinga.OK, nil
}

func checkCertPerSecretKey(secret *core.Secret, req *Request) (icinga.State, error) {
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

func (cc *CertContext) CheckAnyCert(req *Request) (icinga.State, interface{}) {
	secretList, err := cc.getCertSecrets(req)
	if err != nil {
		return icinga.Unknown, err
	}

	for _, secret := range secretList {
		if state, err := checkCertPerSecretKey(&secret, req); err != nil {
			return state, err
		}
	}

	return icinga.OK, fmt.Sprintf("Certificate expirity check is succeeded")
}

func NewCmd() *cobra.Command {
	var req *Request
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

			context, err := NewCertContext(req)
			if err != nil {
				icinga.Output(icinga.Unknown, err)
			}
			icinga.Output(context.CheckAnyCert(req))
		},
	}

	cmd.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	cmd.Flags().StringVarP(&req.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='")
	cmd.Flags().StringVarP(&req.SecretName, "secretName", "s", "", "Name of secret from where certificates are checked")
	cmd.Flags().StringVarP(&commaSeparatedKeys, "secretKey", "k", "", "Name of secret key where certificates are kept")
	cmd.Flags().DurationVarP(&req.Warning, "warning", "w", time.Hour*360, `Remaining duration for Warning state. [Default: 360h]`)
	cmd.Flags().DurationVarP(&req.Critical, "critical", "c", time.Hour*120, `Remaining duration for Critical state. [Default: 120h]`)
	return cmd
}
