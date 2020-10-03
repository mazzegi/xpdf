package xpdf

import (
	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

type table struct {
	style.Styles
	rows []*tableRow

	columnCount int
}

type tableRow struct {
	style.Styles
	cells []*tableCell
}

type tableCell struct {
	style.Styles
	text string
	iss  []xdoc.Instruction

	rowIdx    int
	cellIdx   int // in row
	width     float64
	spannedBy *tableCell
	spans     []*tableCell
	zero      bool
}

func (t *table) assignColumnWidths(width float64) {
	cws := make([]float64, t.columnCount)
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
	for _, row := range t.rows {
		for ic, cell := range row.cells {
			for s := 0; s < cell.ColumnSpan; s++ {
				cell.width = cws[ic]
			}
		}
	}
}

func (tab *table) normalize() {
	tab.columnCount = 0
	for _, row := range tab.rows {
		if len(row.cells) > tab.columnCount {
			tab.columnCount = len(row.cells)
		}
	}
	for _, row := range tab.rows {
		for len(row.cells) < tab.columnCount {
			row.cells = append(row.cells, &tableCell{
				zero: true,
			})
		}
	}
	for ir, row := range tab.rows {
		for ic, cell := range row.cells {
			cell.rowIdx = ir
			cell.cellIdx = ic
		}
	}
}

//
func (t *table) processSpans() {
	//col spans
	for ir, row := range t.rows {
		newCells := []*tableCell{}
		for _, cell := range row.cells {
			cell.rowIdx = ir
			cell.cellIdx = len(newCells)
			newCells = append(newCells, cell)
			if cell.ColumnSpan <= 1 {
				continue
			}
			for s := 0; s < cell.ColumnSpan-1; s++ {
				spannedCell := &tableCell{
					spannedBy: cell,
					rowIdx:    ir,
					cellIdx:   len(newCells),
				}
				spannedCell.RowSpan = cell.RowSpan
				newCells = append(newCells, spannedCell)
				cell.spans = append(cell.spans, spannedCell)
			}
		}
		row.cells = newCells
	}
	//row spans
	for ir, row := range t.rows {
		for ic, cell := range row.cells {
			cell.cellIdx = ic
			if cell.RowSpan <= 1 {
				continue
			}
			if cell.spannedBy != nil || cell.zero {
				continue
			}
			//insert spanned cell in following rows
			for n := 0; n < cell.RowSpan-1; n++ {
				spannedRowIdx := ir + 1 + n

				var spannedRow *tableRow
				if spannedRowIdx >= len(t.rows) {
					spannedRow = &tableRow{}
					t.rows = append(t.rows, spannedRow)
				} else {
					spannedRow = t.rows[spannedRowIdx]
				}
				for ic > len(spannedRow.cells) {
					spannedRow.cells = append(spannedRow.cells, &tableCell{
						zero:   true,
						rowIdx: spannedRowIdx,
					})
				}

				newCell := &tableCell{
					spannedBy: cell,
					rowIdx:    spannedRowIdx,
				}
				cell.spans = append(cell.spans, newCell)

				newCells := make([]*tableCell, ic)
				copy(newCells, spannedRow.cells[:ic])
				newCells = append(newCells, newCell)
				newCells = append(newCells, spannedRow.cells[ic:]...)
				spannedRow.cells = newCells
			}
		}
	}
	t.normalize()
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
	}
	tab.processSpans()
	tab.assignColumnWidths(p.engine.EffectiveWidth(tab.Width))
	return tab
}

//

func (p *Processor) renderTable(xtab *xdoc.Table) {
	defer p.resetStyles()
	tab := p.transformTable(xtab)
	if tab.columnCount == 0 {
		return
	}

	_ = tab
}
