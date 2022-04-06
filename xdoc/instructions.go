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
	ISS []Instruction
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
			i, err := registry.DecodeInstruction(d, t)
			if err != nil {
				return errors.Wrapf(err, "decode token %v", t.Name.Local)
			}
			is.ISS = append(is.ISS, i)
		case xml.CharData:
			v := string(t)
			if v != "" {
				is.ISS = append(is.ISS, &TextBlock{
					Text: string(t),
				})
			}
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

type Box struct {
	Styled
	XMLName xml.Name `xml:"box"`
	//Text    string   `xml:",chardata"`
	Instructions
}

type Text struct {
	Styled
	XMLName xml.Name `xml:"text"`
	Instructions
}

type Image struct {
	Styled
	XMLName xml.Name `xml:"image"`
	Source  string   `xml:",chardata"`
}

type TextBlock struct {
	NoStyles
	Text string
}

type Paragraph struct {
	Styled
	XMLName xml.Name `xml:"p"`
	Text    string   `xml:",chardata"`
}

type LineBreak struct {
	NoStyles
	XMLName xml.Name `xml:"br"`
}

type PageBreak struct {
	NoStyles
	XMLName xml.Name `xml:"newpage"`
}
