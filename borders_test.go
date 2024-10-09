package lipgloss

import "testing"

func TestBorderFunc(t *testing.T) {

	tt := []struct {
		name     string
		text     string
		style    Style
		expected string
	}{
		{
			name: "trunc all string",
			text: "",
			style: NewStyle().
				Width(16).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(BorderTop, Left, "LeftLeftLeftLeft")).
				BorderDecoration(NewBorderDecoration(BorderTop, Center, "CenterCenterCenter")).
				BorderDecoration(NewBorderDecoration(BorderTop, Right, "RightRightRightRight")),
			expected: `┌LeftL─Cent─Right┐
│                │
└────────────────┘`,
		},
		{
			name: "top left title string",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(BorderTop, Left, "TITLE")),
			expected: `┌TITLE─────┐
│          │
└──────────┘`,
		},
		{
			name: "top left title stringer",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(BorderTop, Left, NewStyle().SetString("TITLE").String)),
			expected: `┌TITLE─────┐
│          │
└──────────┘`,
		},
		{
			name: "top left very long title stringer",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(BorderTop, Left, NewStyle().SetString("TitleTitleTitle").String)),
			expected: `┌TitleTitle┐
│          │
└──────────┘`,
		},
		{
			name: "top left title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderTop,
					Left,
					func(width int, middle string) string {
						return "TITLE"
					},
				)),
			expected: `┌TITLE─────┐
│          │
└──────────┘`,
		},
		{
			name: "top center title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderTop,
					Center,
					func(width int, middle string) string {
						return "TITLE"
					},
				)),
			expected: `┌──TITLE───┐
│          │
└──────────┘`,
		},
		{
			name: "top center title even",
			text: "",
			style: NewStyle().
				Width(11).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderTop,
					Center,
					func(width int, middle string) string {
						return "TITLE"
					},
				)),
			expected: `┌───TITLE───┐
│           │
└───────────┘`,
		},
		{
			name: "top right title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderTop,
					Right,
					func(width int, middle string) string {
						return "TITLE"
					},
				)),
			expected: `┌─────TITLE┐
│          │
└──────────┘`,
		},
		{
			name: "bottom left title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderBottom,
					Left,
					func(width int, middle string) string {
						return "STATUS"
					},
				)),
			expected: `┌──────────┐
│          │
└STATUS────┘`,
		},
		{
			name: "bottom center title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderBottom,
					Center,
					func(width int, middle string) string {
						return "STATUS"
					},
				)),
			expected: `┌──────────┐
│          │
└──STATUS──┘`,
		},
		{
			name: "bottom center title odd",
			text: "",
			style: NewStyle().
				Width(11).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderBottom,
					Center,
					func(width int, middle string) string {
						return "STATUS"
					},
				)),
			expected: `┌───────────┐
│           │
└──STATUS───┘`,
		},
		{
			name: "bottom right title",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderBottom,
					Right,
					func(width int, middle string) string {
						return "STATUS"
					},
				)),
			expected: `┌──────────┐
│          │
└────STATUS┘`,
		},
		{
			name: "bottom right padded title",
			text: "",
			style: NewStyle().
				Width(12).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(
					BorderBottom,
					Right,
					func(width int, middle string) string {
						return "│STATUS│" + middle
					},
				)),
			expected: `┌────────────┐
│            │
└───│STATUS│─┘`,
		},
	}

	for i, tc := range tt {
		res := tc.style.Render(tc.text)
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}

}

func TestBorders(t *testing.T) {
	tt := []struct {
		name     string
		text     string
		style    Style
		expected string
	}{
		{
			name:  "border with width",
			text:  "",
			style: NewStyle().Width(10).Border(NormalBorder()),
			expected: `┌──────────┐
│          │
└──────────┘`,
		},
		{
			name:  "top center title",
			text:  "HELLO",
			style: NewStyle().Border(NormalBorder()),
			expected: `┌─────┐
│HELLO│
└─────┘`,
		},
	}

	for i, tc := range tt {
		res := tc.style.Render(tc.text)
		if res != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}

}

func TestTruncateWidths(t *testing.T) {

	tt := []struct {
		name     string
		widths   [3]int
		width    int
		expected [3]int
	}{
		{
			name:     "lll-cc-rrr",
			widths:   [3]int{10, 10, 10},
			width:    10,
			expected: [3]int{3, 2, 3},
		},
		{
			name:     "lll-ccc-rrr",
			widths:   [3]int{10, 10, 10},
			width:    12,
			expected: [3]int{3, 3, 4},
		},
		{
			name:     "lllll-rrrr",
			widths:   [3]int{10, 0, 10},
			width:    10,
			expected: [3]int{5, 0, 4},
		},
		{
			name:     "lllllll-rr",
			widths:   [3]int{10, 0, 2},
			width:    10,
			expected: [3]int{7, 0, 2},
		},
		{
			name:     "ll-rrrrrrr",
			widths:   [3]int{2, 0, 20},
			width:    10,
			expected: [3]int{2, 0, 7},
		},
		{
			name:     "lll-cc----",
			widths:   [3]int{10, 10, 0},
			width:    10,
			expected: [3]int{3, 2, 0},
		},
		{
			name:     "----cc-rrr",
			widths:   [3]int{0, 10, 10},
			width:    10,
			expected: [3]int{0, 3, 3},
		},
	}

	for i, tc := range tt {
		var result [3]int

		result[0], result[1], result[2] = truncateWidths(tc.widths[0], tc.widths[1], tc.widths[2], tc.width)
		if result != tc.expected {
			t.Errorf("Test %d, expected:`%v`Actual output:`%v`", i, tc.expected, result)
		}
	}

}
