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
	return p.engine.Write(w)
}

func (p *Processor) tr(s string) string {
	return strings.ReplaceAll(s, "{cp}", fmt.Sprintf("%d", p.engine.CurrentPage()))
}

func (p *Processor) changeFont(fnt style.Font) {
	p.currStyles.Font = fnt
	p.engine.ChangeFont(p.currStyles.Font)
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
			//p.renderTextBox(i.Text, p.appliedStyles(i))
		case *xdoc.Text:
			//p.renderText(i, p.appliedStyles(i))
		case *xdoc.Table:
			//p.renderTable(i, p.appliedStyles(i))
		case *xdoc.Image:
			//p.renderImage(i, p.appliedStyles(i))
		}
	}
}
