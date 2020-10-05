package hyphenation

import "strings"

func subWords(s string, size int) []string {
	subs := []string{}
	for i := 0; i < len(s)-size+1; i++ {
		var sub string
		for k := i; k < i+size; k++ {
			sub += string(s[k])
		}
		subs = append(subs, sub)
	}
	return subs
}

func allSubWords(s string) []string {
	subs := []string{}
	if len(s) == 0 {
		return subs
	}
	for size := 1; size <= len(s); size++ {
		subs = append(subs, subWords(s, size)...)
	}
	return subs
}

func matchingPatterns(pl *PatternLookup, s string) []Pattern {
	subs := allSubWords(s)
	ps := []Pattern{}
	for _, sub := range subs {
		//TODO: probably stop if all sub words up to here are matched
		if p, ok := pl.Find(sub); ok {
			ps = append(ps, p)
		}
	}
	return ps
}

func maxWeight(r1, r2 rune, ps []Pattern) int {
	max := 0
	for _, p := range ps {
		for i := 0; i < len(p.Letters); i++ {
			var weight int
			if i == 0 && p.Letters[i] == r2 {
				//r1 matches any in case of the first letter
				weight = p.Weights[0]
			} else if i == len(p.Letters)-1 && p.Letters[i] == r1 {
				//r2 matches any in case of the first letter
				weight = p.Weights[i+1]
			} else if p.Letters[i] == r1 && p.Letters[i+1] == r2 {
				weight = p.Weights[i+1]
			}

			if weight > max {
				max = weight
			}
		}
	}
	return max
}

func weightedWordPattern(ps []Pattern, s string) Pattern {
	wp := Pattern{}
	wp.Weights = append(wp.Weights, 0)
	rs := []rune(s)
	for i := 0; i < len(rs)-1; i++ {
		r1 := rs[i]
		r2 := rs[i+1]
		mw := maxWeight(r1, r2, ps)
		wp.Letters = append(wp.Letters, r1)
		wp.Weights = append(wp.Weights, mw)
	}
	wp.Letters = append(wp.Letters, rs[len(rs)-1])
	wp.Weights = append(wp.Weights, 0)
	return wp
}

func hyphenated(pl *PatternLookup, s string) []string {
	s = "." + s + "."
	ps := matchingPatterns(pl, s)
	sl := []string{}
	var curr string
	rs := []rune(s)
	for i := 0; i < len(rs)-1; i++ {
		r1 := rs[i]
		r2 := rs[i+1]
		mw := maxWeight(r1, r2, ps)
		curr += string(r1)
		if mw%2 == 1 {
			//possible hyphen
			curr = strings.Trim(curr, ". ")
			if curr != "" {
				sl = append(sl, curr)
			}
			curr = ""
		}
	}
	curr = strings.Trim(curr, ". ")
	if curr != "" {
		sl = append(sl, curr)
	}
	return sl
}
