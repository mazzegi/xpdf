package xpdf

import (
	"strings"

	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/text"
	"github.com/mazzegi/xpdf/xdoc"
)

func (p *Processor) textLinesHyphenated(iss []xdoc.Instruction, width float64, sty style.Styles) []textLine {
	norm := func(s string) string {
		return text.WhitespaceRectified(p.tr(s))
	}
	tryHyphenate := func(s string, currWidth float64) (s1 string, s2 string, success bool) {
		success = false
		availWidth := width - currWidth
		parts := p.hyphenator.Hyphenate(s)
		if len(parts) < 2 {
			return
		}
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
	for i, is := range iss {
		Logf("generate lines: %d ", i)
		var isitem *textItem
		switch is := is.(type) {
		case *xdoc.LineBreak:
			if len(lines) > 0 {
				lines[len(lines)-1].paragraph = true
			}

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

		p.changeFont(isitem.sty.Font)
		words := p.words(isitem.text)
		for _, word := range words {
			item := &textItem{
				sty:  isitem.sty,
				text: word,
			}
			itemWidth := p.engine.TextWidth(" " + item.text)
			pureItemWidth := p.engine.TextWidth(item.text)
			if curr.width+itemWidth >= width {
				//try hyphenation
				s1, s2, ok := tryHyphenate(item.text, curr.width)
				if !ok {
					lines = append(lines, curr)
					curr = textLine{}
				} else {
					curr.items = append(curr.items, &textItem{
						text: s1,
						sty:  item.sty,
					})
					curr.width += p.engine.TextWidth(s1)
					curr.pureTextWidth += p.engine.TextWidth(strings.Trim(s1, " "))
					lines = append(lines, curr)

					curr = textLine{}
					item.text = s2
					itemWidth = p.engine.TextWidth(s2)
					pureItemWidth = p.engine.TextWidth(strings.Trim(s2, " "))
				}
			}

			if len(curr.items) > 0 {
				item.text = " " + item.text
			}
			curr.items = append(curr.items, item)
			curr.width += itemWidth
			curr.pureTextWidth += pureItemWidth
		}
	}
	if len(curr.items) > 0 {
		curr.paragraph = true
		lines = append(lines, curr)
	}
	return lines
}

func (p *Processor) textHeightHyphenated(iss []xdoc.Instruction, width float64, sty style.Styles) float64 {
	p.engine.ChangeFont(sty.Font)
	lines := p.textLinesHyphenated(iss, width, sty)
	lineHeight := p.engine.FontHeight() * sty.Dimension.LineSpacing
	//subtract line-spacing, to have no space below the last line
	return float64(len(lines))*lineHeight - lineHeight + p.engine.FontHeight()
}

func (p *Processor) writeTextHyphenated(iss []xdoc.Instruction, width float64, sty style.Styles) {
	Logf("write text hyphenated: %d instr", len(iss))
	p.engine.ChangeFont(sty.Font)
	p.engine.SetTextColor(sty.Text.Values())
	lines := p.textLinesHyphenated(iss, width, sty)
	Logf("write text hyphenated: write %d lines", len(lines))
	xLeft, _ := p.engine.GetXY()
	for _, line := range lines {
		p.engine.SetX(xLeft)
		spaceCnt := len(line.items) - 1
		if spaceCnt < 1 || line.paragraph {
			for _, item := range line.items {
				p.changeFont(item.sty.Font)
				p.engine.WriteText(item.text)
			}
		} else {
			//subtract another 0.1 to avoid page breaks on equal widths
			spaceWidth := (width - 0.1 - line.pureTextWidth) / float64(spaceCnt)
			for _, item := range line.items {
				p.changeFont(item.sty.Font)
				p.engine.WriteText(strings.Trim(item.text, " "))
				cx, _ := p.engine.GetXY()
				p.engine.SetX(cx + spaceWidth)
			}
		}
		p.engine.LineFeed(sty.LineSpacing)
	}
}
