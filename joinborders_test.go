package lipgloss

import "testing"

func TestJoinBorderVertical(t *testing.T) {
	s := NewStyle().Border(NormalBorder())
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinBorderVertical(s, Left, "A", "BBBB"), "┌────┐\n│A   │\n├────┤\n│BBBB│\n└────┘"},
		{"pos1", JoinBorderVertical(s, Right, "A", "BBBB"), "┌────┐\n│   A│\n├────┤\n│BBBB│\n└────┘"},
		{"pos0.25", JoinBorderVertical(s, 0.25, "A", "BBBB"), "┌────┐\n│ A  │\n├────┤\n│BBBB│\n└────┘"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}

func TestJoinBorderHorizontal(t *testing.T) {
	s := NewStyle().Border(NormalBorder())
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinBorderHorizontal(s, Top, "A", "B\nB\nB\nB"), "┌─┬─┐\n│A│B│\n│ │B│\n│ │B│\n│ │B│\n└─┴─┘"},
		{"pos1", JoinBorderHorizontal(s, Bottom, "A", "B\nB\nB\nB"), "┌─┬─┐\n│ │B│\n│ │B│\n│ │B│\n│A│B│\n└─┴─┘"},
		{"pos0.25", JoinBorderHorizontal(s, 0.25, "A", "B\nB\nB\nB"), "┌─┬─┐\n│ │B│\n│A│B│\n│ │B│\n│ │B│\n└─┴─┘"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}
