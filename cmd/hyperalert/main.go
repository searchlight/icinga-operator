//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"os"

	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/plugins/hyperalert"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := hyperalert.NewCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
