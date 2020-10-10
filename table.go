package xpdf

import (
	"fmt"
	"strings"

	"github.com/mazzegi/xpdf/style"
	"github.com/mazzegi/xpdf/xdoc"
)

func dumpTableSpans(tab *table) string {
	rowsl := []string{}
	for _, row := range tab.rows {
		rowsl = append(rowsl, dumpRowSpans(row))
	}
	return strings.Join(rowsl, "\n")
}

func dumpRowSpans(row *tableRow) string {
	cellsl := []string{}
	for _, cell := range row.cells {
		if cell.zero {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-zero-cell]", cell.rowIdx, cell.cellIdx))
		} else if cell.spannedBy != nil {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-byspn:%d:%d]", cell.rowIdx, cell.cellIdx, cell.spannedBy.rowIdx, cell.spannedBy.cellIdx))
		} else if len(cell.spansCols) > 0 && len(cell.spansRows) > 0 {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-spans:%d:%d]", cell.rowIdx, cell.cellIdx, cell.colSpan, cell.rowSpan))
		} else if len(cell.spansCols) > 0 {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-colsp:%d:%d]", cell.rowIdx, cell.cellIdx, cell.colSpan, cell.rowSpan))
		} else if len(cell.spansRows) > 0 {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-rowsp:%d:%d]", cell.rowIdx, cell.cellIdx, cell.colSpan, cell.rowSpan))
		} else {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-regl-cell]", cell.rowIdx, cell.cellIdx))
		}
	}
	return strings.Join(cellsl, ", ")
}

func dumpTableDims(tab *table) string {
	rowsl := []string{}
	for _, row := range tab.rows {
		rowsl = append(rowsl, dumpRowDims(row))
	}
	return strings.Join(rowsl, "\n")
}

func dumpRowDims(row *tableRow) string {
	cellsl := []string{}
	for _, cell := range row.cells {
		cellsl = append(cellsl, fmt.Sprintf("[%d:%d:w=%.1f:h=%.1f:mh=%.1f]", cell.rowIdx, cell.cellIdx, cell.width, cell.height, cell.minHeight))
	}
	return strings.Join(cellsl, ", ")
}

//

type table struct {
	style.Styles
	rows []*tableRow

	// auxiliary parameters
	columnCount int
}

type tableRow struct {
	style.Styles
	cells []*tableCell

	// auxiliary parameters
	height float64
}

func (row *tableRow) maxCellHeight() float64 {
	var max float64
	for _, cell := range row.cells {
		_, ch := cell.dim()
		if ch > max {
			max = ch
		}
	}
	return max
}

type tableCell struct {
	style.Styles
	text string
	iss  []xdoc.Instruction

	// auxiliary parameters
	colSpan   int
	rowSpan   int
	rowIdx    int
	cellIdx   int
	minHeight float64
	height    float64
	width     float64
	spannedBy *tableCell
	spansCols []*tableCell
	spansRows []*tableCell
	zero      bool
}

func (c *tableCell) dim() (width, height float64) {
	width = c.width
	height = c.height
	for _, sc := range c.spansCols {
		width += sc.width
	}
	for _, sc := range c.spansRows {
		height += sc.height
	}
	return
}

//

func (t *table) assignColumnWidths(width float64) {
	//TODO: think about considering column-widths as fractions from a sum of widths to avoid overflowing
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
			cell.width = cws[ic]
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
func (t *table) processColumnSpans() {
	for _, row := range t.rows {
		newCells := []*tableCell{}
		for _, cell := range row.cells {
			newCells = append(newCells, cell)
			if cell.colSpan <= 1 {
				continue
			}
			for s := 0; s < cell.colSpan-1; s++ {
				spannedCell := &tableCell{
					spannedBy: cell,
				}
				spannedCell.rowSpan = cell.rowSpan
				newCells = append(newCells, spannedCell)
				cell.spansCols = append(cell.spansCols, spannedCell)
			}
		}
		row.cells = newCells
	}
}

func (t *table) processRowSpans() {
	for ir, row := range t.rows {
		for ic, cell := range row.cells {
			if cell.rowSpan <= 1 {
				continue
			}
			if cell.spannedBy != nil || cell.zero {
				continue
			}
			//insert spanned cell in following rows
			for n := 0; n < cell.rowSpan-1; n++ {
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
						zero: true,
					})
				}

				newCell := &tableCell{
					spannedBy: cell,
				}
				cell.spansRows = append(cell.spansRows, newCell)

				newCells := make([]*tableCell, ic)
				copy(newCells, spannedRow.cells[:ic])
				newCells = append(newCells, newCell)
				newCells = append(newCells, spannedRow.cells[ic:]...)
				spannedRow.cells = newCells
			}
		}
	}
}

func (t *table) processSpans() {
	t.processColumnSpans()
	t.processRowSpans()
	t.normalize()
}

func (p *Processor) assignCellMinHeight(cell *tableCell) {
	if cell.spannedBy != nil || cell.zero {
		return
	}
	if cell.Height > 0 {
		// if set by style
		cell.minHeight = cell.Height
		return
	}
	availableWidth := cell.width
	for _, sc := range cell.spansCols {
		availableWidth += sc.width
	}
	availableWidth -= cell.Padding.Left + cell.Padding.Right

	textHeight := p.textHeightFnc(cell.Styles)(cell.text, availableWidth, cell.Styles)

	cellHeight := textHeight + cell.Padding.Top + cell.Padding.Bottom

	// if cell spans rows, divide height to spanned cells
	heightPerCell := cellHeight / float64(1+len(cell.spansRows))
	cell.minHeight = heightPerCell
	for _, sc := range cell.spansRows {
		sc.minHeight = heightPerCell
	}
}

func (p *Processor) assignHeights(tab *table) {
	for _, row := range tab.rows {
		for _, cell := range row.cells {
			p.assignCellMinHeight(cell)
		}
	}
	//no check for each row
	for _, row := range tab.rows {
		var maxCellHeight float64
		for _, cell := range row.cells {
			if cell.minHeight > maxCellHeight {
				maxCellHeight = cell.minHeight
			}
		}
		row.height = maxCellHeight
	}

	for _, row := range tab.rows {
		for _, cell := range row.cells {
			cell.height = row.height
		}
	}
}

func (p *Processor) transformTable(xtab *xdoc.Table) *table {
	tab := &table{
		Styles: xtab.MutatedStyles(p.doc.StyleClasses(), p.currStyles),
	}
	for ir, xrow := range xtab.Rows {
		var rowSty style.Styles
		switch {
		case ir == 0:
			rowSty = xrow.MutatedStylesWithSelector("first-row", p.doc.StyleClasses(), tab.Styles)
		case ir == len(xtab.Rows)-1:
			rowSty = xrow.MutatedStylesWithSelector("last-row", p.doc.StyleClasses(), tab.Styles)
		default:
			rowSty = xrow.MutatedStyles(p.doc.StyleClasses(), tab.Styles)
		}
		row := &tableRow{
			Styles: rowSty,
		}
		for ic, xcell := range xrow.Cells {
			var cellSty style.Styles
			switch {
			case ic == 0:
				cellSty = xcell.MutatedStylesWithSelector("first-cell", p.doc.StyleClasses(), row.Styles)
			case ic == len(xrow.Cells)-1:
				cellSty = xcell.MutatedStylesWithSelector("last-cell", p.doc.StyleClasses(), row.Styles)
			default:
				cellSty = xcell.MutatedStyles(p.doc.StyleClasses(), row.Styles)
			}
			cell := &tableCell{
				Styles:  cellSty,
				colSpan: xcell.ColSpan,
				rowSpan: xcell.RowSpan,
				text:    xcell.Content,
				iss:     xcell.Instructions,
			}

			row.cells = append(row.cells, cell)
		}
		tab.rows = append(tab.rows, row)
	}
	tab.processSpans()
	tab.assignColumnWidths(p.page().EffectiveWidth(tab.Width))
	p.assignHeights(tab)
	//TODO: reapply styles as first/last row/cell may have changed
	return tab
}

//

func (p *Processor) renderTable(xtab *xdoc.Table) {
	defer p.resetStyles()
	tab := p.transformTable(xtab)
	if tab.columnCount == 0 {
		return
	}

	//TODO: add a "repeat first row on page-break" option
	page := p.page()
	x0, y := p.engine.GetXY()
	for _, row := range tab.rows {
		if !p.preventPageBreak && y+row.maxCellHeight() > page.printableArea.y1 {
			p.engine.AddPage()
			_, y = p.engine.GetXY()
		}
		x := x0
		for _, cell := range row.cells {
			if cell.spannedBy != nil || cell.zero {
				x += cell.width
				continue
			}
			cw, ch := cell.dim()
			p.renderCell(PrintableArea{
				x0: x,
				y0: y,
				x1: x + cw,
				y1: y + ch,
			}, cell)
			x += cell.width
		}
		y += row.height
		p.engine.SetX(x0)
	}
	p.engine.SetY(y)
}

func (p *Processor) renderCell(pa PrintableArea, cell *tableCell) {
	p.drawBox(pa.x0, pa.y0, pa.x1, pa.y1, cell.Styles)

	paddedPa := pa.WithPadding(cell.Padding)
	if cell.text != "" {
		textHeight := p.textHeight(cell.text, paddedPa.Width(), cell.Styles)
		textMargin := paddedPa.Height() - textHeight
		switch cell.VAlign {
		case style.VAlignMiddle:
			p.engine.SetY(paddedPa.y0 + textMargin/2)
		case style.VAlignBottom:
			p.engine.SetY(paddedPa.y0 + textMargin)
		default: //style.VAlignTop:
			p.engine.SetY(paddedPa.y0)
		}
		p.engine.SetX(paddedPa.x0)

		p.writeTextFnc(cell.Styles)(cell.text, pa.Width()-cell.Padding.Left-cell.Padding.Right, cell.Styles)
	}

	for _, is := range cell.iss {
		switch is := is.(type) {
		case *xdoc.Box:
			height := p.textBoxHeight(is, paddedPa)
			switch cell.VAlign {
			case style.VAlignMiddle:
				p.engine.SetY(paddedPa.y0 + height/2)
			case style.VAlignBottom:
				p.engine.SetY(paddedPa.y0 + height)
			default: //style.VAlignTop:
				p.engine.SetY(paddedPa.y0)
			}
			p.engine.SetX(paddedPa.x0)
			p.renderTextBox(is, paddedPa)
		case *xdoc.Image:
			p.engine.SetY(paddedPa.y0)
			p.engine.SetX(paddedPa.x0)
			p.renderImage(is, paddedPa)
		}
	}
}
