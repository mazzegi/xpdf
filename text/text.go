package text

import "strings"

func RemoveTwins(s string, twin string) string {
	rem := strings.Repeat(twin, 2)
	rs := s
	for {
		rs = strings.ReplaceAll(s, rem, twin)
		if rs == s {
			return rs
		}
		s = rs
	}
}

func WhitespaceRectified(s string) string {
	rs := s
	rs = strings.ReplaceAll(rs, "\r", " ")
	rs = strings.ReplaceAll(rs, "\n", " ")
	rs = strings.ReplaceAll(rs, "\t", " ")
	rs = strings.Trim(rs, " ")
	rs = RemoveTwins(rs, " ")
	rs = strings.ReplaceAll(rs, " \\", "\\")
	rs = strings.ReplaceAll(rs, "\\ ", "\\")
	return rs
}
