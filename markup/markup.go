package markup

import (
	"fmt"
	"strings"
)

const asterisk = "*"
const asterisks = "**"
const underscore = "_"
const underscores = "__"
const backtick = "`"
const backslash = `\`

type ControlOp int

const (
	LineFeed ControlOp = 0x01
)

func (c ControlOp) String() string {
	switch c {
	case LineFeed:
		return "linefeed"
	default:
		return "invalid op"
	}
}

type Item interface {
	String() string
}

type TextStyle struct {
	Italic bool
	Bold   bool
	Mono   bool
}

type TextItem struct {
	Text  string
	Style TextStyle
}

func (i *TextItem) append(bs ...byte) {
	cs := []byte(i.Text)
	cs = append(cs, bs...)
	i.Text = string(cs)
}

func (i TextItem) String() string {
	sl := []string{}
	if i.Style.Italic {
		sl = append(sl, "italic")
	}
	if i.Style.Bold {
		sl = append(sl, "bold")
	}
	if i.Style.Mono {
		sl = append(sl, "mono")
	}
	return fmt.Sprintf("text: %q (%s)", i.Text, strings.Join(sl, ", "))
}

func (i TextItem) Words() Items {
	is := Items{}
	words := []string{}
	currWord := ""
	for _, r := range i.Text {
		if r == ' ' {
			words = append(words, currWord)
			currWord = ""
		} else {
			currWord += string(r)
		}
	}
	if currWord != "" {
		words = append(words, currWord)
	}
	for _, word := range words {
		is = append(is, TextItem{
			Text:  word,
			Style: i.Style,
		})
	}
	return is
}

type ControlItem struct {
	Op ControlOp
}

func (i ControlItem) String() string {
	return fmt.Sprintf("control: %q ", i.Op.String())
}

type Items []Item

func (is Items) Words() Items {
	wis := Items{}
	for _, i := range is {
		switch i := i.(type) {
		case TextItem:
			wis = append(wis, i.Words()...)
		default:
			wis = append(wis, i)
		}
	}
	return wis
}

func Parse(s string) Items {
	items := Items{}
	if s == "" {
		return items
	}
	bs := []byte(s)
	currStyle := TextStyle{}
	currTextItem := &TextItem{
		Text:  "",
		Style: currStyle,
	}
	flush := func() {
		if currTextItem.Text != "" {
			items = append(items, currTextItem)
		}
		currTextItem = &TextItem{
			Text:  "",
			Style: currStyle,
		}
	}

	i := 0
	for {
		if i >= len(s) {
			flush()
			return items
		}
		b := bs[i]
		if strings.HasPrefix(s[i:], asterisks) || strings.HasPrefix(s[i:], underscores) {
			currStyle.Bold = !currStyle.Bold
			flush()
			i += 2
		} else if strings.HasPrefix(s[i:], asterisk) || strings.HasPrefix(s[i:], underscore) {
			currStyle.Italic = !currStyle.Italic
			flush()
			i += 1
		} else if strings.HasPrefix(s[i:], backtick) {
			currStyle.Mono = !currStyle.Mono
			flush()
			i += 1
		} else if strings.HasPrefix(s[i:], backslash) {
			flush()
			items = append(items, &ControlItem{
				Op: LineFeed,
			})
			i += 1
		} else {
			currTextItem.append(b)
			i += 1
		}
	}
}
