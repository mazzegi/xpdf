package xpdf

import "github.com/mazzegi/xpdf/style"

func DefaultStyle() style.Styles {
	return style.Styles{
		Font: style.Font{
			Family:     "Arial",
			PointSize:  12,
			Style:      style.FontStyleNormal,
			Weight:     style.FontWeightNormal,
			Decoration: style.FontDecorationNormal,
		},
		Box: style.Box{
			Border:  style.Border{Left: 0, Top: 0, Right: 0, Bottom: 0},
			Padding: style.Padding{Left: 0, Top: 0, Right: 0, Bottom: 0},
			Margin:  style.Margin{Left: 0, Top: 0, Right: 0, Bottom: 0},
		},
		Dimension: style.Dimension{
			Width:       -1,
			Height:      -1,
			LineSpacing: 1.5,
			OffsetX:     0,
			OffsetY:     0,
		},
		Table: style.Table{
			ColumnWidth: -1,
		},
		Align: style.Align{
			HAlign: style.HAlignLeft,
			VAlign: style.VAlignTop,
		},
		Color: style.Color{
			Foreground: style.Black,
			Text:       style.Black,
			Background: style.White,
		},
		Draw: style.Draw{
			LineWidth: 1,
		},
	}
}
