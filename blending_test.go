package lipgloss

import (
	"image/color"
	"testing"
)

func TestBlend1D(t *testing.T) {
	tests := []struct {
		name     string
		steps    int
		stops    []color.Color
		expected []color.Color
	}{
		{
			name:  "2-colors-10-steps",
			steps: 10,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 246, G: 0, B: 45, A: 255},
				&color.RGBA{R: 235, G: 0, B: 73, A: 255},
				&color.RGBA{R: 223, G: 0, B: 99, A: 255},
				&color.RGBA{R: 210, G: 0, B: 124, A: 255},
				&color.RGBA{R: 193, G: 0, B: 149, A: 255},
				&color.RGBA{R: 173, G: 0, B: 175, A: 255},
				&color.RGBA{R: 147, G: 0, B: 201, A: 255},
				&color.RGBA{R: 109, G: 0, B: 228, A: 255},
				&color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
		},
		{
			name:  "3-colors-4-steps",
			steps: 4,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 0, G: 255, B: 0, A: 255},
				&color.RGBA{R: 0, G: 255, B: 0, A: 255},
				&color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
		},
		{
			name:  "black-to-white-5-steps",
			steps: 5,
			stops: []color.Color{
				color.RGBA{R: 0, G: 0, B: 0, A: 255},
				color.RGBA{R: 255, G: 255, B: 255, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 0, G: 0, B: 0, A: 255},
				&color.RGBA{R: 59, G: 59, B: 59, A: 255},
				&color.RGBA{R: 119, G: 119, B: 119, A: 255},
				&color.RGBA{R: 185, G: 185, B: 185, A: 255},
				&color.RGBA{R: 255, G: 255, B: 255, A: 255},
			},
		},
		{
			name:  "4-colors-6-steps",
			steps: 6,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 255, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 255, G: 255, B: 0, A: 255},
				&color.RGBA{R: 255, G: 255, B: 0, A: 255},
				&color.RGBA{R: 0, G: 255, B: 0, A: 255},
				&color.RGBA{R: 0, G: 255, B: 0, A: 255},
				&color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
		},
		{
			name:  "2-steps-5-stops",
			steps: 2,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
				color.RGBA{R: 255, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 0, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 0, G: 255, B: 0, A: 255},
			},
		},
		{
			name:  "insufficient-stops",
			steps: 3,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
		},
		{
			name:  "insufficient-steps",
			steps: 1,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expected: []color.Color{
				&color.RGBA{R: 255, G: 0, B: 0, A: 255},
				&color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Blend1D(tt.steps, tt.stops...)

			if len(got) != len(tt.expected) {
				t.Errorf("Blend() = %v length, want %v length", len(got), len(tt.expected))
			}

			for i := range tt.expected {
				expectColorMatches(t, got[i], tt.expected[i])
			}
		})
	}
}

func TestBlend2D(t *testing.T) {
	tests := []struct {
		name           string
		width, height  int
		angle          float64
		stops          []color.Color
		expectedLength int
	}{
		{
			name:   "2x2-red-to-blue-0deg",
			width:  2,
			height: 2,
			angle:  0,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 4,
		},
		{
			name:   "3x2-red-to-blue-90deg",
			width:  3,
			height: 2,
			angle:  90,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 6,
		},
		{
			name:   "2x3-red-to-blue-180deg",
			width:  2,
			height: 3,
			angle:  180,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 6,
		},
		{
			name:   "2x2-red-to-blue-270deg",
			width:  2,
			height: 2,
			angle:  270,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 4,
		},
		{
			name:   "1x1-single-color",
			width:  1,
			height: 1,
			angle:  0,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
			expectedLength: 1,
		},
		{
			name:   "3-colors-2x2-0deg",
			width:  2,
			height: 2,
			angle:  0,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 255, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 4,
		},
		{
			name:   "invalid-dimensions-fallback",
			width:  0,
			height: -1,
			angle:  0,
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
			},
			expectedLength: 1,
		},
		{
			name:   "angle-normalization-450",
			width:  2,
			height: 2,
			angle:  450, // Should normalize to 90
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 4,
		},
		{
			name:   "negative-angle-normalization",
			width:  2,
			height: 2,
			angle:  -90, // Should normalize to 270
			stops: []color.Color{
				color.RGBA{R: 255, G: 0, B: 0, A: 255},
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
			},
			expectedLength: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Blend2D(tt.width, tt.height, tt.angle, tt.stops...)

			if len(got) != tt.expectedLength {
				t.Errorf("Blend2D() = %v length, want %v length", len(got), tt.expectedLength)
			}

			// Verify row-major order by checking that the width matches.
			if tt.width > 0 && tt.height > 0 {
				expectedTotal := max(tt.width, 1) * max(tt.height, 1)
				if len(got) != expectedTotal {
					t.Errorf("Blend2D() total pixels = %v, want %v", len(got), expectedTotal)
				}
			}

			// Verify that we have valid colors (not nil).
			for i, color := range got {
				if color == nil {
					t.Errorf("Blend2D() color at index %d is nil", i)
				}
			}

			// For single color tests, verify all colors are the same.
			if len(tt.stops) == 1 && len(got) > 0 {
				firstColor := got[0]
				for _, color := range got {
					expectColorMatches(t, color, firstColor)
				}
			}
		})
	}
}

func TestBlend2DEdgeCases(t *testing.T) {
	t.Run("nil-stops", func(t *testing.T) {
		t.Parallel()
		got := Blend2D(2, 2, 0, nil, nil)
		if got != nil {
			t.Errorf("Blend2D() with nil stops = %v, want nil", got)
		}
	})

	t.Run("empty-stops", func(t *testing.T) {
		t.Parallel()
		got := Blend2D(2, 2, 0)
		if got != nil {
			t.Errorf("Blend2D() with empty stops = %v, want nil", got)
		}
	})

	t.Run("nil-color-in-stops", func(t *testing.T) {
		t.Parallel()
		got := Blend2D(2, 2, 0, color.RGBA{R: 255, G: 0, B: 0, A: 255}, nil, color.RGBA{R: 0, G: 0, B: 255, A: 255})
		if len(got) != 4 {
			t.Errorf("Blend2D() with nil color in stops = %v length, want 4", len(got))
		}
		// Should still work with the non-nil colors and produce valid colors
		for i, color := range got {
			if color == nil {
				t.Errorf("Blend2D() color at index %d is nil", i)
			}
		}
	})
}

func BenchmarkBlend1D(b *testing.B) {
	stops := []color.Color{
		hex("#FF0000"), // Red
		hex("#00FF00"), // Green
		hex("#0000FF"), // Blue
		hex("#FFFF00"), // Yellow
		hex("#FF00FF"), // Magenta
	}

	for b.Loop() {
		Blend1D(100, stops...)
	}
}

func BenchmarkBlend2D(b *testing.B) {
	stops := []color.Color{
		hex("#FF0000"), // Red
		hex("#00FF00"), // Green
		hex("#0000FF"), // Blue
		hex("#FFFF00"), // Yellow
		hex("#FF00FF"), // Magenta
	}

	for b.Loop() {
		Blend2D(100, 50, 45, stops...)
	}
}
