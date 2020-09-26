package style

import (
	"bytes"
	"testing"
)

func TestClasses(t *testing.T) {
	mutated := func(t *testing.T, classes string, useClasses []string, sty Styles) (Styles, error) {
		ms := sty
		css, err := DecodeClasses(bytes.NewBufferString(classes))
		if err != nil {
			return ms, err
		}
		css.Mutate(&ms, useClasses...)
		return ms, nil
	}

	tests := []struct {
		name       string
		inStyles   Styles
		classes    string
		useClases  []string
		outStyles  Styles
		decodeFail bool
	}{
		{
			name:     "classes #1",
			inStyles: Styles{},
			classes: `
			padded{
				padding: 1,2,3,4;
			}
			bold-unerline{
				font-family: space-type;
				font-point-size: 14;
				font-style: italic;
				font-weight: bold;
				font-decoration: underline;
			}
			`,
			useClases: []string{"padded", "bold-unerline"},
			outStyles: Styles{
				Box: Box{
					Padding: Padding{
						Left:   1,
						Top:    2,
						Right:  3,
						Bottom: 4,
					},
				},
				Font: Font{
					Family:     "space-type",
					PointSize:  14,
					Style:      "italic",
					Weight:     "bold",
					Decoration: "underline",
				},
			},
			decodeFail: false,
		},
		{
			name:       "classes fail",
			inStyles:   Styles{},
			classes:    `xt{`,
			outStyles:  Styles{},
			decodeFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := mutated(t, test.classes, test.useClases, test.inStyles)
			if test.decodeFail {
				if err == nil {
					t.Fatalf("decode %s should fail but did not", test.classes)
				}
				return
			}
			if err != nil {
				t.Fatalf("decode %s failed (where it should not): %v", test.classes, err)
			}
			if res != test.outStyles {
				t.Fatalf("have:\n%s\nwant:\n%s", dumpStyles(res), dumpStyles(test.outStyles))
			}
		})
	}
}

func TestClassesSelectors(t *testing.T) {
	mutated := func(t *testing.T, classes string, useClasses []string, selector string, sty Styles) (Styles, error) {
		ms := sty
		css, err := DecodeClasses(bytes.NewBufferString(classes))
		if err != nil {
			return ms, err
		}
		css.MutateWithSelector(selector, &ms, useClasses...)
		return ms, nil
	}

	tests := []struct {
		name       string
		inStyles   Styles
		classes    string
		selector   string
		useClases  []string
		outStyles  Styles
		decodeFail bool
	}{
		{
			name:     "classes - no selector #1",
			inStyles: Styles{},
			classes: `
			padded{
				padding: 1,2,3,4;
			}
			bold-unerline{
				font-family: space-type;
				font-point-size: 14;
				font-style: italic;
				font-weight: bold;
				font-decoration: underline;
			}
			`,
			useClases: []string{"padded", "bold-unerline"},
			outStyles: Styles{
				Box: Box{
					Padding: Padding{
						Left:   1,
						Top:    2,
						Right:  3,
						Bottom: 4,
					},
				},
				Font: Font{
					Family:     "space-type",
					PointSize:  14,
					Style:      "italic",
					Weight:     "bold",
					Decoration: "underline",
				},
			},
			decodeFail: false,
		},
		{
			name:     "classes - selector not applied #1",
			inStyles: Styles{},
			classes: `
			padded{
				padding: 1,2,3,4;
			}
			padded:reverse{
				padding: 4,3,2,1;
			}			
			`,
			useClases: []string{"padded"},
			outStyles: Styles{
				Box: Box{
					Padding: Padding{
						Left:   1,
						Top:    2,
						Right:  3,
						Bottom: 4,
					},
				},
			},
			decodeFail: false,
		},
		{
			name:     "classes - selector applied #1",
			inStyles: Styles{},
			classes: `
			padded{
				padding: 1,2,3,4;
			}
			padded:reverse{
				padding: 4,3,2,1;
			}			
			`,
			useClases: []string{"padded"},
			selector:  "reverse",
			outStyles: Styles{
				Box: Box{
					Padding: Padding{
						Left:   4,
						Top:    3,
						Right:  2,
						Bottom: 1,
					},
				},
			},
			decodeFail: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := mutated(t, test.classes, test.useClases, test.selector, test.inStyles)
			if test.decodeFail {
				if err == nil {
					t.Fatalf("decode %s should fail but did not", test.classes)
				}
				return
			}
			if err != nil {
				t.Fatalf("decode %s failed (where it should not): %v", test.classes, err)
			}
			if res != test.outStyles {
				t.Fatalf("have:\n%s\nwant:\n%s", dumpStyles(res), dumpStyles(test.outStyles))
			}
		})
	}
}
