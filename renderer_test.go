package lipgloss

import (
	"os"
	"testing"
)

func TestRendererHasDarkBackground(t *testing.T) {
	r1 := NewRenderer(DetectColorProfile(os.Stdout, nil), true)
	r1.SetHasDarkBackground(false)
	if r1.HasDarkBackground() {
		t.Error("Expected renderer to have light background")
	}
	r2 := NewRenderer(DetectColorProfile(os.Stdout, nil), false)
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
	r := NewRenderer(DetectColorProfile(f, nil), true)
	r.SetColorProfile(TrueColor)
	if r.ColorProfile() != TrueColor {
		t.Error("Expected renderer to use true color")
	}
}

func TestRace(t *testing.T) {
	r := NewRenderer(DetectColorProfile(os.Stdout, nil), true)

	for i := 0; i < 100; i++ {
		t.Run("SetColorProfile", func(t *testing.T) {
			t.Parallel()
			r.SetHasDarkBackground(false)
			r.HasDarkBackground()
			r.SetColorProfile(ANSI256)
			r.SetHasDarkBackground(true)
		})
	}
}
