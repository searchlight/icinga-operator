package check_json_path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"

	"github.com/appscode/envconfig"
	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/appscode/searchlight/pkg/util"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

type Request struct {
	URL             string
	Secret          string
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

const (
	auth = "auth"
)

func getData(req *Request) (string, error) {
	var hc *httpclient.Client

	if req.InClusterConfig {
		kubeClient, err := util.NewClient()
		if err != nil {
			return "", err
		}
		cc := kubeClient.Client.CoreV1().RESTClient().(*rest.RESTClient)
		hc = httpclient.New(cc.Client, nil, nil)
	} else {
		hc = httpclient.Default().WithBaseURL(req.URL)

		if req.Secret != "" {
			kubeClient, err := util.NewClient()
			if err != nil {
				return "", err
			}
			secret, err := kubeClient.Client.CoreV1().Secrets(req.Namespace).Get(req.Secret, metav1.GetOptions{})
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

func checkResult(evalDataString, checkQuery string) (bool, error) {
	jqData := &JQ{
		J: string(evalDataString),
		Q: checkQuery,
	}
	result, err := jqData.eval()
	if err != nil {
		return false, err
	}
	if reflect.TypeOf(result).Kind() != reflect.Bool {
		return false, fmt.Errorf("Invalid check query: %v", checkQuery)
	}
	return result.(bool), nil
}

func CheckJsonPath(req *Request) (icinga.State, interface{}) {
	jsonData, err := getData(req)
	if err != nil {
		return icinga.UNKNOWN, err

	}

	//jqData := &JQ{
	//	J: jsonData,
	//	Q: req.Query,
	//}
	//
	//evalData, err := jqData.eval()
	//if err != nil {
	//	return icinga.UNKNOWN, "Invalid query. No data found"
	//}
	//
	//evalDataByte, err := json.Marshal(evalData)
	//if err != nil {
	//	return icinga.UNKNOWN, err
	//
	//}

	//evalDataString := string(evalDataByte)
	if req.Critical != "" {
		isCritical, err := checkResult(jsonData, req.Critical)
		if err != nil {
			return icinga.UNKNOWN, err
		}
		if isCritical {
			return icinga.CRITICAL, fmt.Sprintf("%v", req.Critical)
		}
	}
	if req.Warning != "" {
		isWarning, err := checkResult(jsonData, req.Warning)
		if err != nil {
			return icinga.UNKNOWN, err
		}
		if isWarning {
			return icinga.WARNING, fmt.Sprintf("%v", req.Warning)
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
				fmt.Fprintln(os.Stdout, icinga.WARNING, "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = host.AlertNamespace

			flags.EnsureRequiredFlags(cmd, "url", "query")
			flags.EnsureAlterableFlags(cmd, "warning", "critical")
			icinga.Output(CheckJsonPath(&req))
		},
	}

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.URL, "url", "u", "", "URL to get data")
	c.Flags().StringVarP(&req.Secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().BoolVar(&req.InClusterConfig, "in_cluster_config", false, `Use Kubernetes InCluserConfig`)
	c.Flags().StringVarP(&req.Warning, "warning", "w", "", `Warning JQ query which returns [true/false]`)
	c.Flags().StringVarP(&req.Critical, "critical", "c", "", `Critical JQ query which returns [true/false]`)
	return c
}
