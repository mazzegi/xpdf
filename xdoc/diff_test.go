package xdoc

import (
	"testing"

	"github.com/mazzegi/xpdf/style"
)

func TestDiff(t *testing.T) {
	org := style.Styles{}

	mod := style.Styles{}
	mod.Font.Family = "space-font"
	mod.Box.Padding.Left = 2

	sdiff, err := stylesDiff(org, mod)
	if err != nil {
		t.Fatalf("diff-error: %v", err)
	}
	t.Logf("diff: %s", sdiff)
}
