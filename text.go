package xpdf

import (
	"github.com/mazzegi/xpdf/markup"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/text"
)

func (p *Processor) setMarkupFont(textStyles markup.TextStyle, baseFont style.Font) {
	fnt := baseFont
	if textStyles.Italic {
		fnt.Style = style.FontStyleItalic
	}
	if textStyles.Bold {
		fnt.Weight = style.FontWeightBold
	}
	if textStyles.Mono {
		fnt.Family = "Courier"
	}
	p.engine.ChangeFont(fnt)
}

type textLine struct {
	items []*markup.TextItem
	width float64
}

func (p *Processor) textLines(items markup.Items, width float64, baseFont style.Font) []textLine {
	lines := []textLine{}
	curr := textLine{}
	for _, item := range items {
		if ci, ok := item.(*markup.ControlItem); ok && ci.Op == markup.LineFeed {
			lines = append(lines, curr)
			curr = textLine{}
			continue
		}
		textItem, ok := item.(*markup.TextItem)
		if !ok {
			continue
		}

		p.setMarkupFont(textItem.Style, baseFont)
		itemWidth := p.engine.TextWidth(" " + textItem.Text)
		if curr.width+itemWidth > width {
			lines = append(lines, curr)
			curr = textLine{}
		}

		if len(curr.items) > 0 {
			textItem.Text = " " + textItem.Text
		}
		curr.items = append(curr.items, textItem)
		curr.width += itemWidth
	}
	if len(curr.items) > 0 {
		lines = append(lines, curr)
	}
	return lines
}

func (p *Processor) writeText(s string, width float64, sty style.Styles) {
	p.engine.ChangeFont(sty.Font)
	p.engine.SetTextColor(sty.Text.R, sty.Text.G, sty.Text.B)
	s = p.tr(s)
	s = text.WhitespaceRectified(s)
	items := markup.Parse(s).Words()
	lines := p.textLines(items, width, sty.Font)
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
			p.setMarkupFont(item.Style, sty.Font)
			p.engine.WriteText(item.Text)
		}
		p.engine.LineFeed(sty.LineHeight)
	}
}
