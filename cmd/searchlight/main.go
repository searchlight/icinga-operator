//go:generate stringer -type=IcingaState ../../pkg/icinga/types.go
package main

import (
	"log"

	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdSearchlight(Version).Execute(); err != nil {
		log.Fatal(err)
	}
}
