package xpdf

type PrintableArea struct {
	x0, y0, x1, y1 float64
}

func (pa PrintableArea) Width() float64 {
	return pa.x1 - pa.x0
}

func (pa PrintableArea) Height() float64 {
	return pa.y1 - pa.y0
}

type Page struct {
	width         float64
	height        float64
	printableArea PrintableArea
}

func (p Page) EffectiveWidth(width float64) float64 {
	ew := p.printableArea.x1 - p.printableArea.x0
	if width < 0 || width > ew {
		return ew
	}
	return width
}
