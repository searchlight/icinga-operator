package main

import (
	"log"
	"os"

	logs "github.com/appscode/go/log/golog"
	"github.com/appscode/searchlight/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdHostfacts(Version).Execute(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
