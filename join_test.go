package lipgloss

import "testing"

func TestJoinVertical(t *testing.T) {
	type test struct {
		result   string
		expected string
	}
	tests := []test{
		{JoinVertical(0, "A", "BBBB"), "A   \nBBBB"},
		{JoinVertical(1, "A", "BBBB"), "   A\nBBBB"},
		{JoinVertical(0.25, "A", "BBBB"), " A  \nBBBB"},
	}

	for _, test := range tests {
		if test.result != test.expected {
			t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
		}
	}
}
