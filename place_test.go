package lipgloss

import "testing"

func TestPlace(t *testing.T) {
	type test struct {
		name     string
		result   string
		expected string
	}
	tests := []test{
		{"pos Left,Top", Place(5, 3, Left, Top, "A", WithWhitespaceChars(".")), "A....\n.....\n....."},
		{"pos Center,Top", Place(5, 3, Center, Top, "A", WithWhitespaceChars(".")), "..A..\n.....\n....."},
		{"pos Right,Top", Place(5, 3, Right, Top, "A", WithWhitespaceChars(".")), "....A\n.....\n....."},

		{"pos Left,Center", Place(5, 3, Left, Center, "A", WithWhitespaceChars(".")), ".....\nA....\n....."},
		{"pos Center,Center", Place(5, 3, Center, Center, "A", WithWhitespaceChars(".")), ".....\n..A..\n....."},
		{"pos Right,Center", Place(5, 3, Right, Center, "A", WithWhitespaceChars(".")), ".....\n....A\n....."},

		{"pos Left,Bottom", Place(5, 3, Left, Bottom, "A", WithWhitespaceChars(".")), ".....\n.....\nA...."},
		{"pos Center,Bottom", Place(5, 3, Center, Bottom, "A", WithWhitespaceChars(".")), ".....\n.....\n..A.."},
		{"pos Right,Bottom", Place(5, 3, Right, Bottom, "A", WithWhitespaceChars(".")), ".....\n.....\n....A"},

		{"pos 0.01,0.01", Place(5, 3, 0.01, 0.01, "A", WithWhitespaceChars(".")), "A....\n.....\n....."},
		{"pos 0.99,0.99", Place(5, 3, 0.99, 0.99, "A", WithWhitespaceChars(".")), ".....\n.....\n....A"},

		{"pos 0.2,Top", Place(5, 3, 0.2, Top, "A", WithWhitespaceChars(".")), ".A...\n.....\n....."},
		{"pos 0.8,Top", Place(5, 3, 0.8, Top, "A", WithWhitespaceChars(".")), "...A.\n.....\n....."},

		{"pos Right,0.2", Place(3, 5, Right, 0.2, "A", WithWhitespaceChars(".")), "...\n..A\n...\n...\n..."},
		{"pos Right,0.8", Place(3, 5, Right, 0.8, "A", WithWhitespaceChars(".")), "...\n...\n...\n..A\n..."},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.result != test.expected {
				t.Errorf("Got \n%s\n, expected \n%s\n", test.result, test.expected)
			}
		})
	}
}
