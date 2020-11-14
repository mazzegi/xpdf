package xpdf

import "github.com/mazzegi/xpdf/style"

func (p *Processor) drawBox(x0, y0, x1, y1 float64, sty style.Styles) {
	p.engine.SetLineWidth(sty.Draw.LineWidth)
	p.engine.SetDrawColor(sty.Color.Foreground.Values())
	p.engine.SetFillColor(sty.Color.Background.Values())

	halfLine := sty.Draw.LineWidth / 2
	width := x1 - x0
	height := y1 - y0

	p.engine.FillRect(x0+halfLine, y0+halfLine, width-2*halfLine, height-2*halfLine)
	//p.engine.FillRect(x0, y0, width, height)

	p.engine.MoveTo(x0-halfLine, y0)
	if sty.Box.Border.Top > 0 {
		p.engine.LineTo(x0+width, y0)
	} else {
		p.engine.MoveTo(x0+width, y0)
	}
	if sty.Box.Border.Right > 0 {
		p.engine.LineTo(x0+width, y1)
	} else {
		p.engine.MoveTo(x0+width, y1)
	}
	if sty.Box.Border.Bottom > 0 {
		p.engine.LineTo(x0, y1)
	} else {
		p.engine.MoveTo(x0, y1)
	}
	if sty.Box.Border.Left > 0 {
		p.engine.LineTo(x0, y0)
	} else {
		p.engine.MoveTo(x0, y0)
	}
	p.engine.DrawPath()
}
