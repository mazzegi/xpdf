package hyphenation

import (
	"bytes"
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

func TestParseSet(t *testing.T) {
	buf := bytes.NewBufferString(enUsPatterns)
	t0 := time.Now()
	ps, err := parsePatterns(buf)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	t.Logf("parsed %d patterns in %s", len(ps.patterns), time.Since(t0))
}

func TestSubWords(t *testing.T) {
	s := ".hyphenation."
	size := 1
	subs := subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	size = 2
	subs = subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	size = 3
	subs = subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	size = 8
	subs = subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	size = 13
	subs = subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	size = 15
	subs = subWords(s, size)
	t.Logf("subs(%d): %v", size, subs)

	allSubs := allSubWords(s)
	t.Logf("all-subs: %v", allSubs)
}

func TestPatternMatch(t *testing.T) {
	buf := bytes.NewBufferString(enUsPatterns)
	pl, err := parsePatterns(buf)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	s := ".hyphenation."
	ps := matchingPatterns(pl, s)
	for _, p := range ps {
		t.Logf("p: %q", p.String())
	}
	wps := weightedWordPattern(ps, s)
	t.Logf("weighted-word-pattern: %q", wps.String())

	s = ".concatenation."
	ps = matchingPatterns(pl, s)
	for _, p := range ps {
		t.Logf("p: %q", p.String())
	}
	wps = weightedWordPattern(ps, s)
	t.Logf("weighted-word-pattern: %q", wps.String())

	s = ".supercalifragilisticexpialidocious."
	ps = matchingPatterns(pl, s)
	for _, p := range ps {
		t.Logf("p: %q", p.String())
	}
	wps = weightedWordPattern(ps, s)
	t.Logf("weighted-word-pattern: %q", wps.String())
}

func TestHyhenation(t *testing.T) {
	buf := bytes.NewBufferString(enUsPatterns)
	pl, err := parsePatterns(buf)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	s := "hyphenation"
	t0 := time.Now()
	hsl := hyphenated(pl, s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	s = "concatenation"
	t0 = time.Now()
	hsl = hyphenated(pl, s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))

	s = "supercalifragilisticexpialidocious"
	t0 = time.Now()
	hsl = hyphenated(pl, s)
	t.Logf("hyph: %v (%s)", hsl, time.Since(t0))
}
