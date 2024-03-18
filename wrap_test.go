package lipgloss

import "testing"

func TestWordwrap(t *testing.T) {
	for _, tc := range []struct {
		input    string
		width    int
		expected string
	}{
		{
			"A",
			1,
			"A",
		},
		{
			"A",
			3,
			"A  ",
		},
		{
			"Hello, 世界",
			6,
			"Hello,\n世界 ",
		},
		{
			"Hello",
			3,
			"Hel\nlo ",
		},
	} {
		t.Run(tc.input, func(t *testing.T) {
			res := NewStyle().Width(tc.width).Render(tc.input)
			if tc.expected != res {
				t.Errorf("got %q; want %q", res, tc.expected)
			}
		})
	}
}
