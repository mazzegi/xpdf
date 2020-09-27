package xdoc

import (
	"encoding/json"
	"encoding/xml"
	"strings"

	"github.com/mazzegi/xpdf/style"
)

func clearStr(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "  ", "")
	return strings.Trim(s, " ")
}

type DescribeItem struct {
	Name       string
	Value      interface{}
	StyleDiffs []Diff
	Items      []DescribeItem
}

type Description struct {
	Items []DescribeItem
	doc   *Document
}

type DescribeFormat string

const (
	DescribeXML  DescribeFormat = "xml"
	DescribeJSON DescribeFormat = "json"
)

func (d *Description) Dump(format DescribeFormat) string {
	var bs []byte
	switch format {
	case DescribeXML:
		bs, _ = xml.MarshalIndent(d, "", "  ")
	case DescribeJSON:
		bs, _ = json.MarshalIndent(d, "", "  ")
	default:
		bs, _ = json.MarshalIndent(d, "", "  ")
	}
	return string(bs)
}

func Describe(doc *Document) *Description {
	desc := &Description{
		doc: doc,
	}

	metaItem := DescribeItem{
		Name: "meta",
	}
	metaItem.Items = append(metaItem.Items, desc.describeMeta(doc.Meta)...)
	desc.Items = append(desc.Items, metaItem)

	pageItem := DescribeItem{
		Name: "page",
	}
	pageItem.Items = append(pageItem.Items, desc.describePage(doc.Page)...)
	desc.Items = append(desc.Items, pageItem)

	headerItem := DescribeItem{
		Name: "header",
	}
	headerItem.Items = append(headerItem.Items, desc.describeInstructions(doc.Header)...)
	desc.Items = append(desc.Items, headerItem)

	footerItem := DescribeItem{
		Name: "footer",
	}
	footerItem.Items = append(footerItem.Items, desc.describeInstructions(doc.Footer)...)
	desc.Items = append(desc.Items, footerItem)

	bodyItem := DescribeItem{
		Name: "body",
	}
	bodyItem.Items = append(bodyItem.Items, desc.describeInstructions(doc.Body)...)
	desc.Items = append(desc.Items, bodyItem)

	return desc
}

func (desc *Description) describeMeta(meta Meta) []DescribeItem {
	return []DescribeItem{
		{
			Name:  "author",
			Value: meta.Author,
		},
		{
			Name:  "creator",
			Value: meta.Creator,
		},
		{
			Name:  "subject",
			Value: meta.Subject,
		},
	}
}

func (desc *Description) describePage(page Page) []DescribeItem {
	return []DescribeItem{
		{
			Name:  "orientation",
			Value: page.Orientation,
		},
		{
			Name:  "format",
			Value: page.Format,
		},
		{
			Name:  "margins",
			Items: desc.describeMargins(page.Margins),
		},
	}
}

func (desc *Description) describeMargins(margins Margins) []DescribeItem {
	return []DescribeItem{
		{
			Name:  "left",
			Value: margins.Right,
		},
		{
			Name:  "top",
			Value: margins.Top,
		},
		{
			Name:  "right",
			Value: margins.Right,
		},
		{
			Name:  "bottom",
			Value: margins.Bottom,
		},
	}
}

func (desc *Description) describeTable(t *Table) []DescribeItem {
	i := DescribeItem{
		Name:       "table",
		StyleDiffs: desc.describeMutator(t),
	}
	for _, tr := range t.Rows {
		i.Items = append(i.Items, desc.describeTableRow(tr)...)
	}
	return []DescribeItem{i}
}

func (desc *Description) describeTableRow(tr *TableRow) []DescribeItem {
	i := DescribeItem{
		Name:       "table-row",
		StyleDiffs: desc.describeMutator(tr),
	}
	for _, td := range tr.Cells {
		i.Items = append(i.Items, desc.describeTableCell(td)...)
	}
	return []DescribeItem{i}
}

func (desc *Description) describeTableCell(td *TableCell) []DescribeItem {
	i := DescribeItem{
		Name:       "table-cell",
		StyleDiffs: desc.describeMutator(td),
		Value:      clearStr(td.Content),
	}
	if len(td.Instructions) > 0 {
		i.Items = append(i.Items, desc.describeInstructions(Instructions{
			ISS: td.Instructions,
		})...)
	}
	return []DescribeItem{i}
}

func (desc *Description) describeInstructions(iss Instructions) []DescribeItem {
	dis := []DescribeItem{}
	for _, is := range iss.ISS {
		switch is := is.(type) {
		case *Font:
			dis = append(dis, DescribeItem{
				Name:       "font",
				StyleDiffs: desc.describeMutator(is),
			})
		case *LineFeed:
			dis = append(dis, DescribeItem{
				Name:       "line-feed",
				Value:      is.Lines,
				StyleDiffs: desc.describeMutator(is),
			})
		case *SetX:
			dis = append(dis, DescribeItem{
				Name:       "setx",
				Value:      is.X,
				StyleDiffs: desc.describeMutator(is),
			})
		case *SetY:
			dis = append(dis, DescribeItem{
				Name:       "sety",
				Value:      is.Y,
				StyleDiffs: desc.describeMutator(is),
			})
		case *Box:
			dis = append(dis, DescribeItem{
				Name:       "box",
				Value:      clearStr(is.Text),
				StyleDiffs: desc.describeMutator(is),
			})
		case *Text:
			dis = append(dis, DescribeItem{
				Name:       "text",
				Value:      clearStr(is.Text),
				StyleDiffs: desc.describeMutator(is),
			})
		case *Image:
			dis = append(dis, DescribeItem{
				Name:       "image",
				Value:      is.Source,
				StyleDiffs: desc.describeMutator(is),
			})
		case *Table:
			dis = append(dis, desc.describeTable(is)...)
		}
	}
	return dis
}

func (desc *Description) describeMutator(ins Instruction) []Diff {
	sty := style.Styles{}
	mutsty := ins.MutatedStyles(desc.doc.styleClasses, sty)
	diffs, err := stylesDiff(sty, mutsty)
	if err != nil {
		return []Diff{
			{
				Path: "ERROR",
				Mod:  err.Error(),
			},
		}
	}
	return diffs
}
