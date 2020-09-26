package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mazzegi/xpdf/xdoc"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: xpdf <input-file> <output-file>")
		os.Exit(1)
	}
	in := os.Args[1]
	var out string
	if len(os.Args) >= 3 {
		out = os.Args[2]
	} else {
		base := filepath.Base(in)
		ext := filepath.Ext(in)
		out = strings.TrimSuffix(base, ext) + ".pdf"
	}

	doc, err := xdoc.LoadFromFile(in)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}
	fmt.Printf("to out %q: %#v\n", out, doc)
}
