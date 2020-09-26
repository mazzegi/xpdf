package style

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestMutate(t *testing.T) {
	mutated := func(t *testing.T, phrase string, sty Styles) (Styles, error) {
		ms := sty
		m, err := DecodeMutator(bytes.NewBufferString(phrase))
		if err != nil {
			return ms, err
		}
		m.Mutate(&ms)
		return ms, nil
	}

	tests := []struct {
		name       string
		inStyles   Styles
		phrase     string
		outStyles  Styles
		decodeFail bool
	}{
		{
			name:     "padding",
			inStyles: Styles{},
			phrase:   "padding: 1,2,3,4",
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
			name:       "padding fail",
			inStyles:   Styles{},
			phrase:     "padding: 1,2,3",
			outStyles:  Styles{},
			decodeFail: true,
		},
		{
			name:     "font",
			inStyles: Styles{},
			phrase:   "font-family: space-type; font-point-size: 14; font-style: italic; font-weight: bold; font-decoration: underline;",
			outStyles: Styles{
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := mutated(t, test.phrase, test.inStyles)
			if test.decodeFail {
				if err == nil {
					t.Fatalf("decode %s should fail but did not", test.phrase)
				}
				return
			}
			if err != nil {
				t.Fatalf("decode %s failed (where it should not): %v", test.phrase, err)
			}
			if res != test.outStyles {
				t.Fatalf("have:\n%s\nwant:\n%s", dumpStyles(res), dumpStyles(test.outStyles))
			}
		})
	}
}

func dumpStyles(sty Styles) string {
	bs, _ := json.MarshalIndent(sty, "", "  ")
	return string(bs)
}
