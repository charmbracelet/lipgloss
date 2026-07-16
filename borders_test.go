package lipgloss

import (
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/rivo/uniseg"
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

// Old implementation using rune slice conversion
func getFirstRuneAsStringOld(str string) string {
	if str == "" {
		return str
	}
	r := []rune(str)
	return string(r[0])
}

func TestGetFirstRuneAsString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"Empty", "", ""},
		{"SingleASCII", "A", "A"},
		{"SingleUnicode", "世", "世"},
		{"ASCIIString", "Hello", "H"},
		{"UnicodeString", "你好世界", "你"},
		{"MixedASCIIFirst", "Hello世界", "H"},
		{"MixedUnicodeFirst", "世界Hello", "世"},
		{"Emoji", "😀Happy", "😀"},
		{"MultiByteFirst", "ñoño", "ñ"},
		{"LongString", "The quick brown fox jumps over the lazy dog", "T"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getFirstRuneAsString(tt.input)
			if got != tt.want {
				t.Errorf("getFirstRuneAsString(%q) = %q, want %q", tt.input, got, tt.want)
			}

			// Verify new implementation matches old implementation
			old := getFirstRuneAsStringOld(tt.input)
			if got != old {
				t.Errorf("getFirstRuneAsString(%q) = %q, but old implementation returns %q", tt.input, got, old)
			}
		})
	}
}

func BenchmarkGetFirstRuneAsString(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"ASCII", "Hello, World!"},
		{"Unicode", "你好世界"},
		{"Single", "A"},
		{"Empty", ""},
	}

	b.Run("Old", func(b *testing.B) {
		for _, tc := range testCases {
			b.Run(tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = getFirstRuneAsStringOld(tc.str)
				}
			})
		}
	})

	b.Run("New", func(b *testing.B) {
		for _, tc := range testCases {
			b.Run(tc.name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = getFirstRuneAsString(tc.str)
				}
			})
		}
	})
}

func BenchmarkMaxRuneWidth(b *testing.B) {
	testCases := []struct {
		name string
		str  string
	}{
		{"Blank", " "},
		{"ASCII", "+"},
		{"Markdown", "|"},
		{"Normal", "├"},
		{"Rounded", "╭"},
		{"Block", "█"},
		{"Emoji", "😀"},
	}
	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			b.Run("Before", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					_ = maxRuneWidthOld(tc.str)
				}
			})
			b.Run("After", func(b *testing.B) {
				b.ReportAllocs()
				for b.Loop() {
					_ = maxRuneWidth(tc.str)
				}
			})
		})
	}
}

func maxRuneWidthOld(str string) int {
	var width int

	state := -1
	for len(str) > 0 {
		var w int
		_, str, w, state = uniseg.FirstGraphemeClusterInString(str, state)
		if w > width {
			width = w
		}
	}

	return width
}

func TestRenderHorizontalEdgeWideMiddle(t *testing.T) {
	tests := []struct {
		name      string
		left, mid string
		right     string
		width     int
	}{
		{"WideMid_EvenInner", "X", "あ", "Y", 6}, // inner 4: あ あ -> 4 cells
		{"WideMid_OddInner", "X", "あ", "Y", 7},  // inner 5: must not overshoot to 6
		{"WideMid_NoCorners", "", "あ", "", 5},   // inner 5
		{"NarrowMid_Odd", "X", "-", "Y", 7},     // inner 5, baseline
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := renderHorizontalEdge(tt.left, tt.mid, tt.right, tt.width)
			if w := ansi.StringWidth(got); w != tt.width {
				t.Errorf("renderHorizontalEdge(%q,%q,%q,%d) width = %d, want %d (out=%q)",
					tt.left, tt.mid, tt.right, tt.width, w, tt.width, got)
			}
		})
	}
}
