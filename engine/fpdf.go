package engine

import (
	"io"
	"io/ioutil"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/mazzegi/xpdf/font"
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
	"github.com/pkg/errors"
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

func NewFPDF(fonts *font.Directory, doc *xdoc.Document) (*FPDF, error) {
	e := &FPDF{
		pdf: gofpdf.New(
			fpdfOrientation(doc.Page.Orientation),
			"mm",
			fpdfFormat(doc.Page.Format),
			"",
		),
	}
	err := e.initFonts(fonts)
	if err != nil {
		return nil, errors.Wrap(err, "init-fonts")
	}

	e.pdf.SetAutoPageBreak(true, doc.Page.Margins.Bottom)
	e.pdf.SetMargins(doc.Page.Margins.Left, doc.Page.Margins.Top, doc.Page.Margins.Right)

	//TODO: make code-page for unicode translator an option (per font?)
	//e.translateUnicode = e.pdf.UnicodeTranslatorFromDescriptor("")
	e.translateUnicode = func(s string) string { return s }
	return e, nil
}

func (e *FPDF) initFonts(fonts *font.Directory) error {
	return fonts.Each(func(fd font.Descriptor) error {
		bs, err := ioutil.ReadFile(fd.FilePath)
		if err != nil {
			return errors.Wrapf(err, "read file %q", fd.FilePath)
		}
		var sty string
		switch fd.Style {
		case font.Bold:
			sty = "B"
		case font.Italic:
			sty = "I"
		case font.BoldItalic:
			sty = "BI"
		default:
			sty = ""
		}
		e.pdf.AddUTF8FontFromBytes(fd.Name, sty, bs)
		if e.pdf.Error() != nil {
			return e.pdf.Error()
		}
		return nil
	})
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

func (e *FPDF) PageHeight() float64 {
	_, ph := e.pdf.GetPageSize()
	return ph
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

func (e *FPDF) Margins() (left, top, right, bottom float64) {
	return e.pdf.GetMargins()
}

func (e *FPDF) TextWidth(s string) float64 {
	return e.pdf.GetStringWidth(e.translateUnicode(s))
}

func (e *FPDF) WriteText(s string) {
	_, heightMM := e.pdf.GetFontSize()
	e.pdf.Write(heightMM, e.translateUnicode(s))
}

//drawing
func (e *FPDF) SetLineWidth(w float64) {
	e.pdf.SetLineWidth(w)
}

func (e *FPDF) SetDrawColor(r, g, b int) {
	e.pdf.SetDrawColor(r, g, b)
}

func (e *FPDF) SetFillColor(r, g, b int) {
	e.pdf.SetFillColor(r, g, b)
}

func (e *FPDF) FillRect(x, y, width, height float64) {
	e.pdf.Rect(x, y, width, height, "F")
}

func (e *FPDF) MoveTo(x, y float64) {
	e.pdf.MoveTo(x, y)
}

func (e *FPDF) LineTo(x, y float64) {
	e.pdf.LineTo(x, y)
}

func (e *FPDF) DrawPath() {
	e.pdf.DrawPath("D")
}
