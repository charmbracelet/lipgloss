package colors

import (
	"fmt"
	"image/color"
	"testing"

	"github.com/lucasb-eyer/go-colorful"
)

// fromHex converts a color to a hex string. Could use lipgloss.Color() but in the future,
// this could cause a circular dependency.
func fromHex(hex string) color.Color {
	cf, err := colorful.Hex(hex)
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
		{
			name:     "complementary-red",
			color:    fromHex("#FF0000"), // Red
			expected: fromHex("#00FFFF"), // Cyan
		},
		{
			name:     "complementary-green",
			color:    fromHex("#00FF00"), // Green
			expected: fromHex("#FF00FF"), // Magenta
		},
		{
			name:     "complementary-blue",
			color:    fromHex("#0000FF"), // Blue
			expected: fromHex("#FFFF00"), // Yellow
		},
		{
			name:     "complementary-yellow",
			color:    fromHex("#FFFF00"), // Yellow
			expected: fromHex("#0000FF"), // Blue
		},
		{
			name:     "complementary-cyan",
			color:    fromHex("#00FFFF"), // Cyan
			expected: fromHex("#FF0000"), // Red
		},
		{
			name:     "complementary-magenta",
			color:    fromHex("#FF00FF"), // Magenta
			expected: fromHex("#00FF00"), // Green
		},
		{
			name:     "complementary-black",
			color:    fromHex("#000000"), // Black
			expected: fromHex("#000000"), // Black (achromatic, no hue to complement)
		},
		{
			name:     "complementary-white",
			color:    fromHex("#FFFFFF"), // White
			expected: fromHex("#FFFFFF"), // White (achromatic, no hue to complement)
		},
		{
			name:     "complementary-gray",
			color:    fromHex("#808080"), // Gray
			expected: fromHex("#808080"), // Gray (complementary of gray is gray)
		},
		{
			name:     "complementary-orange",
			color:    fromHex("#FF8000"), // Orange
			expected: fromHex("#007FFF"), // Blue-cyan
		},
		{
			name:     "complementary-purple",
			color:    fromHex("#8000FF"), // Purple
			expected: fromHex("#7FFF00"), // Lime green
		},
		{
			name:     "complementary-nil-color",
			color:    nil,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Complementary(tt.color), tt.expected)
		})
	}
}
