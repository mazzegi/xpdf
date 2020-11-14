package xpdf

import (
	"strings"

	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/text"
	"github.com/mazzegi/xpdf/xdoc"
)

type textItem struct {
	sty  style.Styles
	text string
}

type textLine struct {
	items         []*textItem
	width         float64
	pureTextWidth float64
	paragraph     bool
}

func (p *Processor) words(s string) []string {
	return strings.Split(s, " ")
}

func (p *Processor) textLines(iss []xdoc.Instruction, width float64, sty style.Styles) []textLine {
	norm := func(s string) string {
		return text.WhitespaceRectified(p.tr(s))
	}
	lines := []textLine{}
	curr := textLine{}
	for _, is := range iss {
		var isitem *textItem
		switch is := is.(type) {
		case *xdoc.LineBreak:
			lines = append(lines, curr)
			curr = textLine{}
			continue
		case *xdoc.TextBlock:
			isitem = &textItem{
				sty:  sty,
				text: norm(is.Text),
			}
		case *xdoc.Paragraph:
			isitem = &textItem{
				sty:  is.MutatedStyles(p.doc.StyleClasses(), sty),
				text: norm(is.Text),
			}
		default:
			continue
		}
		if isitem.text == "" {
			continue
		}

		p.changeFont(isitem.sty.Font)
		words := p.words(isitem.text)
		for _, word := range words {
			item := &textItem{
				sty:  isitem.sty,
				text: word,
			}
			itemWidth := p.engine.TextWidth(" " + item.text)
			if curr.width+itemWidth >= width {
				lines = append(lines, curr)
				curr = textLine{}
			}

			if len(curr.items) > 0 {
				item.text = " " + item.text
			}
			curr.items = append(curr.items, item)
			curr.width += itemWidth
		}
	}
	if len(curr.items) > 0 {
		lines = append(lines, curr)
	}
	return lines
}

func (p *Processor) textHeight(iss []xdoc.Instruction, width float64, sty style.Styles) float64 {
	p.engine.ChangeFont(sty.Font)
	lines := p.textLines(iss, width, sty)
	lineHeight := p.engine.FontHeight() * sty.Dimension.LineSpacing
	//subtract line-spacing, to have no space below the last line
	return float64(len(lines))*lineHeight - lineHeight + p.engine.FontHeight()
}

func (p *Processor) writeText(iss []xdoc.Instruction, width float64, sty style.Styles) {
	p.engine.ChangeFont(sty.Font)
	p.engine.SetTextColor(sty.Text.Values())
	lines := p.textLines(iss, width, sty)
	xLeft, _ := p.engine.GetXY()
	for _, line := range lines {
		switch sty.HAlign {
		case style.HAlignLeft:
			p.engine.SetX(xLeft)
		case style.HAlignCenter:
			p.engine.SetX(xLeft + (width-line.width)/2.0)
		case style.HAlignRight:
			p.engine.SetX(xLeft + width - line.width)
		}
		for _, item := range line.items {
			p.changeFont(item.sty.Font)
			p.engine.WriteText(item.text)
		}
		p.engine.LineFeed(sty.LineSpacing)
	}
}
