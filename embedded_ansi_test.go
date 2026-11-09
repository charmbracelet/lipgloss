package lipgloss

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

// TestSpaceStylerPreservesEmbeddedANSI verifies that when a space styler is
// active (underline or strikethrough), Render does not mangle ANSI escape
// sequences already present in the input, such as a previously-rendered string.
// See https://github.com/charmbracelet/lipgloss/issues/233.
func TestSpaceStylerPreservesEmbeddedANSI(t *testing.T) {
	for _, tc := range []struct {
		name  string
		style Style
	}{
		{"strikethrough", NewStyle().Strikethrough(true)},
		{"underline", NewStyle().Underline(true)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			inner := NewStyle().Bold(true).Render("hi") // "\x1b[1mhi\x1b[m"
			got := tc.style.Render(inner)

			if w := ansi.StringWidth(got); w != 2 {
				t.Errorf("visible width = %d, want 2 (embedded ANSI leaked as text): %q", w, got)
			}
			if !strings.Contains(got, "\x1b[1m") {
				t.Errorf("embedded bold sequence was mangled, not preserved verbatim: %q", got)
			}
		})
	}
}
