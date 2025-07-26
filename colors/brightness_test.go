package colors

import (
	"image/color"
	"testing"
)

func TestDarken(t *testing.T) {
	tests := []struct {
		name     string
		color    color.Color
		percent  float64
		expected color.Color
	}{
		{name: "darken-white-50-percent", color: hex("#FFFFFF"), percent: 0.5, expected: hex("#7F7F7F")},
		{name: "darken-red-25-percent", color: hex("#FF0000"), percent: 0.25, expected: hex("#BF0000")},
		{name: "darken-blue-75-percent", color: hex("#0000FF"), percent: 0.75, expected: hex("#00003F")},
		{name: "darken-black-10-percent", color: hex("#000000"), percent: 0.1, expected: hex("#000000")},
		{name: "darken-with-clamp-min", color: hex("#FFFFFF"), percent: 0, expected: hex("#FFFFFF")},
		{name: "darken-with-clamp-max", color: hex("#FFFFFF"), percent: 1, expected: hex("#000000")},
		{name: "darken-nil-color", color: nil, percent: 0.5, expected: nil},
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
		percent  float64
		expected color.Color
	}{
		{name: "lighten-black-50-percent", color: hex("#000000"), percent: 0.5, expected: hex("#7F7F7F")},
		{name: "lighten-red-25-percent", color: hex("#800000"), percent: 0.25, expected: hex("#BF3F3F")},
		{name: "lighten-blue-75-percent", color: hex("#000080"), percent: 0.75, expected: hex("#BFBFFF")},
		{name: "lighten-white-10-percent", color: hex("#FFFFFF"), percent: 0.1, expected: hex("#FFFFFF")},
		{name: "lighten-with-clamp-min", color: hex("#000000"), percent: 0, expected: hex("#000000")},
		{name: "lighten-with-clamp-max", color: hex("#000000"), percent: 1, expected: hex("#FFFFFF")},
		{name: "lighten-nil-color", color: nil, percent: 0.5, expected: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			expectColorMatches(t, Lighten(tt.color, tt.percent), tt.expected)
		})
	}
}
