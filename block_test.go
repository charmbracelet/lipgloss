package lipgloss

import "testing"

func TestBlockRender(t *testing.T) {
	cases := []struct {
		name     string
		block    Block
		content  string
		expected string
	}{
		{
			name:     "empty block",
			block:    Block{},
			expected: "",
		},
		{
			name: "block with size",
			block: NewBlock().Width(5).
				Height(3),
			expected: NewStyle().Width(5).Height(3).Render(""),
		},
		{
			name: "block with content",
			block: NewBlock().Width(7).
				Height(3),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(3).Render("Hi"),
		},
		{
			name: "block with border",
			block: NewBlock().Width(7).
				Height(3).
				Border(RoundedBorder()),
			content:  "Hi",
			expected: NewStyle().Border(RoundedBorder()).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "padding uniform",
			block: NewBlock().Width(7).
				Height(5).
				Padding(1),
			content:  "Hi",
			expected: NewStyle().Padding(1).Width(7).Height(5).Render("Hi"),
		},
		{
			name: "padding asymmetric",
			block: NewBlock().Width(9).
				Height(5).
				Padding(1, 2, 1, 1),
			content:  "Hi",
			expected: NewStyle().Padding(1, 2, 1, 1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "border top and bottom",
			block: NewBlock().Width(7).
				Height(3).
				Border(NormalBorder(), true, false),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), true, false).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "border left and right",
			block: NewBlock().Width(7).
				Height(3).
				Border(NormalBorder(), false, true),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), false, true).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "zero size explicit",
			block: NewBlock().Width(0).
				Height(0),
			content:  "",
			expected: "",
		},
		{
			name: "margin uniform",
			block: NewBlock().Width(7).
				Height(3).
				Margin(1),
			content:  "Hi",
			expected: NewStyle().Margin(1).Width(7).Height(3).Render("Hi"),
		},
		{
			name: "margin asymmetric",
			block: NewBlock().Width(9).
				Height(5).
				Margin(1, 2, 1, 1),
			content:  "Hi",
			expected: NewStyle().Margin(1, 2, 1, 1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "margin with border rounded",
			block: NewBlock().Width(9).
				Height(5).
				Margin(1).
				Border(RoundedBorder()),
			content:  "Hi",
			expected: NewStyle().Border(RoundedBorder()).Margin(1).Width(9).Height(5).Render("Hi"),
		},
		{
			name: "margin with border normal left+right",
			block: NewBlock().Width(9).
				Height(5).
				Margin(1).
				Border(NormalBorder(), false, true),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder(), false, true).Margin(1).Width(9).Height(5).Render("Hi"),
		},
		{
			name:     "style with colors (content only)",
			block:    NewBlock().Width(7).Height(3),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(3).Render("Hi"),
		},
		{
			name: "borders with per-side colors",
			block: NewBlock().Width(9).Height(5).
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
			block: NewBlock().Width(11).Height(7).
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
			block: NewBlock().Width(2).Height(2).
				Border(NormalBorder()),
			content:  "Hi",
			expected: NewStyle().Border(NormalBorder()).Width(2).Height(2).Render("Hi"),
		},
		// Max constraints
		{
			name:     "max width wraps/clips content",
			block:    NewBlock().MaxWidth(5),
			content:  "HelloWorld",
			expected: NewStyle().MaxWidth(5).Render("HelloWorld"),
		},
		{
			name:     "max height clips lines",
			block:    NewBlock().MaxHeight(2),
			content:  "A\nB\nC",
			expected: NewStyle().MaxHeight(2).Render("A\nB\nC"),
		},
		// Horizontal alignment
		{
			name:     "align horizontal left",
			block:    NewBlock().Width(9).Height(3).AlignHorizontal(Left),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Left).Render("Hi"),
		},
		{
			name:     "align horizontal center",
			block:    NewBlock().Width(9).Height(3).AlignHorizontal(Center),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Center).Render("Hi"),
		},
		{
			name:     "align horizontal right",
			block:    NewBlock().Width(9).Height(3).AlignHorizontal(Right),
			content:  "Hi",
			expected: NewStyle().Width(9).Height(3).AlignHorizontal(Right).Render("Hi"),
		},
		// Vertical alignment
		{
			name:     "align vertical top",
			block:    NewBlock().Width(7).Height(5).AlignVertical(Top),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Top).Render("Hi"),
		},
		{
			name:     "align vertical center",
			block:    NewBlock().Width(7).Height(5).AlignVertical(Center),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Center).Render("Hi"),
		},
		{
			name:     "align vertical bottom",
			block:    NewBlock().Width(7).Height(5).AlignVertical(Bottom),
			content:  "Hi",
			expected: NewStyle().Width(7).Height(5).AlignVertical(Bottom).Render("Hi"),
		},
		// Combined border + padding + center align
		{
			name:     "border + padding + center alignment",
			block:    NewBlock().Width(11).Height(7).Border(RoundedBorder()).Padding(1).Align(Center, Center),
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
