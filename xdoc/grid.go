package xdoc

import (
	"encoding/xml"
	"strings"

	"github.com/pkg/errors"
)

type Grid struct {
	Styled
	XMLName xml.Name    `xml:"grid"`
	Rows    []*GridRow  `xml:"rows>gr"`
	Parts   []*GridPart `xml:"parts>part"`
}

type GridRow struct {
	Styled
	XMLName xml.Name `xml:"gr"`
	Areas   []string
}

type GridPart struct {
	Styled
	XMLName xml.Name `xml:"part"`
	Area    string   `xml:"area,attr"`
	Instructions
}

func (g *Grid) Validate() error {
	if len(g.Rows) == 0 {
		return errors.Errorf("grid has no rows")
	}
	ac := len(g.Rows[0].Areas)
	am := map[string]bool{}
	for _, r := range g.Rows {
		if len(r.Areas) != ac {
			return errors.Errorf("grid rows must have same size all")
		}
		for _, a := range r.Areas {
			am[a] = true
		}
	}
	for _, p := range g.Parts {
		if _, ok := am[p.Area]; !ok {
			return errors.Errorf("no definition for area %q", p.Area)
		}
	}
	return nil
}

func (g *Grid) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	for {
		token, err := d.Token()
		if err != nil {
			return err
		}
		switch t := token.(type) {
		case xml.EndElement:
			if t == start.End() {
				return g.Validate()
			}
		case xml.StartElement:
			i, err := registry.DecodeInstruction(d, t)
			if err != nil {
				continue
			}
			switch i := i.(type) {
			case *GridPart:
				g.Parts = append(g.Parts, i)
			case *GridRow:
				g.Rows = append(g.Rows, i)
			}
		}
	}
}

func (gr *GridRow) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
		case xml.CharData:
			gr.Areas = strings.Split(string(t), " ")
			if len(gr.Areas) == 0 {
				return errors.Errorf("invalid grid-row without valid areas")
			}
		}
	}
}

func (p *GridPart) DecodeAttrs(attrs []xml.Attr) error {
	for _, a := range attrs {
		switch a.Name.Local {
		case "area":
			p.Area = a.Value
		}
	}
	return p.Styled.DecodeAttrs(attrs)
}

func (p *GridPart) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				continue
			}
			p.Instructions.ISS = append(p.Instructions.ISS, i)
		}
	}
}
