package tree_test

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"charm.land/lipgloss/v2/tree"
	"github.com/charmbracelet/x/ansi"
)

func TestRootStyleWidthOnce(t *testing.T) {
	transform := lipgloss.NewStyle().Transform(func(s string) string { return s + "!" })
	padding := lipgloss.NewStyle().Padding(0, 1)
	for _, tc := range []struct {
		name, root, want string
		width            int
		style            lipgloss.Style
	}{
		{"plain", "Root", "Root", 0, lipgloss.NewStyle()},
		{"plain-padded", "Root", "Root      ", 10, lipgloss.NewStyle()},
		{"transform-no-width", "Root", "Root!", 0, transform},
		{"transform-narrow-width", "Root", "Root!", 3, transform},
		{"transform-padded", "Root", "Root!     ", 10, transform},
		{"style-padding-no-width", "Root", " Root ", 0, padding},
		{"style-padding-padded", "Root", " Root     ", 10, padding},
		{"unicode-narrow-width", "界", "界!", 2, transform},
		{"unicode-padded", "界", "界!     ", 8, transform},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tr := tree.Root(tc.root).Child("child").Width(tc.width).RootStyle(tc.style)
			line, _, _ := strings.Cut(ansi.Strip(tr.String()), "\n")
			if line != tc.want {
				t.Errorf("root = %q (width %d), want %q (width %d)",
					line, lipgloss.Width(line), tc.want, lipgloss.Width(tc.want))
			}
		})
	}
}
