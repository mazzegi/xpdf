package main

import (
	"fmt"
	"os"

	"github.com/mazzegi/xpdf/xdoc"
)

func main() {
	var in string
	if len(os.Args) < 2 {
		in = "../../examples/table1.xml"
	} else {
		in = os.Args[1]
	}
	// var out string
	// if len(os.Args) >= 3 {
	// 	out = os.Args[2]
	// } else {
	// 	base := filepath.Base(in)
	// 	ext := filepath.Ext(in)
	// 	out = strings.TrimSuffix(base, ext) + ".pdf"
	// }

	doc, err := xdoc.LoadFromFile(in)
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(2)
	}
	desc := xdoc.Describe(doc)
	fmt.Printf("%s\n", desc.Dump(xdoc.DescribeXML))
}
