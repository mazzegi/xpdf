package xpdf

import (
	"testing"

	"github.com/mazzegi/xpdf/style"
)

func TestTableSpans(t *testing.T) {
	cell := func(colSpan, rowSpan int) *tableCell {
		sty := style.Styles{}
		return &tableCell{
			Styles:  sty,
			colSpan: colSpan,
			rowSpan: rowSpan,
		}
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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))

	desc = "start span and end span"
	tab = &table{
		rows: []*tableRow{
			{
				cells: []*tableCell{
					cell(1, 3), cell(1, 1), cell(1, 1),
				},
			},
			{
				cells: []*tableCell{
					cell(1, 1), cell(1, 1), cell(1, 3),
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
	t.Logf("*** %s\n%s", desc, dumpTableSpans(tab))
}
