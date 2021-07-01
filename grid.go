package xpdf

import (
	"github.com/mazzegi/log"
	"github.com/mazzegi/xpdf/xdoc"
	"github.com/pkg/errors"
)

type areaRow struct {
	idx  int
	cols []int
}

type gridBox struct {
	area                     string
	rows                     []*areaRow
	left, top, right, bottom int
	pa                       PrintableArea
}

func (b *gridBox) add(row, col int) error {
	for _, r := range b.rows {
		if r.idx == row {
			if col != r.cols[len(r.cols)-1]+1 {
				return errors.Errorf("not sequential cols")
			}
			r.cols = append(r.cols, col)
			if col > b.right {
				b.right = col
			}
			return nil
		}
	}
	if len(b.rows) > 0 && row != b.rows[len(b.rows)-1].idx+1 {
		return errors.Errorf("not sequential rows")
	}
	if len(b.rows) == 0 {
		b.left, b.right = col, col
		b.top, b.bottom = row, row
	} else {
		if col != b.left {
			return errors.Errorf("new row must have same first col index")
		}
		b.bottom = row
	}
	b.rows = append(b.rows, &areaRow{
		idx:  row,
		cols: []int{col},
	})
	return nil
}

func (b *gridBox) validate() error {
	rowsize := len(b.rows[0].cols)
	for _, r := range b.rows {
		if len(r.cols) != rowsize {
			return errors.Errorf("non equal rowsize for area %q", b.area)
		}
	}
	return nil
}

func (p *Processor) renderGrid(g *xdoc.Grid, pa PrintableArea) {
	boxes := map[string]*gridBox{}
	for ir, r := range g.Rows {
		for ia, a := range r.Areas {
			b, ok := boxes[a]
			if !ok {
				b = &gridBox{
					area: a,
				}
				boxes[a] = b
			}
			err := b.add(ir, ia)
			if err != nil {
				log.Errorf("render-grid: %v", err)
				return
			}
		}
	}
	for a, b := range boxes {
		if err := b.validate(); err != nil {
			log.Errorf("validation for %q: %v", a, err)
			return
		}
	}
}
