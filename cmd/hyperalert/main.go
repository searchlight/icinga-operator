//go:generate stringer -type=State ../../pkg/icinga/types.go
package main

import (
	"os"

	"github.com/appscode/searchlight/plugins/hyperalert"
)

func main() {
	if err := hyperalert.NewCmd().Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
