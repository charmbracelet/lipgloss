package lipgloss

import "testing"

func TestJoinVertical(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinVertical(Left, "A", "BBBB"), "A   \nBBBB"},
		{"pos1", JoinVertical(Right, "A", "BBBB"), "   A\nBBBB"},
		{"pos0.25", JoinVertical(0.25, "A", "BBBB"), " A  \nBBBB"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}

func TestJoinHorizontal(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos0", JoinHorizontal(Top, "A", "B\nB\nB\nB"), "AB\n B\n B\n B"},
		{"pos1", JoinHorizontal(Bottom, "A", "B\nB\nB\nB"), " B\n B\n B\nAB"},
		{"pos0.25", JoinHorizontal(0.25, "A", "B\nB\nB\nB"), " B\nAB\n B\n B"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}
