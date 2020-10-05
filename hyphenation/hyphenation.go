package hyphenation

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
		if p, ok := pl.Find(sub); ok {
			ps = append(ps, p)
		}
	}
	return ps
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

func maxWeight(r1, r2 rune, ps []Pattern) int {
	max := 0
	for _, p := range ps {
		for i := 0; i < len(p.Letters)-1; i++ {
			if p.Letters[i] == r1 && p.Letters[i+1] == r2 {
				weight := p.Weights[i+1]
				if weight > max {
					max = weight
				}
			}
		}
	}
	return max
}
