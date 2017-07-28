package hello

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	c := &cobra.Command{
		Use: "env",
		Run: func(cmd *cobra.Command, args []string) {
			envList := os.Environ()
			fmt.Fprintln(os.Stdout, "Total ENV: ", len(envList))
			fmt.Fprintln(os.Stdout)
			for _, env := range envList {
				fmt.Fprintln(os.Stdout, env)
			}
		},
	}
	return c
}
