package lipgloss

import (
	"testing"

	"github.com/muesli/termenv"
)

func TestStyleRange(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		rng      Range
		expected string
	}{
		{
			name:     "empty ranges",
			input:    "hello world",
			rng:      Range{},
			expected: "hello world",
		},
		{
			name:     "single range in middle",
			input:    "hello world",
			rng:      NewRange(6, 11, NewStyle().Bold(true)),
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:     "multiple ranges",
			input:    "hello world",
			rng:      NewRange(0, 5, NewStyle().Bold(true)),
			expected: "\x1b[1mhello\x1b[0m world",
		},
		{
			name:     "overlapping with existing ANSI",
			input:    "hello \x1b[32mworld\x1b[0m",
			rng:      NewRange(0, 5, NewStyle().Bold(true)),
			expected: "\x1b[1mhello\x1b[0m \x1b[32mworld\x1b[0m",
		},
		{
			name:     "style at start",
			input:    "hello world",
			rng:      NewRange(0, 5, NewStyle().Bold(true)),
			expected: "\x1b[1mhello\x1b[0m world",
		},
		{
			name:     "style at end",
			input:    "hello world",
			rng:      NewRange(6, 11, NewStyle().Bold(true)),
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:     "multiple styles with gap",
			input:    "hello beautiful world",
			rng:      NewRange(0, 5, NewStyle().Bold(true)),
			expected: "\x1b[1mhello\x1b[0m beautiful world",
		},
		{
			name:     "adjacent ranges",
			input:    "hello world",
			rng:      NewRange(6, 11, NewStyle().Italic(true)),
			expected: "hello \x1b[3mworld\x1b[0m",
		},
		{
			name:     "wide-width characters",
			input:    "Hello 你好 世界",
			rng:      NewRange(11, 50, NewStyle().Bold(true)), // "世界"
			expected: "Hello 你好 \x1b[1m世界\x1b[0m",
		},
	}

	for _, tt := range tests {
		renderer.SetColorProfile(termenv.ANSI)
		t.Run(tt.name, func(t *testing.T) {
			result := StyleRange(tt.input, tt.rng.Start, tt.rng.End, tt.rng.Style)
			if result != tt.expected {
				t.Errorf("StyleRanges()\n got = %q\nwant = %q\n", result, tt.expected)
			}
		})
	}
}

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
				NewRange(6, 11, NewStyle().Bold(true)),
			},
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:  "multiple ranges",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),
				NewRange(6, 11, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[3mworld\x1b[0m",
		},
		{
			name:  "overlapping with existing ANSI",
			input: "hello \x1b[32mworld\x1b[0m",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[32mworld\x1b[0m",
		},
		{
			name:  "style at start",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),
			},
			expected: "\x1b[1mhello\x1b[0m world",
		},
		{
			name:  "style at end",
			input: "hello world",
			ranges: []Range{
				NewRange(6, 11, NewStyle().Bold(true)),
			},
			expected: "hello \x1b[1mworld\x1b[0m",
		},
		{
			name:  "multiple styles with gap",
			input: "hello beautiful world",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),
				NewRange(16, 23, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m beautiful \x1b[3mworld\x1b[0m",
		},
		{
			name:  "adjacent ranges",
			input: "hello world",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),
				NewRange(6, 11, NewStyle().Italic(true)),
			},
			expected: "\x1b[1mhello\x1b[0m \x1b[3mworld\x1b[0m",
		},
		{
			name:  "wide-width characters",
			input: "Hello 你好 世界",
			ranges: []Range{
				NewRange(0, 5, NewStyle().Bold(true)),    // "Hello"
				NewRange(7, 10, NewStyle().Italic(true)), // "你好"
				NewRange(11, 50, NewStyle().Bold(true)),  // "世界"
			},
			expected: "\x1b[1mHello\x1b[0m \x1b[3m你好\x1b[0m \x1b[1m世界\x1b[0m",
		},
	}

	for _, tt := range tests {
		renderer.SetColorProfile(termenv.ANSI)
		t.Run(tt.name, func(t *testing.T) {
			result := StyleRanges(tt.input, tt.ranges)
			if result != tt.expected {
				t.Errorf("StyleRanges()\n got = %q\nwant = %q\n", result, tt.expected)
			}
		})
	}
}
