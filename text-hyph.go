package xpdf

import (
	"strings"

	"github.com/mazzegi/xpdf/markup"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/text"
	"github.com/mazzegi/xpdf/xdoc"
)

func (p *Processor) textLinesHyphenated(s string, width float64, baseFont style.Font) []textLine {
	s = p.tr(s)
	s = text.WhitespaceRectified(s)
	items := markup.Parse(s).Words()

	tryHyphenate := func(s string, currWidth float64) (s1 string, s2 string, success bool) {
		success = false
		availWidth := width - currWidth
		parts := p.hyphenator.Hyphenate(s)
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
			curr.paragraph = true
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
		pureItemWidth := p.engine.TextWidth(textItem.Text)
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
				curr.pureTextWidth += p.engine.TextWidth(strings.Trim(s1, " "))
				lines = append(lines, curr)

				curr = textLine{}
				textItem.Text = s2
				itemWidth = p.engine.TextWidth(s2)
				pureItemWidth = p.engine.TextWidth(strings.Trim(s2, " "))
			}
		}

		if len(curr.items) > 0 {
			textItem.Text = " " + textItem.Text
		}
		curr.items = append(curr.items, textItem)
		curr.width += itemWidth
		curr.pureTextWidth += pureItemWidth
	}
	if len(curr.items) > 0 {
		lines = append(lines, curr)
	}
	return lines
}

func (p *Processor) textHeightHyphenated(iss []xdoc.Instruction, width float64, sty style.Styles) float64 {
	p.engine.ChangeFont(sty.Font)
	lines := p.textLinesHyphenated(s, width, sty.Font)
	lineHeight := p.engine.FontHeight() * sty.Dimension.LineSpacing
	//subtract line-spacing, to have no space below the last line
	return float64(len(lines))*lineHeight - lineHeight + p.engine.FontHeight()
}

func (p *Processor) writeTextHyphenated(iss []xdoc.Instruction, width float64, sty style.Styles) {
	p.engine.ChangeFont(sty.Font)
	p.engine.SetTextColor(sty.Text.Values())
	lines := p.textLinesHyphenated(s, width, sty.Font)
	xLeft, _ := p.engine.GetXY()
	for _, line := range lines {
		p.engine.SetX(xLeft)
		spaceCnt := len(line.items) - 1
		if spaceCnt < 1 || line.paragraph {
			for _, item := range line.items {
				p.setMarkupFont(item.Style, sty.Font)
				p.engine.WriteText(item.Text)
			}
		} else {
			//subtract another 0.1 to avoid page breaks on equal widths
			spaceWidth := (width - 0.1 - line.pureTextWidth) / float64(spaceCnt)
			for _, item := range line.items {
				p.setMarkupFont(item.Style, sty.Font)
				//p.engine.WriteText(item.Text)
				p.engine.WriteText(strings.Trim(item.Text, " "))
				cx, _ := p.engine.GetXY()
				p.engine.SetX(cx + spaceWidth)
			}
		}
		p.engine.LineFeed(sty.LineSpacing)
	}
}
