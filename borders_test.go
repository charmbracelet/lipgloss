package lipgloss

import (
	"testing"

	"github.com/rivo/uniseg"
)

func TestStyle_GetBorderSizes(t *testing.T) {
	tests := []struct {
		name  string
		style Style
		wantX int
		wantY int
	}{
		{
			name:  "Default style",
			style: NewStyle(),
			wantX: 0,
			wantY: 0,
		},
		{
			name:  "Border(NormalBorder())",
			style: NewStyle().Border(NormalBorder()),
			wantX: 2,
			wantY: 2,
		},
		{
			name:  "Border(NormalBorder(), true)",
			style: NewStyle().Border(NormalBorder(), true),
			wantX: 2,
			wantY: 2,
		},
		{
			name:  "Border(NormalBorder(), true, false)",
			style: NewStyle().Border(NormalBorder(), true, false),
			wantX: 0,
			wantY: 2,
		},
		{
			name:  "Border(NormalBorder(), true, true, false)",
			style: NewStyle().Border(NormalBorder(), true, true, false),
			wantX: 2,
			wantY: 1,
		},
		{
			name:  "Border(NormalBorder(), true, true, false, false)",
			style: NewStyle().Border(NormalBorder(), true, true, false, false),
			wantX: 1,
			wantY: 1,
		},
		{
			name:  "BorderTop(true).BorderStyle(NormalBorder())",
			style: NewStyle().BorderTop(true).BorderStyle(NormalBorder()),
			wantX: 0,
			wantY: 1,
		},
		{
			name:  "BorderStyle(NormalBorder())",
			style: NewStyle().BorderStyle(NormalBorder()),
			wantX: 2,
			wantY: 2,
		},
		{
			name:  "Custom BorderStyle",
			style: NewStyle().BorderStyle(Border{Left: "123456789"}),
			wantX: 1, // left and right borders are laid out vertically, one rune per row
			wantY: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX := tt.style.GetHorizontalBorderSize()
			if gotX != tt.wantX {
				t.Errorf("Style.GetHorizontalBorderSize() got %d, want %d", gotX, tt.wantX)
			}

			gotY := tt.style.GetVerticalBorderSize()
			if gotY != tt.wantY {
				t.Errorf("Style.GetVerticalBorderSize() got %d, want %d", gotY, tt.wantY)
			}

			gotX = tt.style.GetHorizontalFrameSize()
			if gotX != tt.wantX {
				t.Errorf("Style.GetHorizontalFrameSize() got %d, want %d", gotX, tt.wantX)
			}

			gotY = tt.style.GetVerticalFrameSize()
			if gotY != tt.wantY {
				t.Errorf("Style.GetVerticalFrameSize() got %d, want %d", gotY, tt.wantY)
			}

			gotX, gotY = tt.style.GetFrameSize()
			if gotX != tt.wantX || gotY != tt.wantY {
				t.Errorf("Style.GetFrameSize() got (%d, %d), want (%d, %d)", gotX, gotY, tt.wantX, tt.wantY)
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

func TestGetBorderSidesWithImplicitBorders(t *testing.T) {
	// When using BorderStyle() without explicitly setting sides,
	// all sides should be implicitly enabled (matching render behavior).
	// See: https://github.com/charmbracelet/lipgloss/issues/522
	t.Run("BorderStyle only (implicit borders)", func(t *testing.T) {
		s := NewStyle().BorderStyle(NormalBorder())

		if !s.GetBorderTop() {
			t.Error("GetBorderTop() = false, want true (implicit border)")
		}
		if !s.GetBorderRight() {
			t.Error("GetBorderRight() = false, want true (implicit border)")
		}
		if !s.GetBorderBottom() {
			t.Error("GetBorderBottom() = false, want true (implicit border)")
		}
		if !s.GetBorderLeft() {
			t.Error("GetBorderLeft() = false, want true (implicit border)")
		}

		_, top, right, bottom, left := s.GetBorder()
		if !top || !right || !bottom || !left {
			t.Errorf("GetBorder() sides = (%v, %v, %v, %v), want all true", top, right, bottom, left)
		}
	})

	// When sides are explicitly set, implicit borders should not interfere.
	t.Run("Border with explicit sides", func(t *testing.T) {
		s := NewStyle().Border(NormalBorder(), true, false, true, false)

		if !s.GetBorderTop() {
			t.Error("GetBorderTop() = false, want true")
		}
		if s.GetBorderRight() {
			t.Error("GetBorderRight() = true, want false")
		}
		if !s.GetBorderBottom() {
			t.Error("GetBorderBottom() = false, want true")
		}
		if s.GetBorderLeft() {
			t.Error("GetBorderLeft() = true, want false")
		}
	})

	// When no border is set at all, all getters should return false.
	t.Run("No border set", func(t *testing.T) {
		s := NewStyle()

		if s.GetBorderTop() {
			t.Error("GetBorderTop() = true, want false")
		}
		if s.GetBorderRight() {
			t.Error("GetBorderRight() = true, want false")
		}
		if s.GetBorderBottom() {
			t.Error("GetBorderBottom() = true, want false")
		}
		if s.GetBorderLeft() {
			t.Error("GetBorderLeft() = true, want false")
		}
	})

	// Frame size should be consistent with border side getters.
	t.Run("Frame size consistency with BorderStyle only", func(t *testing.T) {
		s := NewStyle().BorderStyle(NormalBorder()).Padding(1, 2)

		// GetHorizontalFrameSize should include both border and padding
		hFrame := s.GetHorizontalFrameSize()
		wantH := 2 + 4 // 2 for left+right border, 4 for left+right padding
		if hFrame != wantH {
			t.Errorf("GetHorizontalFrameSize() = %d, want %d", hFrame, wantH)
		}

		vFrame := s.GetVerticalFrameSize()
		wantV := 2 + 2 // 2 for top+bottom border, 2 for top+bottom padding
		if vFrame != wantV {
			t.Errorf("GetVerticalFrameSize() = %d, want %d", vFrame, wantV)
		}
	})

	// Once a side is explicitly set, implicitBorders is no longer in effect.
	t.Run("BorderStyle then explicit side disables implicit", func(t *testing.T) {
		s := NewStyle().BorderStyle(NormalBorder()).BorderTop(true)

		// Top was explicitly set to true
		if !s.GetBorderTop() {
			t.Error("GetBorderTop() = false, want true")
		}
		// Other sides should be false because at least one side was explicitly set,
		// so implicit borders no longer applies.
		if s.GetBorderRight() {
			t.Error("GetBorderRight() = true, want false (implicit borders disabled)")
		}
		if s.GetBorderBottom() {
			t.Error("GetBorderBottom() = true, want false (implicit borders disabled)")
		}
		if s.GetBorderLeft() {
			t.Error("GetBorderLeft() = true, want false (implicit borders disabled)")
		}
	})
}
