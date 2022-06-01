package lipgloss

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestSetColorProfile(t *testing.T) {
	t.Parallel()

	style := NewStyle().Foreground(Color("#5A56E0"))
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
			SetColorProfile(tc.profile)
			res := style.Render(input)
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
