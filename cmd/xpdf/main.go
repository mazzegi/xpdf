package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mazzegi/xpdf"
	"github.com/mazzegi/xpdf/engine"
	"github.com/mazzegi/xpdf/xdoc"
)

func main() {
	var in string
	if len(os.Args) < 2 {
		fmt.Println("usage: xpdf <in> <optional:out>")
		os.Exit(1)
	} else {
		in = os.Args[1]
	}
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
		fmt.Println("ERROR compiling input:", err)
		os.Exit(2)
	}
	outF, err := os.Create(out)
	if err != nil {
		fmt.Println("ERROR create output-file:", err)
		os.Exit(2)
	}
	defer outF.Close()

	engine := engine.NewFPDF(doc)
	p := xpdf.NewProcessor(engine, doc)
	err = p.Process(outF)
	if err != nil {
		fmt.Println("ERROR processing input:", err)
		os.Exit(3)
	}
}