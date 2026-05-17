package compat

import (
	"image/color"
	"testing"
)

func TestAdaptiveColorUsesLightWhenBackgroundIsLight(t *testing.T) {
	SetHasDarkBackground(false)

	light := color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	dark := color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	c := AdaptiveColor{Light: light, Dark: dark}

	r, _, _, _ := c.RGBA()
	wantR, _, _, _ := light.RGBA()
	if r != wantR {
		t.Fatalf("expected light color channel %d, got %d", wantR, r)
	}
}

func TestAdaptiveColorUsesDarkWhenBackgroundIsDark(t *testing.T) {
	SetHasDarkBackground(true)

	light := color.RGBA{R: 0xff, G: 0x00, B: 0x00, A: 0xff}
	dark := color.RGBA{R: 0x00, G: 0x00, B: 0xff, A: 0xff}
	c := AdaptiveColor{Light: light, Dark: dark}

	_, _, b, _ := c.RGBA()
	_, _, wantB, _ := dark.RGBA()
	if b != wantB {
		t.Fatalf("expected dark color blue channel %d, got %d", wantB, b)
	}
}
