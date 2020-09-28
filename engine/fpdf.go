package engine

import (
	"io"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

func fpdfOrientation(o xdoc.Orientation) string {
	switch o {
	case xdoc.OrientationPortrait:
		return "P"
	case xdoc.OrientationLandscape:
		return "L"
	}
	return ""
}

func fpdfFormat(f xdoc.Format) string {
	switch f {
	case xdoc.FormatA3:
		return "A3"
	case xdoc.FormatA4:
		return "A4"
	case xdoc.FormatA5:
		return "A5"
	case xdoc.FormatLetter:
		return "Letter"
	case xdoc.FormatLegal:
		return "Legal"
	default:
		return "A4"
	}
}

func fpdfFontStyle(fnt style.Font) string {
	s := ""
	switch fnt.Style {
	case style.FontStyleItalic:
		s += "I"
	}
	switch fnt.Weight {
	case style.FontWeightBold:
		s += "B"
	}
	switch fnt.Decoration {
	case style.FontDecorationUnderline:
		s += "U"
	}
	return s
}

type FPDF struct {
	pdf              *gofpdf.Fpdf
	translateUnicode func(s string) string
}

func NewFPDF(doc *xdoc.Document) *FPDF {
	e := &FPDF{
		pdf: gofpdf.New(
			fpdfOrientation(doc.Page.Orientation),
			"mm",
			fpdfFormat(doc.Page.Format),
			"",
		),
	}
	e.pdf.SetAutoPageBreak(true, doc.Page.Margins.Bottom)
	e.pdf.SetMargins(doc.Page.Margins.Left, doc.Page.Margins.Top, doc.Page.Margins.Right)

	//TODO: make code-page for unicode translator an option
	e.translateUnicode = e.pdf.UnicodeTranslatorFromDescriptor("")
	return e
}

func (e *FPDF) Error() error {
	return e.pdf.Error()
}

func (e *FPDF) WritePDF(w io.Writer) error {
	return e.pdf.Output(w)
}

func (e *FPDF) SetPageCountAlias(alias string) {
	e.pdf.AliasNbPages(alias)
}

func (e *FPDF) CurrentPage() int {
	return e.pdf.PageNo()
}

func (e *FPDF) OnHeader(f func()) {
	e.pdf.SetHeaderFunc(f)
}

func (e *FPDF) OnFooter(f func()) {
	e.pdf.SetFooterFunc(f)
}

func (e *FPDF) AddPage() {
	e.pdf.AddPage()
}

func (e *FPDF) SetX(x float64) {
	e.pdf.SetX(x)
}

func (e *FPDF) SetY(y float64) {
	e.pdf.SetY(y)
}

func (e *FPDF) GetXY() (float64, float64) {
	return e.pdf.GetXY()
}

func (e *FPDF) LineFeed(lines float64) {
	_, heightMM := e.pdf.GetFontSize()
	e.pdf.Ln(heightMM * lines)
}

func (e *FPDF) ChangeFont(fnt style.Font) {
	e.pdf.SetFont(string(fnt.Family), fpdfFontStyle(fnt), float64(fnt.PointSize))
}

func (e *FPDF) EffectiveWidth(width float64) float64 {
	l, _, r, _ := e.pdf.GetMargins()
	pw, _ := e.pdf.GetPageSize()
	ew := pw - (l + r) - 3 // without subtracting 3 it doesn't fit
	if width < 0 || width > ew {
		return ew
	}
	return width
}

func (e *FPDF) PutImage(src string, x, y, width, height float64) {
	e.pdf.ImageOptions(src, x, y, width, height, false, gofpdf.ImageOptions{}, 0, "")
}

func (e *FPDF) SetTextColor(r, g, b int) {
	e.pdf.SetTextColor(r, g, b)
}

func (e *FPDF) FontHeight() float64 {
	_, heightMM := e.pdf.GetFontSize()
	return heightMM
}

func (e *FPDF) TextWidth(s string) float64 {
	return e.pdf.GetStringWidth(s)
}

func (e *FPDF) WriteText(s string) {
	// x, y := e.pdf.GetXY()
	// e.pdf.Text(x, y, s)
	_, heightMM := e.pdf.GetFontSize()
	e.pdf.Write(heightMM, s)
}
