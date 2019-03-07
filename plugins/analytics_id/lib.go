package analytics_id

import (
	"fmt"

	"github.com/spf13/cobra"
	"kmodules.xyz/client-go/tools/analytics"
)

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "analytics_id",
		Run: func(c *cobra.Command, args []string) {
			fmt.Print(analytics.ClientID())
		},
	}
	return cmd
}
