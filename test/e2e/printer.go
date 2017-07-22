package e2e

import "fmt"

func PrintSeparately(a ...interface{}) {
	fmt.Println()
	fmt.Println(a...)
	fmt.Println()
}
