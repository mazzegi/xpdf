package xdoc

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"

	"github.com/mazzegi/xpdf/style"
	"github.com/pkg/errors"
)

func Load(r io.Reader) (*Document, error) {
	doc := &Document{}
	err := xml.NewDecoder(r).Decode(doc)
	if err != nil {
		return nil, err
	}
	doc.styleClasses, err = style.DecodeClasses(bytes.NewBufferString(doc.Style))
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func LoadFromFile(file string) (*Document, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Errorf("open (%s)", file)
	}
	defer f.Close()
	return Load(f)
}

type Orientation string

const (
	OrientationPortrait  Orientation = "portrait"
	OrientationLandscape Orientation = "landscape"
)

type PaperFormat string

const (
	FormatA3     PaperFormat = "a3"
	FormatA4     PaperFormat = "a4"
	FormatA5     PaperFormat = "a5"
	FormatLetter PaperFormat = "letter"
	FormatLegal  PaperFormat = "legal"
)

type Document struct {
	XMLName      xml.Name     `xml:"document"`
	Meta         Meta         `xml:"meta"`
	Page         Page         `xml:"page"`
	Style        string       `xml:"style"`
	Header       Instructions `xml:"header"`
	Footer       Instructions `xml:"footer"`
	Body         Instructions `xml:"body"`
	styleClasses style.Classes
}

type Meta struct {
	XMLName xml.Name `xml:"meta"`
	Author  string   `xml:"author"`
	Creator string   `xml:"creator"`
	Subject string   `xml:"subject"`
}

type Margins struct {
	XMLName xml.Name `xml:"margins"`
	Left    float64  `xml:"left"`
	Top     float64  `xml:"top"`
	Right   float64  `xml:"right"`
	Bottom  float64  `xml:"bottom"`
}

type Page struct {
	XMLName     xml.Name    `xml:"page"`
	Orientation Orientation `xml:"orientation"`
	Format      PaperFormat `xml:"format"`
	Margins     Margins     `xml:"margins"`
}

func (doc *Document) StyleClasses() style.Classes {
	return doc.styleClasses
}
