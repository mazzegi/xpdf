package xpdf

import "fmt"

func Logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
