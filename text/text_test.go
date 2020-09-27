package text

import (
	"fmt"
	"testing"
)

func TestRemoveTwins(t *testing.T) {
	tests := []struct {
		in   string
		twin string
		exp  string
	}{
		{"foo bar baz", " ", "foo bar baz"},
		{"foo bar  baz", " ", "foo bar baz"},
		{"foo bar bar baz", "bar", "foo bar bar baz"},
		{"foo barbar baz", "bar", "foo bar baz"},
		{"foo \n\n baz", "\n", "foo \n baz"},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			have := RemoveTwins(test.in, test.twin)
			if have != test.exp {
				t.Fatalf("have %q, want %q", have, test.exp)
			}
		})
	}
}

func TestWhitespaceRectified(t *testing.T) {
	tests := []struct {
		in  string
		exp string
	}{
		{"foo bar baz", "foo bar baz"},
		{"foo bar  baz", "foo bar baz"},
		{"foo bar bar baz", "foo bar bar baz"},
		{"foo \\bar \\ baz", "foo\\bar\\baz"},
		{"foo \n\n baz", "foo baz"},
		{"  \rfoo    \\\\bar \n\n baz\n  ", "foo\\\\bar baz"},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("#%d", i+1), func(t *testing.T) {
			have := WhitespaceRectified(test.in)
			if have != test.exp {
				t.Fatalf("have %q, want %q", have, test.exp)
			}
		})
	}
}
