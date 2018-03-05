//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"log"
	"os"

	logs "github.com/appscode/go/log/golog"
	_ "github.com/appscode/kutil/extensions/v1beta1"
	"github.com/appscode/searchlight/pkg/cmds"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/client-go/kubernetes/fake"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdSearchlight().Execute(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
