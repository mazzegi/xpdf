package xpdf

import (
	"fmt"
	"io"
	"strings"

	"github.com/mazzegi/xpdf/engine"
	"github.com/mazzegi/xpdf/hyphenation"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

type Processor struct {
	engine     engine.Engine
	doc        *xdoc.Document
	currStyles style.Styles
	hyphenator *hyphenation.PatternLookup
}

func NewProcessor(engine engine.Engine, hyphenator *hyphenation.PatternLookup, doc *xdoc.Document) *Processor {
	p := &Processor{
		engine:     engine,
		hyphenator: hyphenator,
		doc:        doc,
		currStyles: DefaultStyle(),
	}
	return p
}

func (p *Processor) Process(w io.Writer) error {
	//TODO: make page-count and current-page aliases options
	p.engine.SetPageCountAlias("{np}")
	p.engine.OnHeader(func() {
		p.processInstructions(p.doc.Header)
	})
	p.engine.OnFooter(func() {
		p.processInstructions(p.doc.Footer)
	})

	//Change font to initial default font
	p.changeFont(p.currStyles.Font)

	p.engine.AddPage()
	p.processInstructions(p.doc.Body)

	err := p.engine.Error()
	if err != nil {
		return err
	}
	return p.engine.WritePDF(w)
}

func (p *Processor) tr(s string) string {
	return strings.ReplaceAll(s, "{cp}", fmt.Sprintf("%d", p.engine.CurrentPage()))
}

func (p *Processor) changeFont(fnt style.Font) {
	p.currStyles.Font = fnt
	p.engine.ChangeFont(p.currStyles.Font)
}

func (p *Processor) page() Page {
	x0, y0, x1, y1 := p.engine.PrintableArea()
	page := Page{
		width:  p.engine.PageWidth(),
		height: p.engine.PageHeight(),
		printableArea: PrintableArea{
			x0: x0,
			y0: y0,
			x1: x1,
			y1: y1,
		},
	}
	return page
}

func (p *Processor) resetStyles() {
	p.engine.ChangeFont(p.currStyles.Font)
	p.engine.SetTextColor(p.currStyles.Text.R, p.currStyles.Text.G, p.currStyles.Text.B)
}

func (p *Processor) processInstructions(is xdoc.Instructions) {
	for _, i := range is.ISS {
		switch i := i.(type) {
		case *xdoc.Font:
			p.changeFont(i.MutatedStyles(p.doc.StyleClasses(), p.currStyles).Font)
		case *xdoc.LineFeed:
			p.engine.LineFeed(i.Lines)
		case *xdoc.SetX:
			p.engine.SetX(i.X)
		case *xdoc.SetY:
			p.engine.SetY(i.Y)
		case *xdoc.Box:
			p.renderTextBox(i, p.page().printableArea)
		case *xdoc.Text:
			p.renderText(i)
		case *xdoc.Table:
			p.renderTable(i)
		case *xdoc.Image:
			p.renderImage(i, p.page().printableArea)
		}
	}
}

func (p *Processor) renderText(text *xdoc.Text) {
	if text.Text == "" {
		return
	}
	defer p.resetStyles()
	sty := text.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	width := p.page().EffectiveWidth(sty.Width)

	p.writeTextFnc(sty)(text.Text, width, sty)
}

func (p *Processor) textBoxHeight(box *xdoc.Box, pa PrintableArea) float64 {
	defer p.resetStyles()
	sty := box.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	width := pa.EffectiveWidth(sty.Width) - sty.Padding.Left - sty.Padding.Right
	var height float64
	if sty.Dimension.Height <= 0 {
		if box.Text == "" {
			height = p.engine.FontHeight()
		} else {
			height = p.textHeightFnc(sty)(box.Text, width, sty)
		}
	} else {
		height = sty.Dimension.Height
	}
	return height
}

func (p *Processor) renderTextBox(box *xdoc.Box, pa PrintableArea) {
	defer p.resetStyles()
	sty := box.MutatedStyles(p.doc.StyleClasses(), p.currStyles)

	width := pa.EffectiveWidth(sty.Width) - sty.Padding.Left - sty.Padding.Right
	var height float64
	if sty.Dimension.Height <= 0 {
		if box.Text != "" {
			height = p.textHeightFnc(sty)(box.Text, width, sty)
		} else {
			height = p.engine.FontHeight()
		}
	} else {
		height = sty.Dimension.Height
	}

	x0, y0 := p.engine.GetXY()
	_, _, _, bm := p.engine.Margins()
	effHeight := p.engine.PageHeight() - bm
	if y0+height >= effHeight {
		p.engine.AddPage()
		x0, y0 = p.engine.GetXY()
	}

	x0 += sty.Dimension.OffsetX
	y0 += sty.Dimension.OffsetY
	y1 := y0 + height + sty.Box.Padding.Top + sty.Box.Padding.Bottom
	x1 := x0 + width + sty.Padding.Left + sty.Padding.Right
	p.drawBox(x0, y0, x1, y1, sty)

	p.engine.SetY(y0 + sty.Box.Padding.Top)
	p.engine.SetX(x0 + sty.Box.Padding.Left)
	if box.Text != "" {
		p.writeTextFnc(sty)(box.Text, width, sty)
	}
	p.engine.SetY(y1)
}

func (p *Processor) renderImage(img *xdoc.Image, pa PrintableArea) {
	sty := img.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	x, y := p.engine.GetXY()
	x += sty.Dimension.OffsetX
	y += sty.Dimension.OffsetY
	width, height := sty.Width, sty.Height
	if width <= 0 || height <= 0 || width > pa.Width() || height > pa.Height() {
		width, height = pa.Width(), pa.Height()
	}
	p.engine.PutImage(img.Source, x, y, width, height)
}

//
func (p *Processor) textHeightFnc(sty style.Styles) func(string, float64, style.Styles) float64 {
	switch sty.HAlign {
	case style.HAlignBlock:
		return p.textHeightHyphenated
	default:
		return p.textHeight
	}
}

func (p *Processor) writeTextFnc(sty style.Styles) func(string, float64, style.Styles) {
	switch sty.HAlign {
	case style.HAlignBlock:
		return p.writeTextHyphenated
	default:
		return p.writeText
	}
}
