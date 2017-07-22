package main

import (
	"fmt"
	"log"
	"os"

	"github.com/appscode/go/runtime"
	"github.com/appscode/searchlight/pkg/cmds"
	"github.com/appscode/searchlight/plugins/hyperalert"
	"github.com/spf13/cobra/doc"
)

// ref: https://github.com/spf13/cobra/blob/master/doc/md_docs.md
func main() {
	haCmd := hyperalert.NewCmd()
	dir := runtime.GOPath() + "/src/github.com/appscode/searchlight/docs/reference/hyperalert"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	doc.GenMarkdownTree(haCmd, dir)

	opCmd := cmds.NewCmdSearchlight("")
	dir = runtime.GOPath() + "/src/github.com/appscode/searchlight/docs/reference/searchlight"
	fmt.Printf("Generating cli markdown tree in: %v\n", dir)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatal(err)
	}
	doc.GenMarkdownTree(opCmd, dir)
}
