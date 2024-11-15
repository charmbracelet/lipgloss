package lipgloss

import (
	"testing"
)

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
		r, g, b, _ := Color(tc.input).RGBA()
		o := uint(r>>8)<<16 + uint(g>>8)<<8 + uint(b>>8)
		if o != tc.expected {
			t.Errorf("expected %X, got %X (test #%d)", tc.expected, o, i+1)
		}
	}
}

func TestRGBA(t *testing.T) {
	tt := []struct {
		input    string
		expected uint
	}{
		// lipgloss.Color
		{
			"#FF0000",
			0xFF0000,
		},
		{
			"9",
			9,
		},
		{
			"21",
			21,
		},
	}

	for i, tc := range tt {
		r, g, b, _ := Color(tc.input).RGBA()
		o := uint(r/256)<<16 + uint(g/256)<<8 + uint(b/256)

		if o != tc.expected {
			t.Errorf("expected %X, got %X (test #%d)", tc.expected, o, i+1)
		}
	}
}
