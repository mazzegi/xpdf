package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mazzegi/xpdf"
	"github.com/mazzegi/xpdf/engine"
	"github.com/mazzegi/xpdf/font"
	"github.com/mazzegi/xpdf/xdoc"
)

func main() {
	var in string
	if len(os.Args) < 2 {
		// fmt.Println("usage: xpdf <in> <optional:out>")
		// os.Exit(1)
		//in = "../../examples/measure.xml"
		//in = "../../examples/table1.xml"
		in = "../../examples/doc1.xml"
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

	fonts := font.NewDirectory()
	fonts.Register(font.Descriptor{
		Name:     "chin",
		Style:    font.Regular,
		FilePath: "fonts/DroidSansFallback.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu",
		Style:    font.Regular,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu",
		Style:    font.Bold,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu",
		Style:    font.Italic,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSans-Oblique.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu-serif",
		Style:    font.Regular,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSerif.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu-serif",
		Style:    font.Bold,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSerif-Bold.ttf",
	})
	fonts.Register(font.Descriptor{
		Name:     "dejavu-serif",
		Style:    font.Italic,
		FilePath: "/usr/share/fonts/truetype/dejavu/DejaVuSerif-Italic.ttf",
	})

	engine, err := engine.NewFPDF(fonts, doc)
	if err != nil {
		fmt.Println("ERROR create fpdf-engine:", err)
		os.Exit(3)
	}

	p := xpdf.NewProcessor(engine, doc)
	err = p.Process(outF)
	if err != nil {
		fmt.Println("ERROR processing input:", err)
		os.Exit(4)
	}
}
