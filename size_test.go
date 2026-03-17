package lipgloss

import (
	"strings"
	"testing"
)

func BenchmarkWidthSimple(b *testing.B) {
	simpleStrings := []string{
		"ab",
		"abcdef",
		"abcdefghij",
	}

	for _, str := range simpleStrings {
		b.Run("len-"+str, func(b *testing.B) {
			for b.Loop() {
				_ = Width(str)
			}
		})
	}
}

func BenchmarkWidthMultiLine(b *testing.B) {
	multiLineStrings := []struct {
		name string
		str  string
	}{
		{"2-lines", "Line 1\nLine 2"},
		{"10-lines", "Line 1\nLine 2\nLine 3\nLine 4\nLine 5\nLine 6\nLine 7\nLine 8\nLine 9\nLine 10"},
		{"50-lines", strings.Repeat("Line\n", 49) + "Line"},
	}

	for _, tc := range multiLineStrings {
		b.Run(tc.name, func(b *testing.B) {
			for b.Loop() {
				_ = Width(tc.str)
			}
		})
	}
}
