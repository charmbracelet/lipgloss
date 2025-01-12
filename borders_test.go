package lipgloss

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
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

func TestBorderStyle(t *testing.T) {
	tests := []struct {
		name     string
		title    string
		expected string
	}{
		{
			name:  "standard case",
			title: "Test",
			expected: strings.TrimSpace(`
┌─ Test ───┐
│          │
│          │
│          │
│          │
└──────────┘
`),
		},
		{
			name:  "ignores title if does not fit",
			title: "Title is too long a string and exceeds width",
			expected: strings.TrimSpace(`
┌──────────┐
│          │
│          │
│          │
│          │
└──────────┘
`),
		},
		{
			name:  "works with ansi escapes",
			title: NewStyle().Foreground(Color("#0ff")).Render("Test"),
			expected: strings.TrimSpace(`
┌─ Test ───┐
│          │
│          │
│          │
│          │
└──────────┘
`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := NewStyle().
				Width(10).
				Height(4).
				Border(NormalBorder()).
				BorderTitle(tt.title).
				Render()

			actual = ansi.Strip(actual)

			if actual != tt.expected {
				t.Errorf("expected:\n%s\n but got:\n%s", tt.expected, actual)
			}
		})
	}
}
