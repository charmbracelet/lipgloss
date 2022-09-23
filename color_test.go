package lipgloss

import (
	"image/color"
	"testing"

	"github.com/muesli/termenv"
)

func TestSetColorProfile(t *testing.T) {
	r := renderer
	input := "hello"

	tt := []struct {
		name     string
		profile  termenv.Profile
		expected string
	}{
		{
			"ascii",
			termenv.Ascii,
			"hello",
		},
		{
			"ansi",
			termenv.ANSI,
			"\x1b[94mhello\x1b[0m",
		},
		{
			"ansi256",
			termenv.ANSI256,
			"\x1b[38;5;62mhello\x1b[0m",
		},
		{
			"truecolor",
			termenv.TrueColor,
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r.SetColorProfile(tc.profile)
			style := NewStyle().Foreground(Color("#5A56E0"))
			res := Render(style, input)

			if res != tc.expected {
				t.Errorf("Expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
					tc.expected, formatEscapes(tc.expected),
					res, formatEscapes(res))
			}
		})
	}
}

func TestHexToColor(t *testing.T) {
	t.Parallel()

	tt := []struct {
		input    string
		expected uint
	}{
		{
			"#FF0000",
			0xFF0000,
		},
		{
			"#00F",
			0x0000FF,
		},
		{
			"#6B50FF",
			0x6B50FF,
		},
		{
			"invalid color",
			0x0,
		},
	}

	for i, tc := range tt {
		h := hexToColor(tc.input)
		o := uint(h.R)<<16 + uint(h.G)<<8 + uint(h.B)
		if o != tc.expected {
			t.Errorf("expected %X, got %X (test #%d)", o, tc.expected, i+1)
		}
	}
}

// hexToColor translates a hex color string (#RRGGBB or #RGB) into a color.RGB,
// which satisfies the color.Color interface. If an invalid string is passed
// black with 100% opacity will be returned: or, in hex format, 0x000000FF.
func hexToColor(hex string) (c color.RGBA) {
	c.A = 0xFF

	if hex == "" || hex[0] != '#' {
		return c
	}

	const (
		fullFormat  = 7 // #RRGGBB
		shortFormat = 4 // #RGB
	)

	switch len(hex) {
	case fullFormat:
		const offset = 4
		c.R = hexToByte(hex[1])<<offset + hexToByte(hex[2])
		c.G = hexToByte(hex[3])<<offset + hexToByte(hex[4])
		c.B = hexToByte(hex[5])<<offset + hexToByte(hex[6])
	case shortFormat:
		const offset = 0x11
		c.R = hexToByte(hex[1]) * offset
		c.G = hexToByte(hex[2]) * offset
		c.B = hexToByte(hex[3]) * offset
	}

	return c
}

func hexToByte(b byte) byte {
	const offset = 10
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + offset
	case b >= 'A' && b <= 'F':
		return b - 'A' + offset
	}
	// Invalid, but just return 0.
	return 0
}
