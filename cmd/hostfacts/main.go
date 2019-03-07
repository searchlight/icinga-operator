package main

import (
	"log"

	"github.com/appscode/searchlight/pkg/cmds"
	"kmodules.xyz/client-go/logs"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := cmds.NewCmdHostfacts().Execute(); err != nil {
		log.Fatal(err)
	}
}
