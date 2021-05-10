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

package cmds

import (
	"encoding/json"
	"fmt"
	"time"

	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/spf13/cobra"
	"gomodules.xyz/x/flags"
)

func NewCmdConfigure() *cobra.Command {
	mgr := &icinga.Configurator{
		Expiry: 10 * 365 * 24 * time.Hour,
	}
	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Generate icinga configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags.SetLogLevel(4)

			cfg, err := mgr.LoadConfig(func(key string) (string, bool) {
				return "", false
			})
			if err != nil {
				return err
			}
			bytes, err := json.MarshalIndent(cfg, "", " ")
			if err != nil {
				return err
			}
			fmt.Println(string(bytes))
			return err
		},
	}

	cmd.Flags().StringVarP(&mgr.ConfigRoot, "config-dir", "s", mgr.ConfigRoot, "Path to directory containing icinga2 config. This should be an emptyDir inside Kubernetes.")

	return cmd
}
