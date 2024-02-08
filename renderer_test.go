package lipgloss

import (
	"testing"
)

func TestRendererHasDarkBackground(t *testing.T) {
	r1 := NewRenderer()
	r1.SetHasDarkBackground(false)
	if r1.HasDarkBackground() {
		t.Error("Expected renderer to have light background")
	}
	r2 := NewRenderer()
	r2.SetHasDarkBackground(true)
	if !r2.HasDarkBackground() {
		t.Error("Expected renderer to have dark background")
	}
}
