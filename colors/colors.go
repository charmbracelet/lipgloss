// Package colors provides a set of color-related utilities for working with colors.
// This includes utilities for blending colors, adjusting brightness/alpha, etc.
package colors

import (
	"cmp"
	"errors"
	"image/color"

	"github.com/lucasb-eyer/go-colorful"
)

func clamp[T cmp.Ordered](v, low, high T) T {
	return min(high, max(low, v))
}

// ensureNotTransparent ensures that the alpha value of a color is not 0, and if
// it is, we will set it to 1. This is useful for when we are converting from
// RGB -> RGBA, and the alpha value is lost in the conversion for gradient purposes.
func ensureNotTransparent(c color.Color) color.Color {
	_, _, _, a := c.RGBA()
	if a == 0 {
		return Alpha(c, 1)
	}
	return c
}

// Alpha adjusts the alpha value of a color using a 0-1 (clamped) float scale
// 0 = transparent, 1 = opaque.
func Alpha(c color.Color, alpha float64) color.Color {
	if c == nil {
		return nil
	}

	r, g, b, _ := c.RGBA()
	return color.RGBA{
		R: uint8(min(255, float64(r>>8))),
		G: uint8(min(255, float64(g>>8))),
		B: uint8(min(255, float64(b>>8))),
		A: uint8(clamp(alpha, 0, 1) * 255),
	}
}

// Complementary returns the complementary color (180° away on color wheel) of
// the given color. This is useful for creating a contrasting color.
func Complementary(c color.Color) color.Color {
	if c == nil {
		return nil
	}

	// Offset hue by 180°.
	cf, _ := colorful.MakeColor(ensureNotTransparent(c))

	h, s, v := cf.Hsv()
	h += 180
	if h >= 360 {
		h -= 360
	} else if h < 0 {
		h += 360
	}

	return colorful.Hsv(h, s, v).Clamped()
}

var errInvalidFormat = errors.New("invalid hex format")

// FromHex parses a hex color string and returns a color.RGBA. The string can be
// in the format #RRGGBB or #RGB. This is a more performant implementation of
// [colorful.Hex].
func FromHex(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if len(s) == 0 || s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return c, err
}
