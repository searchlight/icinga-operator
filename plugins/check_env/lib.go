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

package check_env

import (
	"fmt"
	"os"

	"go.searchlight.dev/icinga-operator/pkg/icinga"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "check_env",
		Run: func(c *cobra.Command, args []string) {
			envList := os.Environ()
			fmt.Fprintln(os.Stdout, "Total ENV: ", len(envList))
			fmt.Fprintln(os.Stdout)
			for _, env := range envList {
				fmt.Fprintln(os.Stdout, env)
			}
			icinga.Output(icinga.OK, "A-OK")
		},
	}
	return cmd
}
