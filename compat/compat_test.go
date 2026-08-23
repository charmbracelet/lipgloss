package compat

import (
	"image/color"
	"testing"
)

func TestAdaptiveColorRGBA(t *testing.T) {
	ac := AdaptiveColor{
		Light: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		Dark:  color.RGBA{R: 0, G: 0, B: 255, A: 255},
	}
	r, g, b, a := ac.RGBA()
	if r == 0 && g == 0 && b == 0 && a == 0 {
		t.Error("RGBA returned zero for a valid AdaptiveColor")
	}
}

func TestCompleteColorRGBA(t *testing.T) {
	cc := CompleteColor{
		TrueColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		ANSI256:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
		ANSI:      color.RGBA{R: 128, G: 128, B: 128, A: 255},
	}
	r, g, b, a := cc.RGBA()
	if r == 0 && g == 0 && b == 0 && a == 0 {
		t.Error("RGBA returned zero for a valid CompleteColor")
	}
}

func TestCompleteAdaptiveColorRGBA(t *testing.T) {
	cac := CompleteAdaptiveColor{
		Light: CompleteColor{
			TrueColor: color.RGBA{R: 255, G: 255, B: 0, A: 255},
			ANSI256:   color.RGBA{R: 200, G: 200, B: 0, A: 255},
			ANSI:      color.RGBA{R: 128, G: 128, B: 0, A: 255},
		},
		Dark: CompleteColor{
			TrueColor: color.RGBA{R: 0, G: 255, B: 255, A: 255},
			ANSI256:   color.RGBA{R: 0, G: 200, B: 200, A: 255},
			ANSI:      color.RGBA{R: 0, G: 128, B: 128, A: 255},
		},
	}
	r, g, b, a := cac.RGBA()
	if r == 0 && g == 0 && b == 0 && a == 0 {
		t.Error("RGBA returned zero for a valid CompleteAdaptiveColor")
	}
}
