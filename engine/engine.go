package engine

import (
	"io"

	"github.com/mazzegi/xpdf/style"
)

type Engine interface {
	Error() error
	WritePDF(io.Writer) error
	SetPageCountAlias(alias string)
	CurrentPage() int
	OnHeader(func())
	OnFooter(func())
	AddPage()
	SetX(x float64)
	SetY(y float64)
	GetXY() (float64, float64)
	LineFeed(lines float64)
	ChangeFont(fnt style.Font)
	EffectiveWidth(width float64) float64
	PutImage(src string, x, y, width, height float64)
	SetTextColor(r, g, b int)
	FontHeight() float64
	TextWidth(s string) float64
	WriteText(s string)
}
