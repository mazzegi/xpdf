package xdoc

import (
	"encoding/xml"

	"github.com/mazzegi/xpdf/style"
	"github.com/pkg/errors"
)

type Instruction interface {
	DecodeAttrs(attrs []xml.Attr) error
	MutatedStyles(cs style.Classes, styles style.Styles) style.Styles
	MutatedStylesWithSelector(sel string, cs style.Classes, styles style.Styles) style.Styles
}

type Instructions struct {
	Styled
	iss []Instruction
}

func (is *Instructions) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	err := is.DecodeAttrs(start.Attr)
	if err != nil {
		return err
	}
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t == start.End() {
				return nil
			}
		case xml.StartElement:
			i, err := instructionRegistry.Decode(d, t)
			if err != nil {
				return errors.Wrapf(err, "decode token %v", token)
			}
			is.iss = append(is.iss, i)
		}
	}
}

type Font struct {
	Styled
	XMLName xml.Name `xml:"font"`
}

type LineFeed struct {
	NoStyles
	XMLName xml.Name `xml:"lf"`
	Lines   float64  `xml:"lines,attr"`
}

type SetX struct {
	NoStyles
	XMLName xml.Name `xml:"setx"`
	X       float64  `xml:"x,attr"`
}

type SetY struct {
	NoStyles
	XMLName xml.Name `xml:"sety"`
	Y       float64  `xml:"y,attr"`
}

type SetXY struct {
	NoStyles
	XMLName xml.Name `xml:"setxy"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
}

type Box struct {
	Styled
	XMLName xml.Name `xml:"box"`
	Text    string   `xml:",chardata"`
}

type Text struct {
	Styled
	XMLName xml.Name `xml:"text"`
	Text    string   `xml:",chardata"`
}

type Image struct {
	Styled
	XMLName xml.Name `xml:"image"`
	Source  string   `xml:",chardata"`
}
