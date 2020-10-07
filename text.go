package xpdf

import (
	"strings"

	"github.com/mazzegi/xpdf/hyphenation"
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

func (p *Processor) textLines(s string, width float64, baseFont style.Font) []textLine {
	s = p.tr(s)
	s = text.WhitespaceRectified(s)
	items := markup.Parse(s).Words()

	tryHyphenate := func(s string, currWidth float64) (s1 string, s2 string, success bool) {
		success = false
		availWidth := width - currWidth
		parts := hyphenation.Hyphenated(p.hyphenator, s)
		for i := len(parts) - 2; i >= 0; i-- {
			trial := " " + strings.Join(parts[:i+1], "") + "-"
			trialWidth := p.engine.TextWidth(trial)
			if trialWidth <= availWidth {
				s1 = trial
				s2 = strings.Join(parts[i+1:], "")
				success = true
				return
			}
		}
		return
	}

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
		if curr.width+itemWidth >= width {
			//try hyphenation
			s1, s2, ok := tryHyphenate(textItem.Text, curr.width)
			if !ok {
				lines = append(lines, curr)
				curr = textLine{}
			} else {
				curr.items = append(curr.items, &markup.TextItem{
					Text:  s1,
					Style: textItem.Style,
				})
				curr.width += p.engine.TextWidth(s1)
				lines = append(lines, curr)

				curr = textLine{}
				textItem.Text = s2
				itemWidth = p.engine.TextWidth(s2)
			}
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

func (p *Processor) textHeight(s string, width float64, sty style.Styles) float64 {
	p.engine.ChangeFont(sty.Font)
	lines := p.textLines(s, width, sty.Font)
	lineHeight := p.engine.FontHeight() * sty.Dimension.LineSpacing
	//subtract line-spacing, to have no space below the last line
	return float64(len(lines))*lineHeight - lineHeight + p.engine.FontHeight()
}

func (p *Processor) writeText(s string, width float64, sty style.Styles) {
	p.engine.ChangeFont(sty.Font)
	p.engine.SetTextColor(sty.Text.Values())
	lines := p.textLines(s, width, sty.Font)
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
		p.engine.LineFeed(sty.LineSpacing)
	}
}
