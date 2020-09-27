package engine

import (
	"io"

	"github.com/mazzegi/xpdf/style"
)

type Engine interface {
	Error() error
	Write(io.Writer) error
	SetPageCountAlias(alias string)
	CurrentPage() int
	OnHeader(func())
	OnFooter(func())
	AddPage()
	SetX(x float64)
	SetY(y float64)
	LineFeed(lines float64)
	ChangeFont(fnt style.Font)
}
