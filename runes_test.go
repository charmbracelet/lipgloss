package lipgloss

import (
	"strings"
	"testing"
)

func TestStyleRunes(t *testing.T) {
	t.Parallel()

	matchedStyle := NewStyle().Reverse(true)
	unmatchedStyle := NewStyle()

	tt := []struct {
		input    string
		indices  []int
		expected string
	}{
		{
			"hello",
			[]int{0},
			"\x1b[7mh\x1b[0mello",
		},
		{
			"你好",
			[]int{1},
			"你\x1b[7m好\x1b[0m",
		},
		{
			"hello 你好",
			[]int{6, 7},
			"hello \x1b[7m你好\x1b[0m",
		},
		{
			"hello",
			[]int{1, 3},
			"h\x1b[7me\x1b[0ml\x1b[7ml\x1b[0mo",
		},
		{
			"你好",
			[]int{0, 1},
			"\x1b[7m你好\x1b[0m",
		},
	}

	fn := func(str string, indices []int) string {
		return StyleRunes(str, indices, matchedStyle, unmatchedStyle)
	}

	for i, tc := range tt {
		res := fn(tc.input, tc.indices)
		if fn(tc.input, tc.indices) != tc.expected {
			t.Errorf("Test %d, expected:\n\n`%s`\n`%s`\n\nActual Output:\n\n`%s`\n`%s`\n\n",
				i, tc.expected, formatEscapes(tc.expected),
				res, formatEscapes(res))
		}
	}
}

func formatEscapes(str string) string {
	return strings.ReplaceAll(str, "\x1b", "\\x1b")
}
