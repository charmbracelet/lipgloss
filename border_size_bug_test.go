package lipgloss

import "testing"

// TestBorderTopBottomSizeWithWideRune tests that GetTopSize and GetBottomSize
// return row count (always 1), not rune width.
// This test FAILS on main (bug: returns 2 for wide runes) and PASSES on fix branch.
func TestBorderTopBottomSizeWithWideRune(t *testing.T) {
	// Use ⏩ emoji which has display width of 2
	borderWithWideTop := Border{
		Top:      "⏩",
		TopLeft:  "⏩",
		TopRight: "⏩",
	}

	borderWithWideBottom := Border{
		Bottom:      "⏩",
		BottomLeft:  "⏩",
		BottomRight: "⏩",
	}

	// Top/bottom borders always occupy exactly 1 row, regardless of rune display width
	if got := borderWithWideTop.GetTopSize(); got != 1 {
		t.Errorf("GetTopSize() with wide rune = %d, want 1 (bug: returns rune width instead of row count)", got)
	}

	if got := borderWithWideBottom.GetBottomSize(); got != 1 {
		t.Errorf("GetBottomSize() with wide rune = %d, want 1 (bug: returns rune width instead of row count)", got)
	}
}
