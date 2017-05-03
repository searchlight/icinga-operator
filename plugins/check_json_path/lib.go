package check_json_path

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/appscode/go/flags"
	"github.com/appscode/go/net/httpclient"
	"github.com/appscode/searchlight/pkg/client/k8s"
	"github.com/appscode/searchlight/util"
	"github.com/spf13/cobra"
	rest "k8s.io/kubernetes/pkg/client/restclient"
)

type Request struct {
	URL             string
	Query           string
	Secret          string
	Namespace       string
	InClusterConfig bool
	Warning         string
	Critical        string
}

type AuthInfo struct {
	CertificateAuthorityData []byte `json:"certificate-authority-data,omitempty"`
	ClientCertificateData    []byte `json:"client-certificate-data,omitempty"`
	Token                    string `json:"token,omitempty"`
	Username                 string `json:"username,omitempty"`
	Password                 string `json:"password,omitempty"`
}

type JQ struct {
	J string `json:"j"`
	Q string `json:"q"`
}

const (
	auth = "auth"
)

func getData(req *Request) (string, error) {
	httpClient := httpclient.Default().WithBaseURL(req.URL)
	if req.Secret != "" {
		kubeClient, err := k8s.NewClient()
		if err != nil {
			return "", err
		}

		name := req.Secret
		namespace := req.Namespace

		secret, err := kubeClient.Client.Core().Secrets(namespace).Get(name)
		if err != nil {
			return "", err
		}

		secretData := new(AuthInfo)
		if data, found := secret.Data[auth]; found {
			if err := json.Unmarshal(data, secretData); err != nil {
				return "", err
			}
		}
		httpClient.WithBearerToken(secretData.Token)
		httpClient.WithBasicAuth(secretData.Username, secretData.Password)
		if secretData.CertificateAuthorityData != nil {
			httpClient.WithTLSConfig(secretData.ClientCertificateData, secretData.CertificateAuthorityData)
		}
	}
	if req.InClusterConfig {
		config, err := rest.InClusterConfig()
		if err != nil {
			return "", err
		}

		httpClient.WithBearerToken(config.BearerToken)
		httpClient.WithInsecureSkipVerify()
	}

	var respJson interface{}
	resp, err := httpClient.Call("GET", "", nil, &respJson, true)
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

func CheckJsonPath(req *Request) (util.IcingaState, interface{}) {
	jsonData, err := getData(req)
	if err != nil {
		return util.Unknown, err

	}

	jqData := &JQ{
		J: jsonData,
		Q: req.Query,
	}

	evalData, err := jqData.eval()
	if err != nil {
		return util.Unknown, "Invalid query. No data found"
	}

	evalDataByte, err := json.Marshal(evalData)
	if err != nil {
		return util.Unknown, err

	}

	evalDataString := string(evalDataByte)
	if req.Critical != "" {
		isCritical, err := checkResult(evalDataString, req.Critical)
		if err != nil {
			return util.Unknown, err
		}
		if isCritical {
			return util.Critical, fmt.Sprintf("%v", req.Critical)
		}
	}
	if req.Warning != "" {
		isWarning, err := checkResult(evalDataString, req.Warning)
		if err != nil {
			return util.Unknown, err
		}
		if isWarning {
			return util.Warning, fmt.Sprintf("%v", req.Warning)
		}
	}
	return util.Ok, "Response looks good"
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

			parts := strings.Split(icingaHost, "@")
			if len(parts) != 2 {
				fmt.Fprintln(os.Stdout, util.State[3], "Invalid icinga host.name")
				os.Exit(3)
			}
			req.Namespace = parts[1]

			flags.EnsureRequiredFlags(cmd, "url", "query")
			flags.EnsureAlterableFlags(cmd, "warning", "critical")
			util.Output(CheckJsonPath(&req))
		},
	}

	c.Flags().StringVarP(&icingaHost, "host", "H", "", "Icinga host name")
	c.Flags().StringVarP(&req.URL, "url", "u", "", "URL to get data")
	c.Flags().StringVarP(&req.Query, "query", "q", "", `JQ query`)
	c.Flags().StringVarP(&req.Secret, "secret", "s", "", `Kubernetes secret name`)
	c.Flags().BoolVar(&req.InClusterConfig, "in_cluster_config", false, `Use Kubernetes InCluserConfig`)
	c.Flags().StringVarP(&req.Warning, "warning", "w", "", `Warning JQ query which returns [true/false]`)
	c.Flags().StringVarP(&req.Critical, "critical", "c", "", `Critical JQ query which returns [true/false]`)
	return c
}
