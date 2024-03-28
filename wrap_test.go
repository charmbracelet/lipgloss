package lipgloss

import (
	"testing"
)

var cases = []struct {
	name     string
	input    string
	expected string
	width    int
}{
	{
		name:     "simple",
		input:    "I really \x1B[38;2;249;38;114mlove\x1B[0m Go!",
		expected: "I really\n\x1B[38;2;249;38;114mlove\x1B[0m Go!",
		width:    8,
	},
	{
		name:     "passthrough",
		input:    "hello world",
		expected: "hello world",
		width:    11,
	},
	{
		name:     "longer",
		input:    "the quick brown foxxxxxxxxxxxxxxxx jumped over the lazy dog.",
		expected: "the quick brown\nfoxxxxxxxxxxxxxx\nxx jumped over\nthe lazy dog.",
		width:    16,
	},
	{
		name:     "long input",
		input:    "Rotated keys for a-good-offensive-cheat-code-incorporated/animal-like-law-on-the-rocks.",
		expected: "Rotated keys for a-good-offensive-cheat-code-incorporated/animal-like-law-on\n-the-rocks.",
		width:    76,
	},
	{
		name:     "long input2",
		input:    "Rotated keys for a-good-offensive-cheat-code-incorporated/crypto-line-operating-system.",
		expected: "Rotated keys for a-good-offensive-cheat-code-incorporated/crypto-line-operat\ning-system.",
		width:    76,
	},
}

func TestWrapWordwrap(t *testing.T) {
	for i, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			output := wrap(tc.input, tc.width)
			if output != tc.expected {
				t.Errorf("case %d, expected %q, got %q", i+1, tc.expected, output)
			}
		})
	}
}
