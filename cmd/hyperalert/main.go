//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"os"

	"github.com/appscode/searchlight/plugins/hyperalert"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/client-go/kubernetes/fake"
	"kmodules.xyz/client-go/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := hyperalert.NewCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
