package lipgloss

import (
	"strings"
	"testing"
)

func TestStyleRunes(t *testing.T) {
	matchedStyle := NewStyle().Reverse(true)
	unmatchedStyle := NewStyle()

	tt := []struct {
		name     string
		input    string
		indices  []int
		expected string
	}{
		{
			"hello 0",
			"hello",
			[]int{0},
			"\x1b[7mh\x1b[0mello",
		},
		{
			"你好 1",
			"你好",
			[]int{1},
			"你\x1b[7m好\x1b[0m",
		},
		{
			"hello 你好 6,7",
			"hello 你好",
			[]int{6, 7},
			"hello \x1b[7m你好\x1b[0m",
		},
		{
			"hello 1,3",
			"hello",
			[]int{1, 3},
			"h\x1b[7me\x1b[0ml\x1b[7ml\x1b[0mo",
		},
		{
			"你好 0,1",
			"你好",
			[]int{0, 1},
			"\x1b[7m你好\x1b[0m",
		},
	}

	fn := func(str string, indices []int) string {
		return StyleRunes(str, indices, matchedStyle, unmatchedStyle)
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			res := fn(tc.input, tc.indices)
			if res != tc.expected {
				t.Errorf("Expected:\n\n`%s`\n`%s`\n\nActual Output:\n\n`%s`\n`%s`\n\n",
					tc.expected, formatEscapes(tc.expected),
					res, formatEscapes(res))
			}
		})
	}
}

func formatEscapes(str string) string {
	return strings.ReplaceAll(str, "\x1b", "\\x1b")
}
