package xpdf

import (
	"fmt"
	"io"
	"strings"

	"github.com/mazzegi/xpdf/engine"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

type Processor struct {
	engine     engine.Engine
	doc        *xdoc.Document
	currStyles style.Styles
}

func NewProcessor(engine engine.Engine, doc *xdoc.Document) *Processor {
	p := &Processor{
		engine:     engine,
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
			p.renderTextBox(i)
		case *xdoc.Text:
			p.renderText(i)
		case *xdoc.Table:
			p.renderTable(i)
		case *xdoc.Image:
			p.renderImage(i)
		}
	}
}

func (p *Processor) renderText(text *xdoc.Text) {
	defer p.resetStyles()
	sty := text.MutatedStyles(p.doc.StyleClasses(), p.currStyles)

	width := p.engine.EffectiveWidth(sty.Dimension.Width)
	p.writeText(text.Text, width, sty)
}

func (p *Processor) renderTextBox(box *xdoc.Box) {
	defer p.resetStyles()
	sty := box.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	_ = sty
}

func (p *Processor) renderTable(table *xdoc.Table) {
	defer p.resetStyles()
	sty := table.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	_ = sty
}

func (p *Processor) renderImage(img *xdoc.Image) {
	sty := img.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	x, y := p.engine.GetXY()
	x += sty.Dimension.OffsetX
	y += sty.Dimension.OffsetY
	p.engine.PutImage(img.Source, x, y, sty.Dimension.Width, sty.Dimension.Height)
}
