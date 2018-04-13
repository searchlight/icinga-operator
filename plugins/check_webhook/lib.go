package check_webhook

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/appscode/searchlight/pkg/icinga"
	"github.com/spf13/cobra"
)

const (
	totalFlag    = 20
	flagURL      = "url"
	stateUnknown = 3
)

type webhookResp struct {
	Code    *int32      `json:"code,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

func NewCmd() *cobra.Command {

	type flag struct {
		key string
		val string
	}

	flags := make([]flag, totalFlag)

	cmd := &cobra.Command{
		Use:   "check_webhook",
		Short: "Check webhook result",

		Run: func(cmd *cobra.Command, args []string) {
			var url string

			data := make(map[string]string)

			for _, v := range flags {
				if v.key == "" {
					continue
				}

				if v.key == flagURL {
					url = v.val
				}
				data[v.key] = v.val
			}

			if url == "" {
				icinga.Output(stateUnknown, errors.New("webhook URL missing"))
			}

			reqData, err := json.Marshal(data)
			if err != nil {
				icinga.Output(stateUnknown, err)
			}

			b := strings.NewReader(string(reqData))
			resp, err := http.Post(url, "application/json", b)
			if err != nil {
				icinga.Output(stateUnknown, err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != 200 {
				icinga.Output(stateUnknown, fmt.Sprintf("status code: %d", resp.StatusCode))
			}

			respDataByte, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				icinga.Output(stateUnknown, err)
			}

			var respData webhookResp
			if err := json.Unmarshal(respDataByte, &respData); err != nil {
				icinga.Output(stateUnknown, err)
			}

			if respData.Code == nil {
				icinga.Output(stateUnknown, errors.New("can't identify State"))
			}

			icinga.Output(icinga.State(*respData.Code), respData.Message)
		},
	}

	for i := 0; i < totalFlag; i++ {
		cmd.Flags().StringVar(&flags[i].key, fmt.Sprintf("key.%d", i), "", "")
		cmd.Flags().StringVar(&flags[i].val, fmt.Sprintf("val.%d", i), "", "")
	}

	return cmd
}
