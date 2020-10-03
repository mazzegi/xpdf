package xpdf

import (
	"fmt"

	"github.com/mazzegi/xpdf/style"
)

type PrintableArea struct {
	x0, y0, x1, y1 float64
}

func (pa PrintableArea) String() string {
	return fmt.Sprintf("x0=%.1f, y0=%.1f, x1=%.1f, y1=%.1f", pa.x0, pa.y0, pa.x1, pa.y1)
}

func (pa PrintableArea) Width() float64 {
	return pa.x1 - pa.x0
}

func (pa PrintableArea) Height() float64 {
	return pa.y1 - pa.y0
}

func (pa PrintableArea) WithPadding(pd style.Padding) PrintableArea {
	return PrintableArea{
		x0: pa.x0 + pd.Left,
		y0: pa.y0 + pd.Top,
		x1: pa.x1 - pd.Right,
		y1: pa.y1 - pd.Bottom,
	}
}

func (pa PrintableArea) EffectiveWidth(width float64) float64 {
	ew := pa.Width()
	if width < 0 || width > ew {
		return ew
	}
	return width
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
