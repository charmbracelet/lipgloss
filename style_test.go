package lipgloss

import (
	"testing"
)

func BenchmarkStyleRender(b *testing.B) {
	s := NewStyle().
		Bold(true).
		Foreground(Color("#ffffff"))

	for i := 0; i < b.N; i++ {
		s.Render("Hello world")
	}
}
