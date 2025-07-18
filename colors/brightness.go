package colors

import (
	"image/color"
)

// Darken takes a color and makes it darker by a specific percentage (0-100, clamped).
func Darken(c color.Color, percent int) color.Color {
	if c == nil {
		return nil
	}

	mult := 1.0 - clamp(float64(percent), 0, 100)/100.0

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(float64(r>>8) * mult),
		G: uint8(float64(g>>8) * mult),
		B: uint8(float64(b>>8) * mult),
		A: uint8(min(255, float64(a>>8))),
	}
}

// Lighten makes a color lighter by a specific percentage (0-100, clamped).
func Lighten(c color.Color, percent int) color.Color {
	if c == nil {
		return nil
	}

	add := 255 * (clamp(float64(percent), 0, 100) / 100.0)

	r, g, b, a := c.RGBA()
	return color.RGBA{
		R: uint8(min(255, float64(r>>8)+add)),
		G: uint8(min(255, float64(g>>8)+add)),
		B: uint8(min(255, float64(b>>8)+add)),
		A: uint8(min(255, float64(a>>8))),
	}
}
