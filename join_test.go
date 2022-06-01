package lipgloss

import "testing"

func TestJoinVertical(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"por0", JoinVertical(0, "A", "BBBB"), "A   \nBBBB"},
		{"pos1", JoinVertical(1, "A", "BBBB"), "   A\nBBBB"},
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
