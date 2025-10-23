package lipgloss

import (
	"testing"
)

func BenchmarkBorderRendering(b *testing.B) {
	dimensions := []struct {
		name   string
		width  int
		height int
	}{
		{"10x5", 10, 5},
		{"20x10", 20, 10},
		{"40x20", 40, 15},
		{"80x40", 80, 20},
		{"120x60", 120, 25},
		{"160x80", 160, 30},
	}

	for _, dim := range dimensions {
		b.Run(dim.name, func(b *testing.B) {
			style := NewStyle().
				Border(RoundedBorder(), true).
				Foreground(Color("#ffffff")).
				Background(Color("#000000")).
				Width(dim.width).
				Height(dim.height)

			b.ResetTimer()
			for b.Loop() {
				_ = style.Render("")
			}
		})
	}
}

func BenchmarkBorderBlend(b *testing.B) {
	dimensions := []struct {
		name   string
		width  int
		height int
	}{
		{"10x5", 10, 5},
		{"20x10", 20, 10},
		{"40x20", 40, 15},
		{"80x40", 80, 20},
		{"120x60", 120, 25},
		{"160x80", 160, 30},
	}

	for _, dim := range dimensions {
		b.Run(dim.name, func(b *testing.B) {
			style := NewStyle().
				Border(RoundedBorder(), true).
				BorderForegroundBlend(
					Color("#00FA68"),
					Color("#9900FF"),
					Color("#ED5353"),
				).
				Width(dim.width).
				Height(dim.height)

			b.ResetTimer()
			for b.Loop() {
				_ = style.Render("")
			}
		})
	}
}

func BenchmarkBorderRenderingNoColors(b *testing.B) {
	dimensions := []struct {
		name   string
		width  int
		height int
	}{
		{"10x5", 10, 5},
		{"20x10", 20, 10},
		{"40x20", 40, 15},
		{"80x40", 80, 20},
		{"120x60", 120, 25},
		{"160x80", 160, 30},
	}

	for _, dim := range dimensions {
		b.Run(dim.name, func(b *testing.B) {
			style := NewStyle().
				Border(RoundedBorder(), true).
				Width(dim.width).
				Height(dim.height)

			b.ResetTimer()
			for b.Loop() {
				_ = style.Render("")
			}
		})
	}
}
