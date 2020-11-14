package hyphenation

import (
	"io"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func LoadFromFile(file string) (*Hyphenator, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.Wrapf(err, "open-file %q", file)
	}
	defer f.Close()
	return Load(f)
}

func Load(r io.Reader) (*Hyphenator, error) {
	pl, err := loadPatternLookup(r)
	if err != nil {
		return nil, errors.Wrap(err, "load-pattern-lookup")
	}
	return new(pl), nil
}

type Hyphenator struct {
	lookup *patternLookup
}

func new(pl *patternLookup) *Hyphenator {
	return &Hyphenator{
		lookup: pl,
	}
}

func (h *Hyphenator) Hyphenate(s string) []string {
	if len(s) < 3 {
		//don't hyphenate words with less than 3 runes
		return []string{s}
	}
	if strings.Contains(s, " ") {
		return []string{s}
	}
	if len(s) > 100 {
		return []string{s}
	}
	rs := []rune("." + strings.ToLower(s) + ".")
	ws := make([]int, len(rs)+1)
	for subSize := 1; subSize <= len(rs); subSize++ {
		for i := 0; i < len(rs)-subSize+1; i++ {
			sub := rs[i : i+subSize]
			pattern, ok := h.lookup.find(string(sub))
			if !ok {
				continue
			}
			for iw, w := range pattern.Weights {
				if w > ws[i+iw] {
					ws[i+iw] = w
				}
			}
		}
	}

	sl := []string{}
	var last int
	//skip first and last for the dots (.).
	//skip next to first and prev. to last for start and end of word
	for i, w := range ws[2 : len(ws)-2] {
		if w%2 == 1 {
			part := s[last : i+1]
			if len(part) > 1 {
				sl = append(sl, part)
				last = i + 1
			}
		}
	}
	part := s[last:]
	if len(part) > 1 {
		sl = append(sl, part)
	}
	return sl
}
