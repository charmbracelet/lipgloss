package lipgloss

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestStyleRanges(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		ranges   []Range
		expected string
	}{
		{
			name:     "empty ranges",
			input:    "hello world",
			ranges:   []Range{},
			expected: "hello world",
		},
		{
			name:  "single range in middle",
			input: "hello world",
			ranges: []Range{
				NewRange(6, 10, NewStyle().Bold(true)),
			},
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:  "multiple ranges",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 4, NewStyle().Bold(true)),
				NewRange(6, 10, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[3mworld\x1b[0m",
		},
		{
			name:  "overlapping with existing ANSI",
			input: "hello \x1b[32mworld\x1b[0m",
			ranges: []Range{
				NewRange(0, 4, NewStyle().Bold(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[32mworld\x1b[0m",
		},
		{
			name:  "style at start",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 4, NewStyle().Bold(true)),
			},
			expected: "\x1b[1mhello\x1b[0m world",
		},
		{
			name:  "style at end",
			input: "hello world",
			ranges: []Range{
				NewRange(6, 10, NewStyle().Bold(true)),
			},
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:  "multiple styles with gap",
			input: "hello beautiful world",
			ranges: []Range{
				NewRange(0, 4, NewStyle().Bold(true)),
				NewRange(16, 23, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m beautiful \x1b[3mworld\x1b[0m",
		},
		{
			name:  "adjacent ranges",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 4, NewStyle().Bold(true)),
				NewRange(6, 10, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[3mworld\x1b[0m",
		},
	}

	for _, tt := range tests {
		renderer.SetColorProfile(termenv.ANSI)
		t.Run(tt.name, func(t *testing.T) {
			result := StyleRanges(tt.input, tt.ranges)
			if result != tt.expected {
				t.Errorf("StyleRanges() = %q, want %q", result, tt.expected)
			}
		})
	}
}
