package lipgloss

import (
	"strings"
	"testing"
)

func TestReverseAppliesToSpacesWhenSpaceStylerIsUsed(t *testing.T) {
	tests := []struct {
		name  string
		style Style
	}{
		{"underline", NewStyle().Reverse(true).Underline(true)},
		{"strikethrough", NewStyle().Reverse(true).Strikethrough(true)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.style.Render("a b")
			parts := strings.Split(got, " ")
			if len(parts) != 2 {
				t.Fatalf("Render(%q) = %q, want one space", "a b", got)
			}
			spaceStyleStart := strings.LastIndex(parts[0], "\x1b[")
			if spaceStyleStart < 0 {
				t.Fatalf("Render(%q) = %q, the space has no style", "a b", got)
			}
			if spaceStyle := parts[0][spaceStyleStart:]; !strings.Contains(spaceStyle, "7") {
				t.Fatalf("space style %q does not enable reverse: %q", spaceStyle, got)
			}
		})
	}
}
