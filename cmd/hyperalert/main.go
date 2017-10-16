//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"os"

	logs "github.com/appscode/go/log/golog"
	_ "github.com/appscode/searchlight/client/scheme"
	"github.com/appscode/searchlight/plugins/hyperalert"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/client-go/kubernetes/fake"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()
	if err := hyperalert.NewCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
