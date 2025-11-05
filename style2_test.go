package lipgloss

import (
	"testing"
)

func TestBlockRender(t *testing.T) {
	cases := []struct {
		name     string
		block    Style2
		content  string
		expected string
	}{
		{
			name:     "empty block",
			block:    Style2{},
			expected: "",
		},
		{
			name: "block with size",
			block: NewStyle2().Width(5).
				Height(3),
			expected: NewStyle().Width(5).Height(3).Render(""),
		},
		{
			name: "block with content",
			block: NewStyle2().Width(7).
				Height(3),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(3).Render("Hi"),
		},
		{
			name: "block with border",
			block: NewStyle2().Width(7).
				Height(3).
				Border(RoundedBorder()),
			content:  "Hi",
			expected: NewStyle().Border(RoundedBorder()).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "padding uniform",
			block: NewStyle2().Width(7).
				Height(5).
				Padding(1),
			content:  "Hi",
			expected: NewStyle().Padding(1).Width(7).Height(5).Render("Hi"),
		},
		{
			name: "padding asymmetric",
			block: NewStyle2().Width(9).
				Height(5).
				Padding(1, 2, 1, 1),
			content:  "Hi",
			expected: NewStyle().Padding(1, 2, 1, 1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "border top and bottom",
			block: NewStyle2().Width(7).
				Height(3).
				Border(NormalBorder(), true, false),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), true, false).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "border left and right",
			block: NewStyle2().Width(7).
				Height(3).
				Border(NormalBorder(), false, true),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), false, true).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "zero size explicit",
			block: NewStyle2().Width(0).
				Height(0),
			content:  "",
			expected: "",
		},
		{
			name: "margin uniform",
			block: NewStyle2().Width(7).
				Height(3).
				Margin(1),
			content:  "Hi",
			expected: NewStyle().Margin(1).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "margin asymmetric",
			block: NewStyle2().Width(9).
				Height(5).
				Margin(1, 2, 1, 1),
			content:  "Hi",
			expected: NewStyle().Margin(1, 2, 1, 1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "margin with border rounded",
			block: NewStyle2().Width(9).
				Height(5).
				Margin(1).
				Border(RoundedBorder()),
			content:  "Hi",
			expected: NewStyle().Border(RoundedBorder()).Margin(1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "margin with border normal left+right",
			block: NewStyle2().Width(9).
				Height(5).
				Margin(1).
				Border(NormalBorder(), false, true),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), false, true).Margin(1).Width(9).Height(5).Render("Hi"),
		},
		{
			name:     "style with colors (content only)",
			block:    NewStyle2().Width(7).Height(3),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(3).Render("Hi"),
		},
		{
			name: "borders with per-side colors",
			block: NewStyle2().Width(9).Height(5).
				Border(RoundedBorder()).
				BorderTopForeground(Red).
				BorderRightForeground(Green).
				BorderBottomForeground(Yellow).
				BorderLeftForeground(Blue),
			content: "Hi",
			expected: NewStyle().Border(RoundedBorder()).
				BorderTopForeground(Red).
				BorderRightForeground(Green).
				BorderBottomForeground(Yellow).
				BorderLeftForeground(Blue).
				Width(9).Height(5).Render("Hi"),
		},
		{
			name: "margin + border + padding + colors",
			block: NewStyle2().Width(11).Height(7).
				Margin(1).
				Padding(1).
				Border(RoundedBorder()).
				BorderForeground(Red, Green, Yellow, Blue),
			content: "Hi",
			expected: NewStyle().Margin(1).Padding(1).Border(RoundedBorder()).
				BorderForeground(Red, Green, Yellow, Blue).Width(11).Height(7).Render("Hi"),
		},
		{
			name: "inner too small: skip content, draw borders",
			block: NewStyle2().Width(2).Height(2).
				Border(NormalBorder()),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder()).Width(2).Height(2).Render("Hi"),
		},
		// Max constraints
		{
			name:     "max width wraps/clips content",
			block:    NewStyle2().MaxWidth(5),
			content:  "HelloWorld",
			expected: NewStyle().MaxWidth(5).Render("HelloWorld"),
		},
		{
			name:     "max height clips lines",
			block:    NewStyle2().MaxHeight(2),
			content:  "A\nB\nC",
			expected: NewStyle().MaxHeight(2).Render("A\nB\nC"),
		},
		// Horizontal alignment
		{
			name:     "align horizontal left",
			block:    NewStyle2().Width(9).Height(3).AlignHorizontal(Left),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Left).Render("Hi"),
		},
		{
			name:     "align horizontal center",
			block:    NewStyle2().Width(9).Height(3).AlignHorizontal(Center),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Center).Render("Hi"),
		},
		{
			name:     "align horizontal right",
			block:    NewStyle2().Width(9).Height(3).AlignHorizontal(Right),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Right).Render("Hi"),
		},
		// Vertical alignment
		{
			name:     "align vertical top",
			block:    NewStyle2().Width(7).Height(5).AlignVertical(Top),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Top).Render("Hi"),
		},
		{
			name:     "align vertical center",
			block:    NewStyle2().Width(7).Height(5).AlignVertical(Center),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Center).Render("Hi"),
		},
		{
			name:     "align vertical bottom",
			block:    NewStyle2().Width(7).Height(5).AlignVertical(Bottom),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Bottom).Render("Hi"),
		},
		// Combined border + padding + center align
		{
			name:     "border + padding + center alignment",
			block:    NewStyle2().Width(11).Height(7).Border(RoundedBorder()).Padding(1).Align(Center, Center),
			content:  "Hi",
			expected: NewStyle().Border(RoundedBorder()).Padding(1).Width(11).Height(7).Align(Center, Center).Render("Hi"),
		},
	}

	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.block.Render(tc.content)
			if result != tc.expected {
				t.Errorf("case %d (%s): expected:\n%q\nbut got:\n%q", i, tc.name, tc.expected, result)
			}
		})
	}
}

func TestBlockSettersAndGetters(t *testing.T) {
	t.Run("Background", func(t *testing.T) {
		b := NewStyle2().Background(Red)
		if b.GetBackground() != Red {
			t.Errorf("expected Red, got %v", b.GetBackground())
		}
	})

	t.Run("Foreground", func(t *testing.T) {
		b := NewStyle2().Foreground(Blue)
		if b.GetForeground() != Blue {
			t.Errorf("expected Blue, got %v", b.GetForeground())
		}
	})

	t.Run("Width", func(t *testing.T) {
		b := NewStyle2().Width(42)
		if b.GetWidth() != 42 {
			t.Errorf("expected 42, got %d", b.GetWidth())
		}
	})

	t.Run("Height", func(t *testing.T) {
		b := NewStyle2().Height(24)
		if b.GetHeight() != 24 {
			t.Errorf("expected 24, got %d", b.GetHeight())
		}
	})

	t.Run("MaxWidth", func(t *testing.T) {
		b := NewStyle2().MaxWidth(100)
		if b.GetMaxWidth() != 100 {
			t.Errorf("expected 100, got %d", b.GetMaxWidth())
		}
	})

	t.Run("MaxHeight", func(t *testing.T) {
		b := NewStyle2().MaxHeight(50)
		if b.GetMaxHeight() != 50 {
			t.Errorf("expected 50, got %d", b.GetMaxHeight())
		}
	})

	t.Run("Align", func(t *testing.T) {
		b := NewStyle2().Align(Center, Bottom)
		h, v := b.GetAlign()
		if h != Center || v != Bottom {
			t.Errorf("expected (Center, Bottom), got (%v, %v)", h, v)
		}
	})

	t.Run("AlignHorizontal", func(t *testing.T) {
		b := NewStyle2().AlignHorizontal(Right)
		if b.GetAlignHorizontal() != Right {
			t.Errorf("expected Right, got %v", b.GetAlignHorizontal())
		}
	})

	t.Run("AlignVertical", func(t *testing.T) {
		b := NewStyle2().AlignVertical(Top)
		if b.GetAlignVertical() != Top {
			t.Errorf("expected Top, got %v", b.GetAlignVertical())
		}
	})

	t.Run("Padding", func(t *testing.T) {
		b := NewStyle2().Padding(1, 2, 3, 4)
		top, right, bottom, left := b.GetPadding()
		if top != 1 || right != 2 || bottom != 3 || left != 4 {
			t.Errorf("expected (1, 2, 3, 4), got (%d, %d, %d, %d)", top, right, bottom, left)
		}
	})

	t.Run("PaddingLeft", func(t *testing.T) {
		b := NewStyle2().PaddingLeft(5)
		if b.GetPaddingLeft() != 5 {
			t.Errorf("expected 5, got %d", b.GetPaddingLeft())
		}
	})

	t.Run("PaddingRight", func(t *testing.T) {
		b := NewStyle2().PaddingRight(6)
		if b.GetPaddingRight() != 6 {
			t.Errorf("expected 6, got %d", b.GetPaddingRight())
		}
	})

	t.Run("PaddingTop", func(t *testing.T) {
		b := NewStyle2().PaddingTop(7)
		if b.GetPaddingTop() != 7 {
			t.Errorf("expected 7, got %d", b.GetPaddingTop())
		}
	})

	t.Run("PaddingBottom", func(t *testing.T) {
		b := NewStyle2().PaddingBottom(8)
		if b.GetPaddingBottom() != 8 {
			t.Errorf("expected 8, got %d", b.GetPaddingBottom())
		}
	})

	t.Run("PaddingChar", func(t *testing.T) {
		b := NewStyle2().PaddingChar('*')
		if b.GetPaddingChar() != '*' {
			t.Errorf("expected '*', got %c", b.GetPaddingChar())
		}
	})

	t.Run("Margin", func(t *testing.T) {
		b := NewStyle2().Margin(1, 2, 3, 4)
		top, right, bottom, left := b.GetMargin()
		if top != 1 || right != 2 || bottom != 3 || left != 4 {
			t.Errorf("expected (1, 2, 3, 4), got (%d, %d, %d, %d)", top, right, bottom, left)
		}
	})

	t.Run("MarginLeft", func(t *testing.T) {
		b := NewStyle2().MarginLeft(9)
		if b.GetMarginLeft() != 9 {
			t.Errorf("expected 9, got %d", b.GetMarginLeft())
		}
	})

	t.Run("MarginRight", func(t *testing.T) {
		b := NewStyle2().MarginRight(10)
		if b.GetMarginRight() != 10 {
			t.Errorf("expected 10, got %d", b.GetMarginRight())
		}
	})

	t.Run("MarginTop", func(t *testing.T) {
		b := NewStyle2().MarginTop(11)
		if b.GetMarginTop() != 11 {
			t.Errorf("expected 11, got %d", b.GetMarginTop())
		}
	})

	t.Run("MarginBottom", func(t *testing.T) {
		b := NewStyle2().MarginBottom(12)
		if b.GetMarginBottom() != 12 {
			t.Errorf("expected 12, got %d", b.GetMarginBottom())
		}
	})

	t.Run("MarginChar", func(t *testing.T) {
		b := NewStyle2().MarginChar('.')
		if b.GetMarginChar() != '.' {
			t.Errorf("expected '.', got %c", b.GetMarginChar())
		}
	})

	t.Run("MarginBackground", func(t *testing.T) {
		b := NewStyle2().MarginBackground(Green)
		if b.GetMarginBackground() != Green {
			t.Errorf("expected Green, got %v", b.GetMarginBackground())
		}
	})

	t.Run("BorderStyle", func(t *testing.T) {
		b := NewStyle2().BorderStyle(RoundedBorder())
		border := b.GetBorderStyle()
		expected := RoundedBorder()
		if border != expected {
			t.Errorf("expected RoundedBorder, got %v", border)
		}
	})

	t.Run("GetBorder", func(t *testing.T) {
		b := NewStyle2().Border(NormalBorder(), true, false, true, false)
		border, top, right, bottom, left := b.GetBorder()
		expected := NormalBorder()
		if border != expected {
			t.Errorf("expected NormalBorder, got %v", border)
		}
		if !top || right || !bottom || left {
			t.Errorf("expected (true, false, true, false), got (%v, %v, %v, %v)", top, right, bottom, left)
		}
	})

	t.Run("BorderTop", func(t *testing.T) {
		b := NewStyle2().BorderTop(true)
		if !b.GetBorderTop() {
			t.Error("expected true, got false")
		}
	})

	t.Run("BorderRight", func(t *testing.T) {
		b := NewStyle2().BorderRight(true)
		if !b.GetBorderRight() {
			t.Error("expected true, got false")
		}
	})

	t.Run("BorderBottom", func(t *testing.T) {
		b := NewStyle2().BorderBottom(true)
		if !b.GetBorderBottom() {
			t.Error("expected true, got false")
		}
	})

	t.Run("BorderLeft", func(t *testing.T) {
		b := NewStyle2().BorderLeft(true)
		if !b.GetBorderLeft() {
			t.Error("expected true, got false")
		}
	})

	t.Run("BorderTopForeground", func(t *testing.T) {
		b := NewStyle2().BorderTopForeground(Red)
		if b.GetBorderTopForeground() != Red {
			t.Errorf("expected Red, got %v", b.GetBorderTopForeground())
		}
	})

	t.Run("BorderRightForeground", func(t *testing.T) {
		b := NewStyle2().BorderRightForeground(Green)
		if b.GetBorderRightForeground() != Green {
			t.Errorf("expected Green, got %v", b.GetBorderRightForeground())
		}
	})

	t.Run("BorderBottomForeground", func(t *testing.T) {
		b := NewStyle2().BorderBottomForeground(Yellow)
		if b.GetBorderBottomForeground() != Yellow {
			t.Errorf("expected Yellow, got %v", b.GetBorderBottomForeground())
		}
	})

	t.Run("BorderLeftForeground", func(t *testing.T) {
		b := NewStyle2().BorderLeftForeground(Blue)
		if b.GetBorderLeftForeground() != Blue {
			t.Errorf("expected Blue, got %v", b.GetBorderLeftForeground())
		}
	})

	t.Run("BorderTopBackground", func(t *testing.T) {
		b := NewStyle2().BorderTopBackground(Red)
		if b.GetBorderTopBackground() != Red {
			t.Errorf("expected Red, got %v", b.GetBorderTopBackground())
		}
	})

	t.Run("BorderRightBackground", func(t *testing.T) {
		b := NewStyle2().BorderRightBackground(Green)
		if b.GetBorderRightBackground() != Green {
			t.Errorf("expected Green, got %v", b.GetBorderRightBackground())
		}
	})

	t.Run("BorderBottomBackground", func(t *testing.T) {
		b := NewStyle2().BorderBottomBackground(Yellow)
		if b.GetBorderBottomBackground() != Yellow {
			t.Errorf("expected Yellow, got %v", b.GetBorderBottomBackground())
		}
	})

	t.Run("BorderLeftBackground", func(t *testing.T) {
		b := NewStyle2().BorderLeftBackground(Blue)
		if b.GetBorderLeftBackground() != Blue {
			t.Errorf("expected Blue, got %v", b.GetBorderLeftBackground())
		}
	})

	t.Run("GetHorizontalMargins", func(t *testing.T) {
		b := NewStyle2().MarginLeft(3).MarginRight(5)
		if b.GetHorizontalMargins() != 8 {
			t.Errorf("expected 8, got %d", b.GetHorizontalMargins())
		}
	})

	t.Run("GetVerticalMargins", func(t *testing.T) {
		b := NewStyle2().MarginTop(2).MarginBottom(4)
		if b.GetVerticalMargins() != 6 {
			t.Errorf("expected 6, got %d", b.GetVerticalMargins())
		}
	})

	t.Run("GetHorizontalPadding", func(t *testing.T) {
		b := NewStyle2().PaddingLeft(1).PaddingRight(3)
		if b.GetHorizontalPadding() != 4 {
			t.Errorf("expected 4, got %d", b.GetHorizontalPadding())
		}
	})

	t.Run("GetVerticalPadding", func(t *testing.T) {
		b := NewStyle2().PaddingTop(2).PaddingBottom(2)
		if b.GetVerticalPadding() != 4 {
			t.Errorf("expected 4, got %d", b.GetVerticalPadding())
		}
	})
}

func TestBlockUnsetters(t *testing.T) {
	t.Run("UnsetBackground", func(t *testing.T) {
		b := NewStyle2().Background(Red).UnsetBackground()
		if b.GetBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetBackground())
		}
	})

	t.Run("UnsetForeground", func(t *testing.T) {
		b := NewStyle2().Foreground(Blue).UnsetForeground()
		if b.GetForeground() != nil {
			t.Errorf("expected nil, got %v", b.GetForeground())
		}
	})

	t.Run("UnsetWidth", func(t *testing.T) {
		b := NewStyle2().Width(42).UnsetWidth()
		if b.GetWidth() != 0 {
			t.Errorf("expected 0, got %d", b.GetWidth())
		}
	})

	t.Run("UnsetHeight", func(t *testing.T) {
		b := NewStyle2().Height(24).UnsetHeight()
		if b.GetHeight() != 0 {
			t.Errorf("expected 0, got %d", b.GetHeight())
		}
	})

	t.Run("UnsetMaxWidth", func(t *testing.T) {
		b := NewStyle2().MaxWidth(100).UnsetMaxWidth()
		if b.GetMaxWidth() != 0 {
			t.Errorf("expected 0, got %d", b.GetMaxWidth())
		}
	})

	t.Run("UnsetMaxHeight", func(t *testing.T) {
		b := NewStyle2().MaxHeight(50).UnsetMaxHeight()
		if b.GetMaxHeight() != 0 {
			t.Errorf("expected 0, got %d", b.GetMaxHeight())
		}
	})

	t.Run("UnsetAlign", func(t *testing.T) {
		b := NewStyle2().Align(Center, Bottom).UnsetAlign()
		h, v := b.GetAlign()
		if h != Position(0) || v != Position(0) {
			t.Errorf("expected (0, 0), got (%v, %v)", h, v)
		}
	})

	t.Run("UnsetAlignHorizontal", func(t *testing.T) {
		b := NewStyle2().AlignHorizontal(Right).UnsetAlignHorizontal()
		if b.GetAlignHorizontal() != Position(0) {
			t.Errorf("expected 0, got %v", b.GetAlignHorizontal())
		}
	})

	t.Run("UnsetAlignVertical", func(t *testing.T) {
		b := NewStyle2().AlignVertical(Top).UnsetAlignVertical()
		if b.GetAlignVertical() != Position(0) {
			t.Errorf("expected 0, got %v", b.GetAlignVertical())
		}
	})

	t.Run("UnsetPadding", func(t *testing.T) {
		b := NewStyle2().Padding(1, 2, 3, 4).UnsetPadding()
		top, right, bottom, left := b.GetPadding()
		if top != 0 || right != 0 || bottom != 0 || left != 0 {
			t.Errorf("expected (0, 0, 0, 0), got (%d, %d, %d, %d)", top, right, bottom, left)
		}
	})

	t.Run("UnsetPaddingLeft", func(t *testing.T) {
		b := NewStyle2().PaddingLeft(5).UnsetPaddingLeft()
		if b.GetPaddingLeft() != 0 {
			t.Errorf("expected 0, got %d", b.GetPaddingLeft())
		}
	})

	t.Run("UnsetPaddingRight", func(t *testing.T) {
		b := NewStyle2().PaddingRight(6).UnsetPaddingRight()
		if b.GetPaddingRight() != 0 {
			t.Errorf("expected 0, got %d", b.GetPaddingRight())
		}
	})

	t.Run("UnsetPaddingTop", func(t *testing.T) {
		b := NewStyle2().PaddingTop(7).UnsetPaddingTop()
		if b.GetPaddingTop() != 0 {
			t.Errorf("expected 0, got %d", b.GetPaddingTop())
		}
	})

	t.Run("UnsetPaddingBottom", func(t *testing.T) {
		b := NewStyle2().PaddingBottom(8).UnsetPaddingBottom()
		if b.GetPaddingBottom() != 0 {
			t.Errorf("expected 0, got %d", b.GetPaddingBottom())
		}
	})

	t.Run("UnsetPaddingChar", func(t *testing.T) {
		b := NewStyle2().PaddingChar('*').UnsetPaddingChar()
		if b.GetPaddingChar() != 0 {
			t.Errorf("expected 0, got %c", b.GetPaddingChar())
		}
	})

	t.Run("UnsetMargins", func(t *testing.T) {
		b := NewStyle2().Margin(1, 2, 3, 4).UnsetMargins()
		top, right, bottom, left := b.GetMargin()
		if top != 0 || right != 0 || bottom != 0 || left != 0 {
			t.Errorf("expected (0, 0, 0, 0), got (%d, %d, %d, %d)", top, right, bottom, left)
		}
	})

	t.Run("UnsetMarginLeft", func(t *testing.T) {
		b := NewStyle2().MarginLeft(9).UnsetMarginLeft()
		if b.GetMarginLeft() != 0 {
			t.Errorf("expected 0, got %d", b.GetMarginLeft())
		}
	})

	t.Run("UnsetMarginRight", func(t *testing.T) {
		b := NewStyle2().MarginRight(10).UnsetMarginRight()
		if b.GetMarginRight() != 0 {
			t.Errorf("expected 0, got %d", b.GetMarginRight())
		}
	})

	t.Run("UnsetMarginTop", func(t *testing.T) {
		b := NewStyle2().MarginTop(11).UnsetMarginTop()
		if b.GetMarginTop() != 0 {
			t.Errorf("expected 0, got %d", b.GetMarginTop())
		}
	})

	t.Run("UnsetMarginBottom", func(t *testing.T) {
		b := NewStyle2().MarginBottom(12).UnsetMarginBottom()
		if b.GetMarginBottom() != 0 {
			t.Errorf("expected 0, got %d", b.GetMarginBottom())
		}
	})

	t.Run("UnsetMarginChar", func(t *testing.T) {
		b := NewStyle2().MarginChar('.').UnsetMarginChar()
		if b.GetMarginChar() != 0 {
			t.Errorf("expected 0, got %c", b.GetMarginChar())
		}
	})

	t.Run("UnsetMarginBackground", func(t *testing.T) {
		b := NewStyle2().MarginBackground(Green).UnsetMarginBackground()
		if b.GetMarginBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetMarginBackground())
		}
	})

	t.Run("UnsetBorderStyle", func(t *testing.T) {
		b := NewStyle2().BorderStyle(RoundedBorder()).UnsetBorderStyle()
		if b.GetBorderStyle() != (Border{}) {
			t.Errorf("expected empty Border, got %v", b.GetBorderStyle())
		}
	})

	t.Run("UnsetBorderTop", func(t *testing.T) {
		b := NewStyle2().BorderTop(true).UnsetBorderTop()
		if b.GetBorderTop() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetBorderRight", func(t *testing.T) {
		b := NewStyle2().BorderRight(true).UnsetBorderRight()
		if b.GetBorderRight() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetBorderBottom", func(t *testing.T) {
		b := NewStyle2().BorderBottom(true).UnsetBorderBottom()
		if b.GetBorderBottom() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetBorderLeft", func(t *testing.T) {
		b := NewStyle2().BorderLeft(true).UnsetBorderLeft()
		if b.GetBorderLeft() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetBorderTopForeground", func(t *testing.T) {
		b := NewStyle2().BorderTopForeground(Red).UnsetBorderTopForeground()
		if b.GetBorderTopForeground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderTopForeground())
		}
	})

	t.Run("UnsetBorderRightForeground", func(t *testing.T) {
		b := NewStyle2().BorderRightForeground(Green).UnsetBorderRightForeground()
		if b.GetBorderRightForeground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderRightForeground())
		}
	})

	t.Run("UnsetBorderBottomForeground", func(t *testing.T) {
		b := NewStyle2().BorderBottomForeground(Yellow).UnsetBorderBottomForeground()
		if b.GetBorderBottomForeground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderBottomForeground())
		}
	})

	t.Run("UnsetBorderLeftForeground", func(t *testing.T) {
		b := NewStyle2().BorderLeftForeground(Blue).UnsetBorderLeftForeground()
		if b.GetBorderLeftForeground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderLeftForeground())
		}
	})

	t.Run("UnsetBorderForeground", func(t *testing.T) {
		b := NewStyle2().
			BorderTopForeground(Red).
			BorderRightForeground(Green).
			BorderBottomForeground(Yellow).
			BorderLeftForeground(Blue).
			UnsetBorderForeground()
		if b.GetBorderTopForeground() != nil ||
			b.GetBorderRightForeground() != nil ||
			b.GetBorderBottomForeground() != nil ||
			b.GetBorderLeftForeground() != nil {
			t.Error("expected all border foregrounds to be nil")
		}
	})

	t.Run("UnsetBorderTopBackground", func(t *testing.T) {
		b := NewStyle2().BorderTopBackground(Red).UnsetBorderTopBackground()
		if b.GetBorderTopBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderTopBackground())
		}
	})

	t.Run("UnsetBorderRightBackground", func(t *testing.T) {
		b := NewStyle2().BorderRightBackground(Green).UnsetBorderRightBackground()
		if b.GetBorderRightBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderRightBackground())
		}
	})

	t.Run("UnsetBorderBottomBackground", func(t *testing.T) {
		b := NewStyle2().BorderBottomBackground(Yellow).UnsetBorderBottomBackground()
		if b.GetBorderBottomBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderBottomBackground())
		}
	})

	t.Run("UnsetBorderLeftBackground", func(t *testing.T) {
		b := NewStyle2().BorderLeftBackground(Blue).UnsetBorderLeftBackground()
		if b.GetBorderLeftBackground() != nil {
			t.Errorf("expected nil, got %v", b.GetBorderLeftBackground())
		}
	})

	t.Run("UnsetBorderBackground", func(t *testing.T) {
		b := NewStyle2().
			BorderTopBackground(Red).
			BorderRightBackground(Green).
			BorderBottomBackground(Yellow).
			BorderLeftBackground(Blue).
			UnsetBorderBackground()
		if b.GetBorderTopBackground() != nil ||
			b.GetBorderRightBackground() != nil ||
			b.GetBorderBottomBackground() != nil ||
			b.GetBorderLeftBackground() != nil {
			t.Error("expected all border backgrounds to be nil")
		}
	})
}

func TestBlockChaining(t *testing.T) {
	b := NewStyle2().
		Width(10).
		Height(5).
		Background(Red).
		Foreground(Blue).
		Padding(1).
		Margin(2).
		Border(RoundedBorder()).
		BorderForeground(Green).
		AlignHorizontal(Center)

	if b.GetWidth() != 10 {
		t.Errorf("expected width 10, got %d", b.GetWidth())
	}
	if b.GetHeight() != 5 {
		t.Errorf("expected height 5, got %d", b.GetHeight())
	}
	if b.GetBackground() != Red {
		t.Errorf("expected Red background, got %v", b.GetBackground())
	}
	if b.GetForeground() != Blue {
		t.Errorf("expected Blue foreground, got %v", b.GetForeground())
	}
	top, right, bottom, left := b.GetPadding()
	if top != 1 || right != 1 || bottom != 1 || left != 1 {
		t.Errorf("expected padding (1, 1, 1, 1), got (%d, %d, %d, %d)", top, right, bottom, left)
	}
	top, right, bottom, left = b.GetMargin()
	if top != 2 || right != 2 || bottom != 2 || left != 2 {
		t.Errorf("expected margin (2, 2, 2, 2), got (%d, %d, %d, %d)", top, right, bottom, left)
	}
	if b.GetAlignHorizontal() != Center {
		t.Errorf("expected Center alignment, got %v", b.GetAlignHorizontal())
	}
}

func TestStyle2TextAttributes(t *testing.T) {
	t.Run("Bold", func(t *testing.T) {
		b := NewStyle2().Bold(true)
		if !b.GetBold() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Italic", func(t *testing.T) {
		b := NewStyle2().Italic(true)
		if !b.GetItalic() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Underline", func(t *testing.T) {
		b := NewStyle2().Underline(true)
		if !b.GetUnderline() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Strikethrough", func(t *testing.T) {
		b := NewStyle2().Strikethrough(true)
		if !b.GetStrikethrough() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Reverse", func(t *testing.T) {
		b := NewStyle2().Reverse(true)
		if !b.GetReverse() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Blink", func(t *testing.T) {
		b := NewStyle2().Blink(true)
		if !b.GetBlink() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Faint", func(t *testing.T) {
		b := NewStyle2().Faint(true)
		if !b.GetFaint() {
			t.Error("expected true, got false")
		}
	})

	t.Run("UnderlineSpaces", func(t *testing.T) {
		b := NewStyle2().UnderlineSpaces(true)
		if !b.GetUnderlineSpaces() {
			t.Error("expected true, got false")
		}
	})

	t.Run("StrikethroughSpaces", func(t *testing.T) {
		b := NewStyle2().StrikethroughSpaces(true)
		if !b.GetStrikethroughSpaces() {
			t.Error("expected true, got false")
		}
	})

	t.Run("ColorWhitespace", func(t *testing.T) {
		b := NewStyle2().ColorWhitespace(true)
		if !b.GetColorWhitespace() {
			t.Error("expected true, got false")
		}
	})

	t.Run("Inline", func(t *testing.T) {
		b := NewStyle2().Inline(true)
		if !b.GetInline() {
			t.Error("expected true, got false")
		}
	})

	t.Run("TabWidth", func(t *testing.T) {
		b := NewStyle2().TabWidth(8)
		if b.GetTabWidth() != 8 {
			t.Errorf("expected 8, got %d", b.GetTabWidth())
		}
	})

	t.Run("Transform", func(t *testing.T) {
		fn := func(s string) string { return "transformed" }
		b := NewStyle2().Transform(fn)
		if b.GetTransform() == nil {
			t.Error("expected transform function, got nil")
		}
	})
}

func TestStyle2UnsetTextAttributes(t *testing.T) {
	t.Run("UnsetBold", func(t *testing.T) {
		b := NewStyle2().Bold(true).UnsetBold()
		if b.GetBold() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetItalic", func(t *testing.T) {
		b := NewStyle2().Italic(true).UnsetItalic()
		if b.GetItalic() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetUnderline", func(t *testing.T) {
		b := NewStyle2().Underline(true).UnsetUnderline()
		if b.GetUnderline() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetStrikethrough", func(t *testing.T) {
		b := NewStyle2().Strikethrough(true).UnsetStrikethrough()
		if b.GetStrikethrough() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetReverse", func(t *testing.T) {
		b := NewStyle2().Reverse(true).UnsetReverse()
		if b.GetReverse() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetBlink", func(t *testing.T) {
		b := NewStyle2().Blink(true).UnsetBlink()
		if b.GetBlink() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetFaint", func(t *testing.T) {
		b := NewStyle2().Faint(true).UnsetFaint()
		if b.GetFaint() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetUnderlineSpaces", func(t *testing.T) {
		b := NewStyle2().UnderlineSpaces(true).UnsetUnderlineSpaces()
		if b.GetUnderlineSpaces() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetStrikethroughSpaces", func(t *testing.T) {
		b := NewStyle2().StrikethroughSpaces(true).UnsetStrikethroughSpaces()
		if b.GetStrikethroughSpaces() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetColorWhitespace", func(t *testing.T) {
		b := NewStyle2().ColorWhitespace(true).UnsetColorWhitespace()
		if b.GetColorWhitespace() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetInline", func(t *testing.T) {
		b := NewStyle2().Inline(true).UnsetInline()
		if b.GetInline() {
			t.Error("expected false, got true")
		}
	})

	t.Run("UnsetTabWidth", func(t *testing.T) {
		b := NewStyle2().TabWidth(8).UnsetTabWidth()
		if b.GetTabWidth() != 0 {
			t.Errorf("expected 0, got %d", b.GetTabWidth())
		}
	})

	t.Run("UnsetTransform", func(t *testing.T) {
		fn := func(s string) string { return "transformed" }
		b := NewStyle2().Transform(fn).UnsetTransform()
		if b.GetTransform() != nil {
			t.Error("expected nil, got transform function")
		}
	})
}

// Compatibility tests comparing Style2 with Style
func TestStyle2CompatibilityTextAttributes(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "bold",
			style2:  NewStyle2().Bold(true),
			style:   NewStyle().Bold(true),
			content: "hello",
		},
		{
			name:    "italic",
			style2:  NewStyle2().Italic(true),
			style:   NewStyle().Italic(true),
			content: "hello",
		},
		{
			name:    "underline",
			style2:  NewStyle2().Underline(true),
			style:   NewStyle().Underline(true),
			content: "hello",
		},
		{
			name:    "strikethrough",
			style2:  NewStyle2().Strikethrough(true),
			style:   NewStyle().Strikethrough(true),
			content: "hello",
		},
		{
			name:    "reverse",
			style2:  NewStyle2().Reverse(true),
			style:   NewStyle().Reverse(true),
			content: "hello",
		},
		{
			name:    "blink",
			style2:  NewStyle2().Blink(true),
			style:   NewStyle().Blink(true),
			content: "hello",
		},
		{
			name:    "faint",
			style2:  NewStyle2().Faint(true),
			style:   NewStyle().Faint(true),
			content: "hello",
		},
		{
			name:    "underline with spaces",
			style2:  NewStyle2().Underline(true).UnderlineSpaces(true),
			style:   NewStyle().Underline(true).UnderlineSpaces(true),
			content: "ab c",
		},
		{
			name:    "underline without spaces",
			style2:  NewStyle2().Underline(true).UnderlineSpaces(false),
			style:   NewStyle().Underline(true).UnderlineSpaces(false),
			content: "ab c",
		},
		{
			name:    "strikethrough with spaces",
			style2:  NewStyle2().Strikethrough(true).StrikethroughSpaces(true),
			style:   NewStyle().Strikethrough(true).StrikethroughSpaces(true),
			content: "ab c",
		},
		{
			name:    "strikethrough without spaces",
			style2:  NewStyle2().Strikethrough(true).StrikethroughSpaces(false),
			style:   NewStyle().Strikethrough(true).StrikethroughSpaces(false),
			content: "ab c",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			t.Logf("Style output:\n%s\n\nStyle2 output:\n%s\n", result1, result2)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityColors(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "foreground color",
			style2:  NewStyle2().Foreground(Color("#5A56E0")),
			style:   NewStyle().Foreground(Color("#5A56E0")),
			content: "hello",
		},
		{
			name:    "background color",
			style2:  NewStyle2().Background(Color("#FF00FF")),
			style:   NewStyle().Background(Color("#FF00FF")),
			content: "hello",
		},
		{
			name:    "foreground and background",
			style2:  NewStyle2().Foreground(Color("#FFFFFF")).Background(Color("#000000")),
			style:   NewStyle().Foreground(Color("#FFFFFF")).Background(Color("#000000")),
			content: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityDimensions(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "width",
			style2:  NewStyle2().Width(10),
			style:   NewStyle().Width(10),
			content: "hello",
		},
		{
			name:    "height",
			style2:  NewStyle2().Height(5),
			style:   NewStyle().Height(5),
			content: "hello",
		},
		{
			name:    "width and height",
			style2:  NewStyle2().Width(10).Height(5),
			style:   NewStyle().Width(10).Height(5),
			content: "hello",
		},
		{
			name:    "max width",
			style2:  NewStyle2().MaxWidth(8),
			style:   NewStyle().MaxWidth(8),
			content: "hello world",
		},
		{
			name:    "max height",
			style2:  NewStyle2().MaxHeight(2),
			style:   NewStyle().MaxHeight(2),
			content: "line1\nline2\nline3",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityPadding(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "uniform padding",
			style2:  NewStyle2().Padding(1),
			style:   NewStyle().Padding(1),
			content: "hello",
		},
		{
			name:    "asymmetric padding",
			style2:  NewStyle2().Padding(1, 2, 3, 4),
			style:   NewStyle().Padding(1, 2, 3, 4),
			content: "hello",
		},
		{
			name:    "padding with width",
			style2:  NewStyle2().Padding(1).Width(10),
			style:   NewStyle().Padding(1).Width(10),
			content: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityMargin(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "uniform margin",
			style2:  NewStyle2().Margin(1),
			style:   NewStyle().Margin(1),
			content: "hello",
		},
		{
			name:    "asymmetric margin",
			style2:  NewStyle2().Margin(1, 2, 3, 4),
			style:   NewStyle().Margin(1, 2, 3, 4),
			content: "hello",
		},
		{
			name:    "margin with width",
			style2:  NewStyle2().Margin(1).Width(10),
			style:   NewStyle().Margin(1).Width(10),
			content: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityBorders(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "normal border",
			style2:  NewStyle2().Border(NormalBorder()),
			style:   NewStyle().Border(NormalBorder()),
			content: "hello",
		},
		{
			name:    "rounded border",
			style2:  NewStyle2().Border(RoundedBorder()),
			style:   NewStyle().Border(RoundedBorder()),
			content: "hello",
		},
		{
			name:    "thick border",
			style2:  NewStyle2().Border(ThickBorder()),
			style:   NewStyle().Border(ThickBorder()),
			content: "hello",
		},
		{
			name:    "border top and bottom",
			style2:  NewStyle2().Border(NormalBorder(), true, false),
			style:   NewStyle().Border(NormalBorder(), true, false),
			content: "hello",
		},
		{
			name:    "border left and right",
			style2:  NewStyle2().Border(NormalBorder(), false, true),
			style:   NewStyle().Border(NormalBorder(), false, true),
			content: "hello",
		},
		{
			name:    "border with width",
			style2:  NewStyle2().Border(RoundedBorder()).Width(10),
			style:   NewStyle().Border(RoundedBorder()).Width(10),
			content: "hello",
		},
		{
			name:    "border with colors",
			style2:  NewStyle2().Border(RoundedBorder()).BorderForeground(Red),
			style:   NewStyle().Border(RoundedBorder()).BorderForeground(Red),
			content: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityAlignment(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name:    "align left",
			style2:  NewStyle2().Width(10).AlignHorizontal(Left),
			style:   NewStyle().Width(10).AlignHorizontal(Left),
			content: "hi",
		},
		{
			name:    "align center",
			style2:  NewStyle2().Width(10).AlignHorizontal(Center),
			style:   NewStyle().Width(10).AlignHorizontal(Center),
			content: "hi",
		},
		{
			name:    "align right",
			style2:  NewStyle2().Width(10).AlignHorizontal(Right),
			style:   NewStyle().Width(10).AlignHorizontal(Right),
			content: "hi",
		},
		{
			name:    "vertical align top",
			style2:  NewStyle2().Height(5).AlignVertical(Top),
			style:   NewStyle().Height(5).AlignVertical(Top),
			content: "hi",
		},
		{
			name:    "vertical align center",
			style2:  NewStyle2().Height(5).AlignVertical(Center),
			style:   NewStyle().Height(5).AlignVertical(Center),
			content: "hi",
		},
		{
			name:    "vertical align bottom",
			style2:  NewStyle2().Height(5).AlignVertical(Bottom),
			style:   NewStyle().Height(5).AlignVertical(Bottom),
			content: "hi",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func TestStyle2CompatibilityComplex(t *testing.T) {
	tests := []struct {
		name    string
		style2  Style2
		style   Style
		content string
	}{
		{
			name: "bold + color + padding",
			style2: NewStyle2().
				Bold(true).
				Foreground(Color("#FF0000")).
				Padding(1),
			style: NewStyle().
				Bold(true).
				Foreground(Color("#FF0000")).
				Padding(1),
			content: "hello",
		},
		{
			name: "border + padding + margin",
			style2: NewStyle2().
				Border(RoundedBorder()).
				Padding(1).
				Margin(1).
				Width(10),
			style: NewStyle().
				Border(RoundedBorder()).
				Padding(1).
				Margin(1).
				Width(10),
			content: "hello",
		},
		{
			name: "all text attributes",
			style2: NewStyle2().
				Bold(true).
				Italic(true).
				Underline(true).
				Strikethrough(true),
			style: NewStyle().
				Bold(true).
				Italic(true).
				Underline(true).
				Strikethrough(true),
			content: "hello",
		},
		{
			name: "complex styling",
			style2: NewStyle2().
				Bold(true).
				Foreground(Color("#FFFFFF")).
				Background(Color("#000000")).
				Border(RoundedBorder()).
				BorderForeground(Color("#FF00FF")).
				Padding(1, 2).
				Margin(1).
				Width(15).
				AlignHorizontal(Center),
			style: NewStyle().
				Bold(true).
				Foreground(Color("#FFFFFF")).
				Background(Color("#000000")).
				Border(RoundedBorder()).
				BorderForeground(Color("#FF00FF")).
				Padding(1, 2).
				Margin(1).
				Width(15).
				AlignHorizontal(Center),
			content: "hello",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result2 := tc.style2.Render(tc.content)
			result1 := tc.style.Render(tc.content)
			t.Logf("Style output:\n%s\n\nStyle2 output:\n%s\n", result1, result2)
			if result2 != result1 {
				t.Errorf("Style2 output differs from Style:\nStyle2: %q\nStyle:  %q", result2, result1)
			}
		})
	}
}

func BenchmarkStyle2Render(b *testing.B) {
	style := NewStyle2().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		Foreground(Color("#5A56E0")).
		Background(Color("#FF00FF")).
		Width(40).
		Height(10).
		Padding(1, 2, 1, 2).
		Margin(1, 1, 1, 1).
		Border(RoundedBorder()).
		BorderForeground(Color("#00FF00")).
		AlignHorizontal(Center).
		AlignVertical(Center)

	content := "The quick brown fox jumps over the lazy dog."

	b.Logf("Benchmarking Style2 render output:\n%s\n", style.Render(content))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(content)
	}
}

func BenchmarkStyle1Render(b *testing.B) {
	style := NewStyle().
		Bold(true).
		Italic(true).
		Underline(true).
		Strikethrough(true).
		Foreground(Color("#5A56E0")).
		Background(Color("#FF00FF")).
		Width(40).
		Height(10).
		Padding(1, 2, 1, 2).
		Margin(1, 1, 1, 1).
		Border(RoundedBorder()).
		BorderForeground(Color("#00FF00")).
		AlignHorizontal(Center).
		AlignVertical(Center)

	content := "The quick brown fox jumps over the lazy dog."

	b.Logf("Benchmarking Style render output:\n%s\n", style.Render(content))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = style.Render(content)
	}
}
