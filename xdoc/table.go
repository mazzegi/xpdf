package xdoc

import "encoding/xml"

type Table struct {
	Styled
	XMLName xml.Name    `xml:"table"`
	Rows    []*TableRow `xml:"tr"`
}

type TableRow struct {
	Styled
	XMLName xml.Name     `xml:"tr"`
	Cells   []*TableCell `xml:"td"`
}

type TableCell struct {
	Styled
	XMLName      xml.Name `xml:"td"`
	Content      string   `xml:",chardata"`
	Instructions []Instruction
}

func (tab *Table) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				return err
			}
			switch i := i.(type) {
			case *TableRow:
				tab.Rows = append(tab.Rows, i)
			}
		}
	}
}

func (row *TableRow) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				return err
			}
			switch i := i.(type) {
			case *TableCell:
				row.Cells = append(row.Cells, i)
			}
		}
	}
}

func (cell *TableCell) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
				continue
			}
			cell.Instructions = append(cell.Instructions, i)
		case xml.CharData:
			cell.Content += string(t)
		}
	}
}
