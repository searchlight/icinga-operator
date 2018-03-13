package check_json_path

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Request struct {
	masterURL      string
	kubeconfigPath string

	URL             string
	SecretName      string
	Namespace       string
	InClusterConfig bool
	Warning         string
	Critical        string
}

type AuthInfo struct {
	Username           string `envconfig:"USERNAME"`
	Password           string `envconfig:"PASSWORD"`
	Token              string `envconfig:"TOKEN"`
	CACertData         string `envconfig:"CA_CERT_DATA"`
	ClientCertData     string `envconfig:"CLIENT_CERT_DATA"`
	ClientKeyData      string `envconfig:"CLIENT_KEY_DATA"`
	InsecureSkipVerify bool   `envconfig:"INSECURE_SKIP_VERIFY"`
}

type JQ struct {
	J string `json:"j"`
	Q string `json:"q"`
}

func getData(req *Request) (string, error) {
	var hc *httpclient.Client

	if req.InClusterConfig {
		config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
		if err != nil {
			return "", err
		}
		kubeClient := kubernetes.NewForConfigOrDie(config)
		cc := kubeClient.CoreV1().RESTClient().(*rest.RESTClient)
		hc = httpclient.New(cc.Client, nil, nil)
	} else {
		hc = httpclient.Default().WithBaseURL(req.URL)
		if req.URL == "" {
			return "", errors.New("Missing URL")
		}

		if req.SecretName != "" {
			config, err := clientcmd.BuildConfigFromFlags(req.masterURL, req.kubeconfigPath)
			if err != nil {
				return "", err
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)
			secret, err := kubeClient.CoreV1().Secrets(req.Namespace).Get(req.SecretName, metav1.GetOptions{})
			if err != nil {
				return "", err
			}
			var au AuthInfo
			err = envconfig.Load("", &au, func(key string) (string, bool) {
				v, ok := secret.Data[key]
				if !ok {
					return "", false
				}
				return string(v), true
			})
			if err != nil {
				return "", err
			}
			hc = hc.WithBasicAuth(au.Username, au.Password).WithBearerToken(au.Token)
			if au.CACertData != "" {
				if au.ClientCertData != "" && au.ClientKeyData != "" {
					hc = hc.WithTLSConfig([]byte(au.CACertData), []byte(au.ClientKeyData), []byte(au.ClientKeyData))
				} else {
					hc = hc.WithTLSConfig([]byte(au.CACertData))
				}
			}
			if au.InsecureSkipVerify {
				hc = hc.WithInsecureSkipVerify()
			}
		}
	}

	var respJson interface{}
	resp, err := hc.Call("GET", "", nil, &respJson, true)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Invalid status_code %v", resp.StatusCode)
	}

	data, err := json.Marshal(respJson)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (j *JQ) eval() (res interface{}, err error) {
	cmd := exec.Command("jq", j.Q)
	cmd.Stdin = bytes.NewBufferString(j.J)

	var cmdOut []byte
	if cmdOut, err = cmd.Output(); err != nil {
		return

	}
	err = json.Unmarshal(cmdOut, &res)
	return
}

func checkResult(response, query string) (bool, error) {
	jqData := &JQ{
		J: response,
		Q: query,
	}
	result, err := jqData.eval()
	if err != nil {
		return false, err
	}
	if reflect.TypeOf(result).Kind() != reflect.Bool {
		return false, fmt.Errorf("invalid check query: %v", query)
	}
	return result.(bool), nil
}

func CheckJsonPath(req *Request) (icinga.State, interface{}) {
	jsonData, err := getData(req)
	if err != nil {
		return icinga.Unknown, err
	}

	if req.Critical != "" {
		isCritical, err := checkResult(jsonData, req.Critical)
		if err != nil {
			return icinga.Unknown, err
		}
		if isCritical {
			return icinga.Critical, fmt.Sprintf("%v", req.Critical)
		}
	}
	if req.Warning != "" {
		isWarning, err := checkResult(jsonData, req.Warning)
		if err != nil {
			return icinga.Unknown, err
		}
		if isWarning {
			return icinga.Warning, fmt.Sprintf("%v", req.Warning)
		}
	}
	return icinga.OK, "Response looks good"
}

func NewCmd() *cobra.Command {
	var req Request
	var icingaHost string

	c := &cobra.Command{
		Use:     "check_json_path",
		Short:   "Check Json Object",
		Example: "",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, "host")

			host, err := icinga.ParseHost(icingaHost)
			if err != nil {
				fmt.Fprintln(os.Stdout, icinga.Warning, "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			flags.EnsureAlterableFlags(cmd, "warning", "critical")
			icinga.Output(CheckJsonPath(&req))
		},
	}

	c.Flags().StringVar(&req.masterURL, "master", req.masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	c.Flags().StringVar(&req.kubeconfigPath, "kubeconfig", req.kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.URL, "url", "u", "", "URL to get data")
	c.Flags().StringVarP(&req.SecretName, "secretName", "s", "", `Kubernetes secret name`)
	c.Flags().BoolVar(&req.InClusterConfig, "inClusterConfig", false, `Use Kubernetes InCluserConfig`)
	c.Flags().StringVarP(&req.Warning, "warning", "w", "", `Warning JQ query which returns [true/false]`)
	c.Flags().StringVarP(&req.Critical, "critical", "c", "", `Critical JQ query which returns [true/false]`)
	return c
}
