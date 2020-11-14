package xdoc

import (
	"encoding/xml"
)

// type Block interface {
// 	DecodeAttrs(attrs []xml.Attr) error
// 	MutatedStyles(cs style.Classes, styles style.Styles) style.Styles
// 	MutatedStylesWithSelector(sel string, cs style.Classes, styles style.Styles) style.Styles
// }

// type Blocks struct {
// 	BS []Block
// }

// func (bs *Blocks) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
// 	for {
// 		token, err := d.Token()
// 		if err != nil {
// 			return err
// 		}
// 		switch t := token.(type) {
// 		case xml.EndElement:
// 			if t == start.End() {
// 				return nil
// 			}
// 		case xml.StartElement:
// 			b, err := registry.DecodeBlock(d, t)
// 			if err != nil {
// 				return errors.Wrapf(err, "decode token %v", token)
// 			}
// 			bs.BS = append(bs.BS, b)
// 		case xml.CharData:
// 			v := string(t)
// 			if v != "" {
// 				bs.BS = append(bs.BS, &TextBlock{
// 					Text: string(t),
// 				})
// 			}
// 		}
// 	}
// }

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
