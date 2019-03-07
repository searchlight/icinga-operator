//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"log"

	_ "github.com/appscode/searchlight/apis/incidents/v1alpha1"
	"github.com/appscode/searchlight/pkg/cmds"
	_ "k8s.io/api/core/v1"
	_ "k8s.io/client-go/kubernetes/fake"
	_ "kmodules.xyz/client-go/extensions/v1beta1"
	"kmodules.xyz/client-go/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdSearchlight().Execute(); err != nil {
		log.Fatal(err)
	}
}
