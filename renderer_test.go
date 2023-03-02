package lipgloss

import (
	"os"
	"testing"

	"github.com/muesli/termenv"
)

func TestRendererHasDarkBackground(t *testing.T) {
	r1 := NewRenderer(os.Stdout)
	r1.SetHasDarkBackground(false)
	if r1.HasDarkBackground() {
		t.Error("Expected renderer to have light background")
	}
	r2 := NewRenderer(os.Stdout)
	r2.SetHasDarkBackground(true)
	if !r2.HasDarkBackground() {
		t.Error("Expected renderer to have dark background")
	}
}

func TestRendererWithOutput(t *testing.T) {
	f, err := os.Create(t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	r := NewRenderer(f)
	r.SetColorProfile(termenv.TrueColor)
	if r.output.Profile != termenv.TrueColor {
		t.Error("Expected renderer to use true color")
	}
}
