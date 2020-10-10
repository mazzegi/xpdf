package hyphenation

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	s := "ab5o5liz"
	p, err := parsePattern(s)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	t.Logf("pattern: %q", p.String())

	s = ".me5ter"
	p, err = parsePattern(s)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	t.Logf("pattern: %q", p.String())
}

func TestHyhenation(t *testing.T) {
	hyp := NewEnUs()

	var s string
	var hsl []string
	var t0 time.Time

	// s = "hyphenation"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "concatenation"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "supercalifragilisticexpialidocious"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "Developer"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	// s = "sportsman"
	// t0 = time.Now()
	// hsl = Hyphenated(pl, s)
	// t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	s = "small"
	t0 = time.Now()
	hsl = hyp.Hyphenate(s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	s = "sportsman"
	t0 = time.Now()
	hsl = hyp.Hyphenate(s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))
}

func TestLatinHyhenation(t *testing.T) {
	hyp := NewLatin()

	s := "Lorem"
	t0 := time.Now()
	hsl := hyp.Hyphenate(s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))
}
