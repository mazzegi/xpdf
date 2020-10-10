package hyphenation

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type pattern struct {
	Letters []rune
	Weights []int
}

type patternLookup struct {
	patterns map[string]pattern
}

func newPatternLookup() *patternLookup {
	return &patternLookup{
		patterns: map[string]pattern{},
	}
}

func loadPatternLookup(r io.Reader) (*patternLookup, error) {
	return parsePatterns(r)
}

func (pl *patternLookup) find(key string) (pattern, bool) {
	p, ok := pl.patterns[key]
	return p, ok
}

func (p pattern) String() string {
	var s string
	for i := 0; i < len(p.Weights); i++ {
		s += fmt.Sprintf("%d", p.Weights[i])
		if i < len(p.Letters) {
			s += fmt.Sprintf("%c", p.Letters[i])
		}
	}
	return s
}

func parsePattern(s string) (pattern, error) {
	p := pattern{}
	wantDigit := true
	for _, r := range s {
		if wantDigit {
			if unicode.IsDigit(r) {
				w, err := strconv.ParseInt(string(r), 10, 8)
				if err != nil {
					return p, errors.Wrap(err, "while expecting digit")
				}
				p.Weights = append(p.Weights, int(w))
				wantDigit = false
			} else {
				p.Weights = append(p.Weights, 0)
				p.Letters = append(p.Letters, r)
				wantDigit = true
			}
		} else {
			p.Letters = append(p.Letters, r)
			wantDigit = true
		}
	}
	if wantDigit {
		p.Weights = append(p.Weights, 0)
	}
	if len(p.Letters)+1 != len(p.Weights) {
		return p, errors.Errorf("invalid pattern with %d letters and %d weights", len(p.Letters), len(p.Weights))
	}
	return p, nil
}

func parsePatterns(r io.Reader) (*patternLookup, error) {
	pl := newPatternLookup()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Trim(line, " \n\r\t")
		if line == "" {
			continue
		}
		p, err := parsePattern(line)
		if err != nil {
			return nil, errors.Wrap(err, "parse-pattern")
		}
		pl.patterns[string(p.Letters)] = p
	}
	return pl, nil
}
