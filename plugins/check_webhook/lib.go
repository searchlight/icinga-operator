/*
Copyright AppsCode Inc. and Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package check_webhook

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	api "go.searchlight.dev/icinga-operator/apis/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1"
	cs "go.searchlight.dev/icinga-operator/client/clientset/versioned/typed/monitoring/v1alpha1"
	"go.searchlight.dev/icinga-operator/pkg/icinga"
	"go.searchlight.dev/icinga-operator/plugins"

	"github.com/spf13/cobra"
	"gomodules.xyz/x/flags"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kmodules.xyz/client-go/tools/clientcmd"
)

const (
	totalFlag        = 20
	FlagCheckCommand = "check_command"
	FlagWebhookURL   = "webhook_url"
	stateUnknown     = 3
)

type webhookResp struct {
	Code    *int32      `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type param struct {
	key string
	val string
}

type plugin struct {
	client  v1alpha1.SearchlightPluginInterface
	options options
}

var _ plugins.PluginInterface = &plugin{}

func newPlugin(client v1alpha1.SearchlightPluginInterface, opts options) *plugin {
	return &plugin{client, opts}
}

func newPluginFromConfig(opts options) (*plugin, error) {
	config, err := clientcmd.BuildConfigFromContext(opts.kubeconfigPath, opts.contextName)
	if err != nil {
		return nil, err
	}
	monitoringClient, err := cs.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return newPlugin(monitoringClient.SearchlightPlugins(), opts), nil
}

type options struct {
	kubeconfigPath string
	contextName    string
	// options
	webhookURL   string
	checkCommand string
	params       []param
}

func (o *options) complete(cmd *cobra.Command) error {
	var err error

	o.kubeconfigPath, err = cmd.Flags().GetString(plugins.FlagKubeConfig)
	if err != nil {
		return err
	}
	o.contextName, err = cmd.Flags().GetString(plugins.FlagKubeConfigContext)
	if err != nil {
		return err
	}
	return nil
}

func (o *options) validate() error {
	return nil
}

func (p *plugin) Check() (icinga.State, interface{}) {
	opts := p.options

	sp, err := p.client.Get(context.TODO(), opts.checkCommand, metav1.GetOptions{})
	if err != nil {
		return icinga.Unknown, err
	}

	data := make(map[string]interface{})

	vars := sp.Spec.Arguments.Vars
	if vars != nil {
		for _, p := range opts.params {
			if p.key == "" || p.val == "" {
				continue
			}

			item, found := vars.Fields[p.key]
			if !found {
				return stateUnknown, fmt.Errorf(`var "%s" is not registered in SearchlightPlugin`, p.key)
			}

			switch item.Type {
			case api.VarTypeInteger:
				val, err := strconv.ParseInt(p.val, 10, 64)
				if err != nil {
					return icinga.Unknown, fmt.Errorf(`failed to parse value for key "%s" to Int64`, p.key)
				}
				data[p.key] = val
			case api.VarTypeNumber:
				val, err := strconv.ParseFloat(p.val, 64)
				if err != nil {
					return icinga.Unknown, fmt.Errorf(`failed to parse value for key "%s" to Float64`, p.key)
				}
				data[p.key] = val
			case api.VarTypeString:
				data[p.key] = p.val
			case api.VarTypeBoolean:
				val, err := strconv.ParseBool(p.val)
				if err != nil {
					return icinga.Unknown, fmt.Errorf(`failed to parse value for key "%s" to Bool`, p.key)
				}
				data[p.key] = val
			case api.VarTypeDuration:
				duration, err := time.ParseDuration(p.val)
				if err != nil {
					return icinga.Unknown, fmt.Errorf(`failed to parse value for key "%s" to Duration`, p.key)
				}
				data[p.key] = int64(duration.Nanoseconds() / 1000000)
			}
		}
	} else {
		for _, p := range opts.params {
			data[p.key] = p.val
		}
	}

	reqData, err := json.Marshal(data)
	if err != nil {
		return stateUnknown, err
	}

	b := strings.NewReader(string(reqData))
	resp, err := http.Post(opts.webhookURL, "application/json", b)
	if err != nil {
		return stateUnknown, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return stateUnknown, fmt.Sprintf("status code: %d", resp.StatusCode)
	}

	respDataByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return stateUnknown, err
	}

	var respData webhookResp
	if err := json.Unmarshal(respDataByte, &respData); err != nil {
		return stateUnknown, err
	}

	if respData.Code == nil {
		return stateUnknown, errors.New("can't identify State")
	}

	return icinga.State(*respData.Code), respData.Message
}

func NewCmd() *cobra.Command {
	opts := options{
		params: make([]param, totalFlag),
	}

	cmd := &cobra.Command{
		Use:   "check_webhook",
		Short: "Check webhook result",

		Run: func(cmd *cobra.Command, args []string) {
			flags.EnsureRequiredFlags(cmd, FlagCheckCommand, FlagWebhookURL)

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
			icinga.Output(plugin.Check())
		},
	}

	cmd.Flags().StringVar(&opts.webhookURL, FlagWebhookURL, "", "Call the webhook server using this URL")
	cmd.Flags().StringVar(&opts.checkCommand, FlagCheckCommand, "", "SearchlightPlugin name for this webhook check")

	for i := 0; i < totalFlag; i++ {
		cmd.Flags().StringVar(&opts.params[i].key, fmt.Sprintf("key.%d", i), "", "")
		cmd.Flags().StringVar(&opts.params[i].val, fmt.Sprintf("val.%d", i), "", "")
	}

	return cmd
}
