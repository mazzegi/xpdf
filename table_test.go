package xpdf

import (
	"fmt"
	"strings"
	"testing"

	"github.com/mazzegi/xpdf/style"
)

func TestTableSpans(t *testing.T) {
	cell := func(colSpan, rowSpan int) *tableCell {
		sty := style.Styles{}
		sty.ColumnSpan = colSpan
		sty.RowSpan = rowSpan
		return &tableCell{Styles: sty}
	}
	desc := "some col spans"
	tab := &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(2, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(2, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "some row spans"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 2), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "combined col+row spans"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(2, 2), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "start col span"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(2, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "start row span"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 2), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "end row span"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 2),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "end row span overlap"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 2),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "end row span empty overlap"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 2),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))

	desc = "end row span empty overlap and overlap rows"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 3),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1),
				},
			},
		},
	}
	tab.processSpans()
	t.Logf("*** %s\n%s", desc, dumpTable(tab))
}

func dumpTable(tab *table) string {
	rowsl := []string{}
	for _, row := range tab.rows {
		rowsl = append(rowsl, dumpRow(row))
	}
	return strings.Join(rowsl, "\n")
}

func dumpRow(row *tableRow) string {
	cellsl := []string{}
	for _, cell := range row.cells {
		if cell.zero {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-zero-cell]", cell.rowIdx, cell.cellIdx))
		} else if cell.spannedBy != nil {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-byspn:%d:%d]", cell.rowIdx, cell.cellIdx, cell.spannedBy.rowIdx, cell.spannedBy.cellIdx))
		} else if len(cell.spans) > 0 {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-spans:%d:%d]", cell.rowIdx, cell.cellIdx, cell.ColumnSpan, cell.RowSpan))
		} else {
			cellsl = append(cellsl, fmt.Sprintf("[%d:%d-regl-cell]", cell.rowIdx, cell.cellIdx))
		}
	}
	return strings.Join(cellsl, ", ")
}
