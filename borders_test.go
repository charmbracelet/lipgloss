package lipgloss

import (
	"testing"

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

func TestBorderGetTopBottomSize(t *testing.T) {
	tests := []struct {
		name       string
		border     Border
		wantTop    int
		wantBottom int
		wantLeft   int
		wantRight  int
	}{
		{
			name:       "normal border",
			border:     NormalBorder(),
			wantTop:    1,
			wantBottom: 1,
			wantLeft:   1,
			wantRight:  1,
		},
		{
			name:       "no border",
			border:     Border{},
			wantTop:    0,
			wantBottom: 0,
			wantLeft:   0,
			wantRight:  0,
		},
		{
			name: "wide rune in top-right corner",
			border: Border{
				Top: "─", Bottom: "─", Left: "│", Right: "│",
				TopLeft: "┌", TopRight: "⏩", BottomLeft: "└", BottomRight: "┘",
			},
			wantTop:    1, // top border is always 1 row, regardless of rune width
			wantBottom: 1,
			wantLeft:   1,
			wantRight:  2, // wide corner rune included in right edge width
		},
		{
			name: "wide rune in bottom-right corner",
			border: Border{
				Top: "─", Bottom: "─", Left: "│", Right: "│",
				TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "⏩",
			},
			wantTop:    1,
			wantBottom: 1, // bottom border is always 1 row
			wantLeft:   1,
			wantRight:  2, // wide corner rune included in right edge width
		},
		{
			name: "wide rune as left border char",
			border: Border{
				Top: "─", Bottom: "─", Left: "⏩", Right: "│",
				TopLeft: "┌", TopRight: "┐", BottomLeft: "└", BottomRight: "┘",
			},
			wantTop:    1,
			wantBottom: 1,
			wantLeft:   2, // left border with wide rune is 2 cells wide
			wantRight:  1,
		},
		{
			name: "top-only border parts",
			border: Border{
				Top: "─", TopLeft: "┌", TopRight: "┐",
			},
			wantTop:    1,
			wantBottom: 0, // no bottom border parts
			wantLeft:   1, // TopLeft included in left edge width
			wantRight:  1, // TopRight included in right edge width
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.border.GetTopSize(); got != tt.wantTop {
				t.Errorf("GetTopSize() = %d, want %d", got, tt.wantTop)
			}
			if got := tt.border.GetBottomSize(); got != tt.wantBottom {
				t.Errorf("GetBottomSize() = %d, want %d", got, tt.wantBottom)
			}
			if got := tt.border.GetLeftSize(); got != tt.wantLeft {
				t.Errorf("GetLeftSize() = %d, want %d", got, tt.wantLeft)
			}
			if got := tt.border.GetRightSize(); got != tt.wantRight {
				t.Errorf("GetRightSize() = %d, want %d", got, tt.wantRight)
			}
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
