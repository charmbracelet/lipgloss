package lipgloss

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestSetColorProfile(t *testing.T) {
	tt := []struct {
		profile  termenv.Profile
		input    string
		style    Style
		expected string
	}{
		{
			termenv.Ascii,
			"hello",
			NewStyle().Foreground(Color("#5A56E0")),
			"hello",
		},
		{
			termenv.ANSI,
			"hello",
			NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[94mhello\x1b[0m",
		},
		{
			termenv.ANSI256,
			"hello",
			NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;5;62mhello\x1b[0m",
		},
		{
			termenv.TrueColor,
			"hello",
			NewStyle().Foreground(Color("#5A56E0")),
			"\x1b[38;2;89;86;224mhello\x1b[0m",
		},
	}

	for i, tc := range tt {
		SetColorProfile(tc.profile)
		res := tc.style.Render(tc.input)
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
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
			t.Errorf("expected %X, got %X (test #%d)", tc.expected, o, i+1)
		}
	}
}

func TestRGBA(t *testing.T) {
	tt := []struct {
		profile  termenv.Profile
		darkBg   bool
		input    TerminalColor
		expected uint
	}{
		// lipgloss.Color
		{
			termenv.TrueColor,
			true,
			Color("#FF0000"),
			0xFF0000,
		},
		{
			termenv.TrueColor,
			true,
			Color("9"),
			0xFF0000,
		},
		{
			termenv.TrueColor,
			true,
			Color("21"),
			0x0000FF,
		},
		// lipgloss.AdaptiveColor
		{
			termenv.TrueColor,
			true,
			AdaptiveColor{Dark: "#FF0000", Light: "#0000FF"},
			0xFF0000,
		},
		{
			termenv.TrueColor,
			false,
			AdaptiveColor{Dark: "#FF0000", Light: "#0000FF"},
			0x0000FF,
		},
		{
			termenv.TrueColor,
			true,
			AdaptiveColor{Dark: "9", Light: "21"},
			0xFF0000,
		},
		{
			termenv.TrueColor,
			false,
			AdaptiveColor{Dark: "9", Light: "21"},
			0x0000FF,
		},
		// lipgloss.CompleteColor
		{
			termenv.TrueColor,
			true,
			CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
			0xFF0000,
		},
		{
			termenv.ANSI256,
			true,
			CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
			0xFFFFFF,
		},
		{
			termenv.ANSI,
			true,
			CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
			0x0000FF,
		},
		// lipgloss.CompleteAdaptiveColor
		// dark
		{
			termenv.TrueColor,
			true,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#0000FF", ANSI256: "231", ANSI: "12"},
			},
			0xFF0000,
		},
		{
			termenv.ANSI256,
			true,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "21", ANSI: "12"},
			},
			0xFFFFFF,
		},
		{
			termenv.ANSI,
			true,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "9"},
			},
			0x0000FF,
		},
		// light
		{
			termenv.TrueColor,
			false,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#0000FF", ANSI256: "231", ANSI: "12"},
			},
			0x0000FF,
		},
		{
			termenv.ANSI256,
			false,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "21", ANSI: "12"},
			},
			0x0000FF,
		},
		{
			termenv.ANSI,
			false,
			CompleteAdaptiveColor{
				Dark:  CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "12"},
				Light: CompleteColor{TrueColor: "#FF0000", ANSI256: "231", ANSI: "9"},
			},
			0xFF0000,
		},
	}

	for i, tc := range tt {
		SetColorProfile(tc.profile)
		SetHasDarkBackground(tc.darkBg)

		r, g, b, _ := tc.input.RGBA()
		o := uint(r/256)<<16 + uint(g/256)<<8 + uint(b/256)

		if o != tc.expected {
			t.Errorf("expected %X, got %X (test #%d)", tc.expected, o, i+1)
		}
	}
}
