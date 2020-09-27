package xdoc

import (
	"fmt"
	"strings"
)

func clearStr(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "  ", "")
	return strings.Trim(s, " ")
}

type DescribeItem struct {
	Name  string
	Value interface{}
	Items []DescribeItem
}

func (i DescribeItem) Dump(ident string) string {
	var vs string
	if i.Value != nil {
		vs = fmt.Sprintf("%v", i.Value)
	} else {
		vs = "-"
	}
	s := fmt.Sprintf("%s[%q:%v]", ident, i.Name, vs)
	sl := []string{}
	for _, sub := range i.Items {
		sl = append(sl, sub.Dump(ident+"  "))
	}
	if len(sl) > 0 {
		s += "\n" + strings.Join(sl, "\n")
	}
	return s
}

type Description struct {
	Items []DescribeItem
}

func (d *Description) Dump() string {
	sl := []string{}
	for _, sub := range d.Items {
		sl = append(sl, sub.Dump(""))
	}
	return strings.Join(sl, "\n")
}

func Describe(doc *Document) *Description {
	desc := &Description{}

	metaItem := DescribeItem{
		Name: "meta",
	}
	metaItem.Items = append(metaItem.Items, describeMeta(doc.Meta)...)
	desc.Items = append(desc.Items, metaItem)

	pageItem := DescribeItem{
		Name: "page",
	}
	pageItem.Items = append(pageItem.Items, describePage(doc.Page)...)
	desc.Items = append(desc.Items, pageItem)

	headerItem := DescribeItem{
		Name: "header",
	}
	headerItem.Items = append(headerItem.Items, describeInstructions(doc.Header)...)
	desc.Items = append(desc.Items, headerItem)

	footerItem := DescribeItem{
		Name: "footer",
	}
	footerItem.Items = append(footerItem.Items, describeInstructions(doc.Footer)...)
	desc.Items = append(desc.Items, footerItem)

	bodyItem := DescribeItem{
		Name: "body",
	}
	bodyItem.Items = append(bodyItem.Items, describeInstructions(doc.Body)...)
	desc.Items = append(desc.Items, bodyItem)

	return desc
}

func describeMeta(meta Meta) []DescribeItem {
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

func describePage(page Page) []DescribeItem {
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
			Items: describeMargins(page.Margins),
		},
	}
}

func describeMargins(margins Margins) []DescribeItem {
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

func describeTable(t *Table) []DescribeItem {
	i := DescribeItem{
		Name: "table",
	}
	for _, tr := range t.Rows {
		i.Items = append(i.Items, describeTableRow(tr)...)
	}
	return []DescribeItem{i}
}

func describeTableRow(tr *TableRow) []DescribeItem {
	i := DescribeItem{
		Name: "table-row",
	}
	for _, td := range tr.Cells {
		i.Items = append(i.Items, describeTableCell(td)...)
	}
	return []DescribeItem{i}
}

func describeTableCell(td *TableCell) []DescribeItem {
	i := DescribeItem{
		Name:  "table-cell",
		Value: clearStr(td.Content),
	}
	if len(td.Instructions) > 0 {
		i.Items = append(i.Items, describeInstructions(Instructions{
			iss: td.Instructions,
		})...)
	}
	return []DescribeItem{i}
}

func describeInstructions(iss Instructions) []DescribeItem {
	dis := []DescribeItem{}
	for _, is := range iss.iss {
		switch is := is.(type) {
		case *Font:
			dis = append(dis, DescribeItem{
				Name: "font",
			})
		case *LineFeed:
			dis = append(dis, DescribeItem{
				Name:  "line-feed",
				Value: is.Lines,
			})
		case *SetX:
			dis = append(dis, DescribeItem{
				Name:  "setx",
				Value: is.X,
			})
		case *SetY:
			dis = append(dis, DescribeItem{
				Name:  "sety",
				Value: is.Y,
			})
		case *Box:
			dis = append(dis, DescribeItem{
				Name:  "box",
				Value: clearStr(is.Text),
			})
		case *Text:
			dis = append(dis, DescribeItem{
				Name:  "text",
				Value: clearStr(is.Text),
			})
		case *Image:
			dis = append(dis, DescribeItem{
				Name:  "image",
				Value: is.Source,
			})
		case *Table:
			dis = append(dis, describeTable(is)...)
		}
	}
	return dis
}
