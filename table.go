package xpdf

import (
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

type table struct {
	style.Styles
	rows []*tableRow

	maxColumnCount int
}

type tableRow struct {
	style.Styles
	cells []*tableCell
}

type tableCell struct {
	style.Styles
	text string
	iss  []xdoc.Instruction

	spannedBy *tableCell
	spans     []*tableCell
}

func (t *table) determineColumnWidths(width float64) {
	cws := make([]float64, t.maxColumnCount)
	for _, row := range t.rows {
		for ic, cell := range row.cells {
			if cell.ColumnWidth > cws[ic] {
				cws[ic] = cell.ColumnWidth
			}
		}
	}
	var spaceUsed float64
	var colsZero int
	for _, cw := range cws {
		if cw > 0 {
			spaceUsed += cw
		} else {
			colsZero++
		}
	}
	if colsZero > 0 {
		cw := (width - spaceUsed) / float64(colsZero)
		for i := range cws {
			if cws[i] == 0 {
				cws[i] = cw
			}
		}
	}
}

//
func (t *table) processRowSpans() {
	for ir, row := range t.rows {
		for ic, cell := range row.cells {
			if cell.RowSpan <= 1 {
				continue
			}
			//insert spanned cell in following rows
			for n := 0; n < cell.RowSpan-1; n++ {
				spannedRowIdx := ir + 1 + n
				if spannedRowIdx >= len(t.rows) {
					//row span exceeds table - just do nothing
					continue
				}
				spannedRow := t.rows[spannedRowIdx]
				if ic > len(spannedRow.cells) {
					//if anyway cell-idx is bigger than spanned-row max cell-idx nothing has to be done (but later in render)
					continue
				}

				newCell := &tableCell{spannedBy: cell}
				cell.spans = append(cell.spans, newCell)

				newCells := append(spannedRow.cells[:ic], newCell)
				newCells = append(newCells, spannedRow.cells[ic:]...)
				spannedRow.cells = newCells
			}
		}
	}
}

func (p *Processor) transformTable(xtab *xdoc.Table) *table {
	tab := &table{
		Styles: xtab.MutatedStyles(p.doc.StyleClasses(), p.currStyles),
	}
	for _, xrow := range xtab.Rows {
		row := &tableRow{
			Styles: xrow.MutatedStyles(p.doc.StyleClasses(), tab.Styles),
		}
		for _, xcell := range xrow.Cells {
			cell := &tableCell{
				Styles: xcell.MutatedStyles(p.doc.StyleClasses(), row.Styles),
				text:   xcell.Content,
				iss:    xcell.Instructions,
			}

			row.cells = append(row.cells, cell)
		}
		tab.rows = append(tab.rows, row)
		if len(row.cells) > tab.maxColumnCount {
			tab.maxColumnCount = len(row.cells)
		}
	}
	tab.processRowSpans()
	return tab
}

//

func (p *Processor) renderTable(xtab *xdoc.Table) {
	defer p.resetStyles()
	tab := p.transformTable(xtab)
	if tab.maxColumnCount == 0 {
		return
	}

	_ = tab
}
