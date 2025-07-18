package colors

import (
	"fmt"
	"image/color"
	"testing"
)

// hex converts a color to a hex string or panics if invalid.
func hex(hex string) color.Color {
	cf, err := FromHex(hex)
	if err != nil {
		panic(err)
	}
	return cf
}

func expectColorMatches(t *testing.T, got, want color.Color) {
	t.Helper()

	if (got == nil) != (want == nil) {
		t.Errorf("expectColorMatches() = %s, want %s", rgbaString(t, got), rgbaString(t, want))
	}

	if got == nil {
		return
	}

	gr, gg, gb, ga := got.RGBA()
	wr, wg, wb, wa := want.RGBA()

	gru, ggu, gbu, gau := uint8(gr>>8), uint8(gg>>8), uint8(gb>>8), uint8(ga>>8)
	wru, wgu, wbu, wau := uint8(wr>>8), uint8(wg>>8), uint8(wb>>8), uint8(wa>>8)

	if gru != wru || ggu != wgu || gbu != wbu || gau != wau {
		t.Errorf("expectColorMatches() = %s, want %s", rgbaString(t, got), rgbaString(t, want))
	}
}

func rgbaString(t *testing.T, c color.Color) string {
	t.Helper()

	if c == nil {
		return "nil"
	}

	r, g, b, a := c.RGBA()
	return fmt.Sprintf("rgba(%d,%d,%d,%d)", uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

func TestAlpha(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		alpha    float64
		expected color.Color
	}{
		{
			name:     "alpha-full-opacity",
			color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
			alpha:    1.0,
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:     "alpha-half-opacity",
			color:    color.RGBA{R: 0, G: 255, B: 0, A: 255},
			alpha:    0.5,
			expected: color.RGBA{R: 0, G: 255, B: 0, A: 127},
		},
		{
			name:     "alpha-quarter-opacity",
			color:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			alpha:    0.25,
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 63},
		},
		{
			name:     "alpha-zero-opacity",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			alpha:    0.0,
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 0},
		},
		{
			name:     "alpha-clamp-above-max",
			color:    color.RGBA{R: 255, G: 0, B: 255, A: 255},
			alpha:    1.5,
			expected: color.RGBA{R: 255, G: 0, B: 255, A: 255},
		},
		{
			name:     "alpha-clamp-below-min",
			color:    color.RGBA{R: 255, G: 255, B: 0, A: 255},
			alpha:    -0.5,
			expected: color.RGBA{R: 255, G: 255, B: 0, A: 0},
		},
		{
			name:     "alpha-complex-color",
			color:    color.RGBA{R: 18, G: 52, B: 86, A: 255},
			alpha:    0.75,
			expected: color.RGBA{R: 18, G: 52, B: 86, A: 191},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Alpha(tt.color, tt.alpha), tt.expected)
		})
	}
}

func TestComplementary(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		expected color.Color
	}{
		{name: "red", color: hex("#FF0000"), expected: hex("#00FFFF")},
		{name: "green", color: hex("#00FF00"), expected: hex("#FF00FF")},
		{name: "blue", color: hex("#0000FF"), expected: hex("#FFFF00")},
		{name: "yellow", color: hex("#FFFF00"), expected: hex("#0000FF")},
		{name: "cyan", color: hex("#00FFFF"), expected: hex("#FF0000")},
		{name: "magenta", color: hex("#FF00FF"), expected: hex("#00FF00")},
		// Black has no hue to complement
		{name: "black", color: hex("#000000"), expected: hex("#000000")},
		// White has no hue to complement
		{name: "white", color: hex("#FFFFFF"), expected: hex("#FFFFFF")},
		// Gray has no hue to complement
		{name: "gray", color: hex("#808080"), expected: hex("#808080")},
		{name: "orange", color: hex("#FF8000"), expected: hex("#007FFF")},
		{name: "purple", color: hex("#8000FF"), expected: hex("#7FFF00")},
		{name: "nil-color", color: nil, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Complementary(tt.color), tt.expected)
		})
	}
}

func TestFromHex(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    color.Color
		expectError bool
	}{
		{name: "valid-6-red", input: "#FF0000", expected: hex("#FF0000")},
		{name: "valid-6-green", input: "#00FF00", expected: hex("#00FF00")},
		{name: "valid-6-blue", input: "#0000FF", expected: hex("#0000FF")},
		{name: "valid-6-white", input: "#FFFFFF", expected: hex("#FFFFFF")},
		{name: "valid-6-black", input: "#000000", expected: hex("#000000")},
		{name: "valid-6-gray", input: "#808080", expected: hex("#808080")},
		{name: "valid-3-red", input: "#F00", expected: hex("#FF0000")},
		{name: "valid-3-green", input: "#0F0", expected: hex("#00FF00")},
		{name: "valid-3-blue", input: "#00F", expected: hex("#0000FF")},
		{name: "valid-3-white", input: "#FFF", expected: hex("#FFFFFF")},
		{name: "valid-3-black", input: "#000", expected: hex("#000000")},
		{name: "valid-6-lowercase", input: "#ff0000", expected: hex("#FF0000")},
		{name: "valid-6-mixed-case", input: "#Ff0000", expected: hex("#FF0000")},
		{name: "valid-3-lowercase", input: "#f00", expected: hex("#FF0000")},
		{name: "missing-hash-prefix", input: "FF0000", expectError: true},
		{name: "empty-string", input: "", expectError: true},
		{name: "only-hash", input: "#", expectError: true},
		{name: "too-short-3", input: "#F0", expectError: true},
		{name: "too-long-6", input: "#FF00000", expectError: true},
		{name: "invalid-char", input: "#FG0000", expectError: true},
		{name: "invalid-char-3", input: "#FG0", expectError: true},
		{name: "invalid-char-lowercase", input: "#fg0000", expectError: true},
		{name: "invalid-char-mixed", input: "#Fg0000", expectError: true},
		{name: "wrong-len-5", input: "#FF000", expectError: true},
		{name: "wrong-len-8", input: "#FF000000", expectError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := FromHex(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("FromHex() expected error but got none for input %q", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("FromHex() unexpected error for input %q: %v", tt.input, err)
				return
			}

			expectColorMatches(t, result, tt.expected)
		})
	}
}
