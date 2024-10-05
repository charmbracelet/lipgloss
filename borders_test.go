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
			name: "top left title string",
			text: "",
			style: NewStyle().
				Width(10).
				Border(NormalBorder()).
				BorderDecoration(NewBorderDecoration(Top, Left, "TITLE")),
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
				BorderDecoration(NewBorderDecoration(Top, Left, NewStyle().SetString("TITLE").String)),
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
				BorderDecoration(NewBorderDecoration(Top, Left, NewStyle().SetString("TitleTitleTitle").String)),
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
					Top,
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
					Top,
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
					Top,
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
					Top,
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
					Bottom,
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
					Bottom,
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
					Bottom,
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
					Bottom,
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
					Bottom,
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
