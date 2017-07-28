//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"log"

	logs "github.com/appscode/log/golog"
	"github.com/appscode/searchlight/pkg/cmds"
	"os"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdSearchlight(Version).Execute(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
