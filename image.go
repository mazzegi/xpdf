package xpdf

import (
	"image"
	"os"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/mazzegi/xpdf/xdoc"
	"github.com/pkg/errors"
)

const Dpi96 = 96

type ImageDescriptor struct {
	Format   string
	WidthPx  int
	HeightPx int
}

func (id ImageDescriptor) WidthMm(dpi int) float64 {
	pxPerMm := float64(dpi) / 25.4
	return float64(id.WidthPx) / pxPerMm
}

func (id ImageDescriptor) HeightMm(dpi int) float64 {
	pxPerMm := float64(dpi) / 25.4
	return float64(id.HeightPx) / pxPerMm
}

func DescribeImage(src string) (ImageDescriptor, error) {
	f, err := os.Open(src)
	if err != nil {
		return ImageDescriptor{}, errors.Wrapf(err, "open image %q", src)
	}
	defer f.Close()
	cfg, format, err := image.DecodeConfig(f)
	if err != nil {
		return ImageDescriptor{}, errors.Wrapf(err, "decode config of %q", src)
	}
	return ImageDescriptor{
		Format:   format,
		WidthPx:  cfg.Width,
		HeightPx: cfg.Height,
	}, nil
}

func (p *Processor) renderImage(img *xdoc.Image, pa PrintableArea) {
	imgSrc := p.resolveFile(img.Source)

	iDesc, err := DescribeImage(imgSrc)
	if err != nil {
		Logf("ERROR: describe image: %v", err)
		return
	}
	sty := img.MutatedStyles(p.doc.StyleClasses(), p.currStyles)
	x, y := p.engine.GetXY()
	x += sty.OffsetX
	y += sty.OffsetY
	paWidth := pa.Width() - sty.OffsetX
	paHeight := pa.Height() - sty.OffsetY

	idWidth := iDesc.WidthMm(Dpi96)
	idHeight := iDesc.HeightMm(Dpi96)

	var width, height float64
	switch {
	case sty.Width > 0 && sty.Height > 0:
		width, height = sty.Width, sty.Height
	case sty.Width <= 0 && sty.Height > 0:
		height = sty.Height
		width = float64(idWidth) / float64(idHeight) * sty.Height
	case sty.Width > 0 && sty.Height <= 0:
		width = sty.Width
		height = float64(idHeight) / float64(idWidth) * sty.Width
	default:
		width, height = float64(idWidth), float64(idHeight)
	}

	//scale to printable area
	//TODO: probably width and height are in millimeters here !!!
	if width > paWidth {
		height = height * paWidth / width
		width = paWidth
	}
	if height > paHeight {
		width = width * paHeight / height
		height = paHeight
	}

	p.engine.PutImage(imgSrc, x, y, width, height)
}
