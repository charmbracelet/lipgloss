package colors

import (
	"image/color"
	"testing"
)

func TestDarken(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		percent  int
		expected color.Color
	}{
		{
			name:     "darken-white-50-percent",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			percent:  50,
			expected: color.RGBA{R: 127, G: 127, B: 127, A: 255},
		},
		{
			name:     "darken-red-25-percent",
			color:    color.RGBA{R: 255, G: 0, B: 0, A: 255},
			percent:  25,
			expected: color.RGBA{R: 191, G: 0, B: 0, A: 255},
		},
		{
			name:     "darken-blue-75-percent",
			color:    color.RGBA{R: 0, G: 0, B: 255, A: 255},
			percent:  75,
			expected: color.RGBA{R: 0, G: 0, B: 63, A: 255},
		},
		{
			name:     "darken-black-10-percent",
			color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			percent:  10,
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "darken-with-clamp-min",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			percent:  0,
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "darken-with-clamp-max",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			percent:  100,
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "darken-nil-color",
			color:    nil,
			percent:  50,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Darken(tt.color, tt.percent), tt.expected)
		})
	}
}

func TestLighten(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		percent  int
		expected color.Color
	}{
		{
			name:     "lighten-black-50-percent",
			color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			percent:  50,
			expected: color.RGBA{R: 127, G: 127, B: 127, A: 255},
		},
		{
			name:     "lighten-red-25-percent",
			color:    color.RGBA{R: 128, G: 0, B: 0, A: 255},
			percent:  25,
			expected: color.RGBA{R: 191, G: 63, B: 63, A: 255},
		},
		{
			name:     "lighten-blue-75-percent",
			color:    color.RGBA{R: 0, G: 0, B: 128, A: 255},
			percent:  75,
			expected: color.RGBA{R: 191, G: 191, B: 255, A: 255},
		},
		{
			name:     "lighten-white-10-percent",
			color:    color.RGBA{R: 255, G: 255, B: 255, A: 255},
			percent:  10,
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "lighten-with-clamp-min",
			color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			percent:  0,
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "lighten-with-clamp-max",
			color:    color.RGBA{R: 0, G: 0, B: 0, A: 255},
			percent:  100,
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "lighten-nil-color",
			color:    nil,
			percent:  50,
			expected: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Lighten(tt.color, tt.percent), tt.expected)
		})
	}
}
